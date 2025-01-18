package modules

import (
	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug("test", core.FUNCTION).
			Brief("Create a test case").
			Arg("name", STRING, "A description of the test").
			Arg("fn", FUNC, "The test implementation").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				name := args.EnsureString(0).String()
				fn := args.EnsureFunction(0)

				glue.RegisterTest(name, func() {
					R.InvokeFunction(fn)
				})

				return nil, nil
			})

		return nil
	})
}
