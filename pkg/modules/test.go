package modules

import (
	"github.com/patrixr/glue/pkg/core"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug().
			Name("test").
			Short("Create a test case").
			Arg("name", "string", "A description of the test").
			Arg("fn", "function", "The test implementation").
			Mode(core.NONE).
			Bypass().
			Do(func(L *lua.LState) (int, error) {
				active, err := glue.AtActiveLevel()

				if err != nil {
					return 0, err
				}

				if !active {
					return 0, nil
				}

				L.CheckTypes(1, lua.LTString)
				L.CheckTypes(2, lua.LTFunction)

				name := L.ToString(1)
				fn := L.ToFunction(2)

				glue.RegisterTest(name, func() {
					L.CallByParam(lua.P{
						Fn:      fn,
						NRet:    0,
						Protect: false,
					})
				})

				return 0, nil
			})

		return nil
	})
}
