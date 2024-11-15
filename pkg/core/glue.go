package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/patrixr/glue/pkg/luatools"
	"github.com/patrixr/q"
	lua "github.com/yuin/gopher-lua"
)

type Glue struct {
	afterScriptFuncs  []func() error
	beforeScriptFuncs []func() error
	lstate            *lua.LState
	fileStack         []string
	nesting           []string

	ExecutionTrace []Trace
	DryRun         bool
	Log            *GlueLogger
	Done           bool
	Unsafe         bool
	Modules        []*GlueModule
	Annotations    luatools.LuaAnnotations
	UserSelector   Selector
	FailFast       bool
}

type GlueOptions struct {
	DryRun   bool
	Selector string
}

type Trace struct {
	Name  string
	Args  []string
	Error error
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

	if err := luatools.LoadSafeLibs(L); err != nil {
		panic(err.Error())
	}

	glue := &Glue{
		DryRun:       options.DryRun,
		UserSelector: NewSelectorWithPrefix(options.Selector, []string{RootLevel}),
		Log:          logger,
		lstate:       L,
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

func (glue *Glue) NotifyError(err error) {
	glue.lstate.RaiseError(err.Error())
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
	nesting := glue.nesting

	glue.fileStack = append(glue.fileStack, file)
	glue.nesting = []string{RootLevel}

	defer func() {
		glue.fileStack = glue.fileStack[:len(glue.fileStack)-1]
		glue.nesting = nesting
	}()

	return glue.lstate.DoFile(file)
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

func (glue *Glue) Result() (bool, []Trace) {
	failedTraces := q.Filter(glue.ExecutionTrace, func(trace Trace) bool {
		return trace.Error != nil
	})

	return len(failedTraces) == 0, failedTraces
}

func runAll(funcs []func() error) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
