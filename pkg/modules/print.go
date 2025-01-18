package modules

import (
	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {

		glue.Plug("print", core.FUNCTION).
			Brief("Print a string").
			Arg("obj", ANY, "the message or object to log").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				glue.Log.Info(args.EnsureString(0).String())
				return nil, nil
			})

		return nil
	})
}
