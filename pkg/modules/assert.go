package modules

import (
	"github.com/patrixr/glue/pkg/core"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug().
			Name("assert").
			Short("Asserts the given boolean and raises and error if problematic").
			Arg("value", "boolean", "the condition to assert on").
			Arg("brief", "string", "short explanation of the next step").
			Mode(core.NONE).
			Bypass().
			Do(func(L *lua.LState) (int, error) {
				if !glue.Testing() {
					return 0, nil
				}

				valid := L.ToBool(1)
				msg := L.ToString(2)

				if !valid {
					panic(msg)
				}

				return 0, nil
			})

		return nil
	})
}
