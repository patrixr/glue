package modules

import (
	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/luatools"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {

		glue.Plug().
			Name("print").
			Short("Print a string").
			Long("Print a string").
			Arg("obj", "any", "the message or object to log").
			Example("print('Hello, world!')").
			Bypass().
			Do(func(L *lua.LState) (int, error) {
				input := luatools.GetArgAsString(L, 1)
				glue.Log.Info(input)
				return 0, nil
			})

		return nil
	})
}
