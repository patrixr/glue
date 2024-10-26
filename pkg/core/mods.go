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
	Name     string
	Short    string
	Long     string
	Mode     Mode
	Examples []string
	fn       lua.LGFunction
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

func (plug *glueplug) Do(fn lua.LGFunction) error {
	if len(plug.name) == 0 {
		return errors.New(
			"Trying to install a module with empty name",
		)
	}

	glue := plug.glue
	name := plug.name
	mode := plug.mode

	wrapped := glue.lstate.NewFunction(
		func(L *lua.LState) int {
			glue.recordTrace(name, L)

			if glue.DryRun && (HasMode(mode, WRITE) || HasMode(mode, NETWORK)) {
				// When doing a dry run, we stub out the
				// write methods with mocks
				inputs := luatools.GetAllArgsAsStrings(L)
				text := fmt.Sprintf(
					"(mocked) %s(%s)", name, strings.Join(inputs, ", "))
				glue.Log.Info(text)
				return 0
			}

			return fn(L)
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

func (glue *Glue) recordTrace(name string, L *lua.LState) {
	glue.ExecutionTrace = append(
		glue.ExecutionTrace, FunctionCall{
			name,
			luatools.GetAllArgsAsStrings(L),
		})
}
