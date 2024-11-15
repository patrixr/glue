package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/patrixr/glue/pkg/luatools"
	lua "github.com/yuin/gopher-lua"
)

type Mode uint8

const (
	NONE Mode = 0
	READ Mode = 1 << iota
	WRITE
	NETWORK
)

const DefaultMode = READ | WRITE | NETWORK

func HasMode(b, flag Mode) bool {
	return b&flag != 0
}

// A representation of a module
// Modules can be installed into the Lua state
type GlueModule struct {
	Name       string
	Short      string
	Long       string
	Mode       Mode
	Examples   []string
	Bypass     bool
	fn         luatools.LuaFuncWithError
	mockReturn []lua.LValue
}

// An intermediate builder for creating a module
type glueplug struct {
	name       string
	short      string
	long       string
	mode       Mode
	examples   []string
	annotation luatools.LuaFuncAnnotation
	glue       *Glue
	bypass     bool
	mockReturn []lua.LValue
}

// Entry point for creating a new module
//
// Example:
//
//	glue.Plug().
//	   Name("myModule").
//	   Short("description").
//	   Do(...)
func (glue *Glue) Plug() *glueplug {
	return &glueplug{
		glue: glue,
		mode: DefaultMode,
	}
}

func (plug *glueplug) Name(name string) *glueplug {
	plug.name = name
	return plug
}

func (plug *glueplug) Short(short string) *glueplug {
	plug.short = short
	return plug
}

func (plug *glueplug) Long(lines ...string) *glueplug {
	plug.long = strings.TrimSpace(strings.Join(lines, "\n"))
	return plug
}

func (plug *glueplug) Mode(mode Mode) *glueplug {
	plug.mode = mode
	return plug
}

func (plug *glueplug) Example(ex string) *glueplug {
	plug.examples = append(plug.examples, ex)
	return plug
}

func (plug *glueplug) Bypass() *glueplug {
	plug.bypass = true
	return plug
}

func (plug *glueplug) Arg(name string, valtype string, desc string) *glueplug {
	plug.annotation.Args = append(plug.annotation.Args, luatools.LuaFieldDesc{
		Name: name,
		Type: valtype,
		Desc: desc,
	})
	return plug
}

func (plug *glueplug) Return(valtype string, desc string) *glueplug {
	plug.annotation.Returns = append(plug.annotation.Returns, luatools.LuaReturnDesc{
		Type: valtype,
		Desc: desc,
	})
	return plug
}

func (plug *glueplug) MockReturn(returnArgs ...lua.LValue) *glueplug {
	plug.mockReturn = returnArgs
	return plug
}

func (plug *glueplug) Do(fn luatools.LuaFuncWithError) error {
	if len(plug.name) == 0 {
		return errors.New(
			"Trying to install a module with empty name",
		)
	}

	glue := plug.glue
	name := plug.name
	mode := plug.mode
	bypass := plug.bypass

	wrapped := glue.lstate.NewFunction(
		func(L *lua.LState) int {
			if !bypass {
				active, err := glue.AtActiveLevel()

				if err != nil {
					L.RaiseError(err.Error())
					return 0
				}

				skip := !active || glue.DryRun

				if glue.DryRun && active {
					inputs := luatools.GetAllArgsAsStrings(L)
					text := fmt.Sprintf(
						"(mocked) %s(%s)", name, strings.Join(inputs, ", "))
					glue.Log.Info(text)
				}

				// We don't run functions in dry mode or when the current nesting level is not active
				if skip {
					for _, arg := range plug.mockReturn {
						L.Push(arg)
					}
					return len(plug.mockReturn)
				}
			}

			res, err := fn(L)

			glue.recordTrace(name, L, err)

			if err != nil && glue.FailFast {
				L.RaiseError(err.Error()) // boom
			}

			return res
		})

	path, err := luatools.SetNestedGlobalValue(
		glue.lstate,
		name,
		wrapped,
	)

	glue.Annotations.AddNestedTable(path)

	if err != nil {
		return err
	}

	mod := &GlueModule{
		Name:     name,
		Long:     plug.long,
		Short:    plug.short,
		Mode:     mode,
		Examples: plug.examples,
		fn:       fn,
	}

	glue.Modules = append(glue.Modules, mod)
	glue.Annotations.Add(&luatools.LuaFuncAnnotation{
		Name:    name,
		Args:    plug.annotation.Args,
		Returns: plug.annotation.Returns,
		Desc:    plug.short,
	})

	return err
}

func (glue *Glue) recordTrace(name string, L *lua.LState, err error) {
	glue.ExecutionTrace = append(glue.ExecutionTrace, Trace{
		Name:  name,
		Args:  luatools.GetAllArgsAsStrings(L),
		Error: err,
	})
}
