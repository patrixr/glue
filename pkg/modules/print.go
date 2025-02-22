package modules

import (
	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	// @auteur("Helpers/Print")
	//
	// # Print
	//
	// Prints a message or object to the log.
	// This helper takes one argument, `obj`, which can be any type.
	// The `obj` is converted to a string and logged as an informational message.
	//
	// E.g.
	//
	// ```lua
	// print("hello world")
	// ```
	//
	Registry.RegisterModule(func(glue *core.Glue) error {

		glue.Plug("print", core.FUNCTION).
			Brief("Print a string").
			Arg("obj", ANY, "the message or object to log").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				glue.Log.Info(args.Get(0).String())
				return nil, nil
			})

		return nil
	})
}
