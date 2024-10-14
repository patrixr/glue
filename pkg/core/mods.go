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
)

const defaultMode = READ | WRITE

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
type gluePlug struct {
	name     string
	short    string
	long     string
	mode     Mode
	examples []string
	glue     *Glue
}

// Entry point for creating a new module
//
// Example:
//
//	glue.Plug()
//	   .Name("myModule")
//	   .Do(func(L *lua.LState) int { ... })
func (glue *Glue) Plug() *gluePlug {
	return &gluePlug{
		glue: glue,
		mode: defaultMode,
	}
}

func (plug *gluePlug) Name(name string) *gluePlug {
	plug.name = name
	return plug
}

func (plug *gluePlug) Short(short string) *gluePlug {
	plug.short = short
	return plug
}

func (plug *gluePlug) Long(lines ...string) *gluePlug {
	plug.long = strings.TrimSpace(strings.Join(lines, "\n"))
	return plug
}

func (plug *gluePlug) Mode(mode Mode) *gluePlug {
	plug.mode = mode
	return plug
}

func (plug *gluePlug) Example(ex string) *gluePlug {
	plug.examples = append(plug.examples, ex)
	return plug
}

func (plug *gluePlug) Do(fn lua.LGFunction) error {
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

			if glue.DryRun && HasMode(mode, WRITE) {
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

	err := luatools.SetNestedGlobalValue(
		glue.lstate,
		name,
		wrapped,
	)

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

	return err
}

func (glue *Glue) recordTrace(name string, L *lua.LState) {
	glue.ExecutionTrace = append(
		glue.ExecutionTrace, FunctionCall{
			name,
			luatools.GetAllArgsAsStrings(L),
		})
}
