package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/patrixr/glue/pkg/blueprint"
	"github.com/patrixr/glue/pkg/machine"
	"github.com/patrixr/glue/pkg/runtime"
	"github.com/patrixr/glue/pkg/runtime/lua"
	"github.com/patrixr/q"
)

type Glue struct {
	q.Eventful
	Testable

	Stack        GlueStack
	BluePrint    Blueprint
	Verbose      bool
	Done         bool
	Unsafe       bool
	FailFast     bool
	Log          *GlueLogger
	Modules      []*GluePlugin
	UserSelector Selector
	Cache        q.Cache[string]
	Context      context.Context
	Runtime      runtime.Runtime
	Machine      machine.Machine
}

type GlueOptions struct {
	Selector string
	Verbose  bool
}

func NewGlue() *Glue {
	return NewGlueWithOptions(GlueOptions{})
}

func NewGlueWithOptions(options GlueOptions) *Glue {
	logger := CreateLogger()

	ctx := context.Background()

	glue := &Glue{
		Runtime:      lua.NewLuaRuntime(),
		Eventful:     q.NewEventEmitter(ctx, 1),
		Testable:     NewTestSuite(),
		Verbose:      options.Verbose,
		UserSelector: NewSelectorWithPrefix(options.Selector, []string{RootLevel}),
		Log:          logger,
		Cache:        q.NewInMemoryCache[string](time.Hour * 8760),
		Context:      ctx,
		BluePrint:    nil,
		Machine:      machine.NewLocalMachine(),
	}

	InstallNativeGlueModules(glue)

	return glue
}

// Main entry point for glue
// Given a script file, it compiles it into an executable plan
func (glue *Glue) CompilePlan(file string) (Blueprint, error) {
	if glue.Done {
		return nil, errors.New("Unable to reuse the same Glue instance")
	}

	glue.BluePrint = NewSerialBlueprint("<root>")

	defer func() {
		glue.BluePrint = nil
	}()

	_, errors := glue.Fire(EV_GLUE_PLAN_END, glue)

	if len(errors) > 0 {
		return nil, errors[0]
	}

	if err := glue.execFile(file); err != nil {
		return nil, err
	}

	_, errors = glue.Fire(EV_GLUE_PLAN_END, glue)

	if len(errors) > 0 {
		return nil, errors[0]
	}

	return glue.BluePrint, nil
}

// (internal)
// Checks if the current group is allowed to run based on the user's filter options
func (glue *Glue) canRunGroup(group string) (bool, error) {
	if glue.Testing() {
		return true, nil
	}

	if len(glue.Stack.ExecutionStack) == 0 {
		return false, nil
	}

	script := glue.Stack.ActiveScript()
	groups := q.Map(script.GroupStack, func(grp *GlueCodeGroup) string {
		return grp.Name
	})

	return glue.UserSelector.Test(append(groups, group))
}

// (internal)
// Executes a script from a string
func (glue *Glue) execString(script string) error {
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

// (internal)
// Executes a script from a file
func (glue *Glue) execFile(file string) error {
	path, err := glue.SmartPath(file)

	if err != nil {
		return err
	}

	glue.Stack.PushScript(path, FILE)

	defer glue.Stack.PopScript()

	return glue.Runtime.ExecFile(path)
}

// Getwd returns the working directory of the active script or the current working directory
func (glue *Glue) Getwd() (string, error) {
	if glue.Stack.HasActiveScript() {
		script := glue.Stack.ActiveScript()
		return filepath.Dir(script.Uri), nil
	}

	return os.Getwd()
}

// SmartPath resolves a path to an absolute path
// If called from within a script, it resolves the path relative to the script's directory
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
