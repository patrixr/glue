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
	nesting           []string

	ExecutionTrace []FunctionCall
	DryRun         bool
	Log            *GlueLogger
	Done           bool
	Unsafe         bool
	Modules        []*GlueModule
	Annotations    luatools.LuaAnnotations
	UserSelector   Selector
}

type GlueOptions struct {
	DryRun   bool
	Selector string
}

type FunctionCall struct {
	Name string
	Args []string
}

func NewGlue() *Glue {
	return NewGlueWithOptions(GlueOptions{
		DryRun: false,
	})
}

func NewGlueWithOptions(options GlueOptions) *Glue {
	logger := CreateLogger()

	L := lua.NewState(lua.Options{
		SkipOpenLibs: true,
	})

	glue := &Glue{
		DryRun:       options.DryRun,
		UserSelector: NewSelectorWithPrefix(options.Selector, []string{RootLevel}),
		Log:          logger,

		lstate:  L,
		nesting: []string{RootLevel},
	}

	InstallNativeGlueModules(glue)

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

func (glue *Glue) AtActiveLevel() (bool, error) {
	return glue.UserSelector.Test(glue.nesting)
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
