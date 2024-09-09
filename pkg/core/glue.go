package core

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/patrixr/glue/pkg/luatools"
	lua "github.com/yuin/gopher-lua"
)

type Glue struct {
	afterScriptFuncs  []func() error
	beforeScriptFuncs []func() error
	lstate            *lua.LState
	fileStack         []string

	ExecutionTrace []FunctionCall
	DryRun         bool
	Log            Logger
	Done           bool
	Unsafe         bool
}

type GlueOptions struct {
	Unsafe bool
	DryRun bool
}

type FunctionCall struct {
	Name string
	Args []string
}

func NewGlue() *Glue {
	return NewGlueWithOptions(GlueOptions{
		Unsafe: false,
		DryRun: false,
	})
}

func NewGlueWithOptions(options GlueOptions) *Glue {
	logger := CreateLogger()

	if options.DryRun && options.Unsafe {
		logger.Error("Unable to initialize glue in both DryRun and Unsafe modes")
		os.Exit(1)
	}

	L := lua.NewState(lua.Options{
		SkipOpenLibs: !options.Unsafe,
	})

	glue := &Glue{
		DryRun: options.DryRun,
		lstate: L,
		Log:    logger,
	}

	installNativeGlueModules(glue)

	return glue
}

func (glue *Glue) Execute(file string) error {
	if glue.Done {
		return errors.New("Unable to reuse the same Glue instance")
	}

	if err := glue.RunBeforeScript(); err != nil {
		return err
	}

	if err := glue.RunFileRaw(file); err != nil {
		return err
	}

	return glue.RunAfterScript()
}

func (glue *Glue) ExecuteString(script string) error {
	if glue.Done {
		return errors.New("Unable to reuse the same Glue instance")
	}

	if err := glue.RunBeforeScript(); err != nil {
		return err
	}

	if err := glue.lstate.DoString(script); err != nil {
		return err
	}

	return glue.RunAfterScript()
}

func (glue *Glue) RunFileRaw(file string) error {
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

func (glue *Glue) Getwd() (string, error) {
	current, err := glue.GetCurrentScript()
	if err == nil {
		return filepath.Dir(current), nil
	}
	return os.Getwd()
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
			glue.WrapFunc(name, mock),
		)

		if err != nil {
			return err
		}

		return nil
	}

	err := luatools.SetNestedGlobalValue(
		glue.lstate,
		name,
		glue.WrapFunc(name, fn),
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

func (glue *Glue) WrapFunc(name string, fn lua.LGFunction) *lua.LFunction {
	return glue.lstate.NewFunction(func(L *lua.LState) int {
		glue.ExecutionTrace = append(glue.ExecutionTrace, FunctionCall{
			name,
			luatools.GetAllArgsAsStrings(L),
		})
		return fn(L)
	})
}

func runAll(funcs []func() error) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
