package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/patrixr/glue/pkg/luatools"
	"github.com/patrixr/q"
	lua "github.com/yuin/gopher-lua"
)

type Glue struct {
	q.Eventful

	lstate *lua.LState

	Stack          GlueStack
	ExecutionTrace []Trace
	DryRun         bool
	Verbose        bool
	Log            *GlueLogger
	Done           bool
	Unsafe         bool
	Modules        []*GlueModule
	Annotations    luatools.LuaAnnotations
	UserSelector   Selector
	FailFast       bool
	Cache          q.Cache[string]
	Context        context.Context
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

	L := lua.NewState(lua.Options{
		SkipOpenLibs: true,
	})

	if err := luatools.LoadSafeLibs(L); err != nil {
		panic(err.Error())
	}

	glue := &Glue{
		Eventful:     q.NewEventEmitter(ctx, 1),
		DryRun:       options.DryRun,
		Verbose:      options.Verbose,
		UserSelector: NewSelectorWithPrefix(options.Selector, []string{RootLevel}),
		Log:          logger,
		lstate:       L,
		Cache:        q.NewInMemoryCache[string](time.Hour * 8760),
		Context:      ctx,
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

	_, errors = glue.Fire(EV_GLUE_END, glue)

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

func (glue *Glue) NotifyError(err error) {
	glue.lstate.RaiseError(err.Error())
}

func (glue *Glue) AtActiveLevel() (bool, error) {
	script := glue.Stack.ActiveScript()
	return glue.UserSelector.Test(
		q.Map(script.GroupStack, func(grp *GlueCodeGroup) string {
			return grp.Name
		}),
	)
}

func (glue *Glue) ExecuteString(script string) error {
	if glue.Done {
		return errors.New("Unable to reuse the same Glue instance")
	}

	glue.Stack.PushScript(":memory:", STRING)

	defer glue.Stack.PopScript()

	if err := glue.lstate.DoString(script); err != nil {
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

	return glue.lstate.DoFile(path)
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
	glue.lstate.Close()
}

func (glue *Glue) Result() (bool, int, []Trace) {
	failedTraces := q.Filter(glue.ExecutionTrace, func(trace Trace) bool {
		return trace.Error != nil
	})

	return len(failedTraces) == 0, len(failedTraces), glue.ExecutionTrace
}
