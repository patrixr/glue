package modules

import (
	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug("assert", core.FUNCTION).
			Brief("Asserts the given boolean and raises and error if problematic").
			Arg("value", BOOL, "the condition to assert on").
			Arg("brief", STRING, "short explanation of the next step").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				if !glue.Testing() {
					return nil, nil
				}

				valid := args.EnsureBool(0).Value()
				msg := args.EnsureString(1).String()

				if !valid {
					panic(msg)
				}

				return nil, nil
			})

		return nil
	})
}
