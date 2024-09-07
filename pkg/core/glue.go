package core

import (
	"errors"

	"github.com/patrixr/glue/pkg/luatools"
	lua "github.com/yuin/gopher-lua"
)

type Glue struct {
	afterScriptFuncs  []func() error
	beforeScriptFuncs []func() error
	lstate            *lua.LState
	fileStack         []string

	DryRun bool
	Log    Logger
	Done   bool
}

type GlueOption func(glue *Glue)

func NewGlue(options ...GlueOption) *Glue {
	L := lua.NewState()

	glue := &Glue{
		DryRun: false,
		lstate: L,
		Log:    CreateLogger(),
	}

	for _, opt := range options {
		opt(glue)
	}

	installNativeGlueModules(glue)

	return glue
}

func (glue *Glue) Execute(file string) error {
	if glue.Done {
		return errors.New("Unable to reuse the same Glue instance")
	}

	defer glue.Close()

	if err := glue.RunBeforeScript(); err != nil {
		return err
	}

	if err := glue.RunRaw(file); err != nil {
		return err
	}

	return glue.RunAfterScript()
}

func (glue *Glue) RunRaw(file string) error {
	glue.fileStack = append(glue.fileStack, file)
	err := glue.lstate.DoFile(file)
	glue.fileStack = glue.fileStack[:len(glue.fileStack)-1]
	return err
}

func (glue *Glue) GetCurrentScript() (string, error) {
	if len(glue.fileStack) == 0 {
		return "", errors.New("No script is running at the moment")
	}

	return glue.fileStack[len(glue.fileStack)-1], nil
}

func (glue *Glue) Close() {
	glue.Done = true
	glue.lstate.Close()
}

func (glue *Glue) AddFunction(name string, fn lua.LGFunction) error {
	if len(name) == 0 {
		return errors.New("Trying to register a module with empty name")
	}

	if glue.DryRun {
		mock := glue.generateMockMethod(name)
		err := luatools.SetNestedGlobalValue(
			glue.lstate,
			name,
			mock,
		)

		if err != nil {
			return err
		}

		return nil
	}

	err := luatools.SetNestedGlobalValue(
		glue.lstate,
		name,
		glue.lstate.NewFunction(fn),
	)

	if err != nil {
		return err
	}
	return nil
}

func (glue *Glue) AfterScript(f func() error) {
	glue.afterScriptFuncs = append(glue.afterScriptFuncs, f)
}

func (glue *Glue) BeforeScript(f func() error) {
	glue.beforeScriptFuncs = append(glue.beforeScriptFuncs, f)
}

func (glue *Glue) RunBeforeScript() error {
	return runAll(glue.beforeScriptFuncs)
}

func (glue *Glue) RunAfterScript() error {
	return runAll(glue.afterScriptFuncs)
}

func runAll(funcs []func() error) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
