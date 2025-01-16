package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/patrixr/glue/pkg/luatools"
	"github.com/patrixr/glue/pkg/runtime"
	"github.com/patrixr/glue/pkg/runtime/lua"
	"github.com/patrixr/q"
)

type Glue struct {
	q.Eventful
	Testable

	Stack          GlueStack
	ExecutionTrace []Trace
	Plan           Plan
	DryRun         bool
	Verbose        bool
	RunTests       bool
	Done           bool
	Unsafe         bool
	FailFast       bool
	Log            *GlueLogger
	Modules        []*GlueModule
	Annotations    luatools.LuaAnnotations
	UserSelector   Selector
	Cache          q.Cache[string]
	Context        context.Context
	Runtime        runtime.Runtime
}

type GlueOptions struct {
	DryRun   bool
	Selector string
	Verbose  bool
}

func NewGlue() *Glue {
	return NewGlueWithOptions(GlueOptions{
		DryRun: false,
	})
}

func NewGlueWithOptions(options GlueOptions) *Glue {
	logger := CreateLogger()

	ctx := context.Background()

	glue := &Glue{
		Runtime:      lua.NewLuaRuntime(),
		Eventful:     q.NewEventEmitter(ctx, 1),
		Testable:     NewTestSuite(),
		DryRun:       options.DryRun,
		Verbose:      options.Verbose,
		UserSelector: NewSelectorWithPrefix(options.Selector, []string{RootLevel}),
		Log:          logger,
		Cache:        q.NewInMemoryCache[string](time.Hour * 8760),
		Context:      ctx,
		RunTests:     true,
		Plan:         NewPlan("<root>"),
	}

	InstallNativeGlueModules(glue)

	return glue
}

func (glue *Glue) Execute(file string) error {
	if glue.Done {
		return errors.New("Unable to reuse the same Glue instance")
	}

	_, errors := glue.Fire(EV_GLUE_END, glue)

	if len(errors) > 0 {
		return errors[0]
	}

	if err := glue.RunFileRaw(file); err != nil {
		return err
	}

	fmt.Println("THE PLAN:")
	fmt.Println(glue.Plan)

	if glue.RunTests {
		glue.Test()
	}

	_, errors = glue.Fire(EV_GLUE_END, glue)

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

func (glue *Glue) AtActiveLevel() (bool, error) {
	return true, nil
	// if glue.Testing() {
	// 	return true, nil
	// }

	// if len(glue.Stack.ExecutionStack) == 0 {
	// 	return false, nil
	// }

	// script := glue.Stack.ActiveScript()
	// return glue.UserSelector.Test(
	// 	q.Map(script.GroupStack, func(grp *GlueCodeGroup) string {
	// 		return grp.Name
	// 	}),
	// )
}

func (glue *Glue) ExecuteString(script string) error {
	if glue.Done {
		return errors.New("Unable to reuse the same Glue instance")
	}

	glue.Stack.PushScript(":memory:", STR)

	defer glue.Stack.PopScript()

	if err := glue.Runtime.ExecString(script); err != nil {
		return err
	}

	return nil
}

func (glue *Glue) RunFileRaw(file string) error {
	path, err := glue.SmartPath(file)

	if err != nil {
		return err
	}

	glue.Stack.PushScript(path, FILE)

	defer glue.Stack.PopScript()

	return glue.Runtime.ExecFile(path)
}

func (glue *Glue) Getwd() (string, error) {
	if glue.Stack.HasActiveScript() {
		script := glue.Stack.ActiveScript()
		return filepath.Dir(script.Uri), nil
	}

	return os.Getwd()
}

func (glue *Glue) SmartPath(path string) (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if path == "~" {
		return homedir, nil
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homedir, path[2:]), nil
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	wd, err := glue.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, path), nil
}

func (glue *Glue) Close() {
	glue.Done = true
	glue.Runtime.Close()
}

func (glue *Glue) Result() (bool, int, []Trace) {
	failedTraces := q.Filter(glue.ExecutionTrace, func(trace Trace) bool {
		return trace.Error != nil
	})

	return len(failedTraces) == 0, len(failedTraces), glue.ExecutionTrace
}
