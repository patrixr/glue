package core

import (
	"errors"
	"strings"

	"github.com/golang-cz/textcase"
	"github.com/patrixr/glue/pkg/runtime"
)

type PluginKind uint8

const (
	MODULE PluginKind = 1 << iota
	FUNCTION
)

// A representation of a module
// Modules can be installed into the Lua state
type GluePlugin struct {
	Name       string
	Brief      string
	Args       []runtime.ArgDef
	ReturnType runtime.Type
	Kind       PluginKind
}

// An intermediate builder for creating a module
type plugin struct {
	name       string
	brief      string
	kind       PluginKind
	returnType runtime.Type
	args       []runtime.ArgDef
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
func (glue *Glue) Plug(name string, kind PluginKind) *plugin {
	path := strings.Split(name, ".")

	for _, key := range path {
		runtime.AssertValidSymbolName(key)
	}

	if kind == MODULE {
		path[len(path)-1] = textcase.PascalCase(path[len(path)-1])
	} else {
		path[len(path)-1] = textcase.CamelCase(path[len(path)-1])
	}

	name = strings.Join(path, ".")

	return &plugin{
		glue:       glue,
		name:       name,
		kind:       kind,
		returnType: runtime.NIL,
	}
}

func (plug *plugin) Brief(short string) *plugin {
	plug.brief = short
	return plug
}

func (plug *plugin) Arg(name string, valtype runtime.Type, desc string) *plugin {
	plug.args = append(plug.args, runtime.ArgDef{
		Name: name,
		Type: valtype,
		Desc: desc,
	})
	return plug
}

func (plug *plugin) Return(typ runtime.Type, desc string) *plugin {
	if typ != runtime.NIL && plug.kind == MODULE {
		panic("Glue modules cannot return values, are you looking to create a function?")
	}
	plug.returnType = runtime.TypeWithDesc(typ, desc)
	return plug
}

func (plug *plugin) Do(fn func(R runtime.Runtime, args *runtime.Arguments) (runtime.RTValue, error)) error {
	if len(plug.name) == 0 {
		return errors.New(
			"Trying to install a module with empty name",
		)
	}

	glue := plug.glue
	name := plug.name

	glue.Runtime.SetFunction(
		name,
		plug.brief,
		plug.args,
		func(R runtime.Runtime, args *runtime.Arguments) runtime.RTValue {
			if plug.kind == FUNCTION {
				res, err := fn(R, args)
				if err != nil {
					R.RaiseError("%s", err.Error())
				}
				return res
			}

			glue.BluePrint.Action(name, "", "", func() error {
				_, err := fn(R, args)
				return err
			})

			return nil
		})

	mod := &GluePlugin{
		Name:       name,
		Kind:       plug.kind,
		Brief:      plug.brief,
		Args:       plug.args,
		ReturnType: plug.returnType,
	}

	glue.Modules = append(glue.Modules, mod)

	return nil
}
