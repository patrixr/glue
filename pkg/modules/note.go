package modules

import (
	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

const ABOUT_CACHE_KEY = "annotation:about"

func init() {
	// @auteur("Modules/Note")
	//
	// # Note
	//
	// The note module allows users to annotate the execution scope.
	// This can be useful for adding context or explanations to different parts of the trace.
	//
	// ```lua
	// note("some information")
	// ```
	Registry.RegisterModule(func(glue *core.Glue) error {

		// When a module is run, we check if there is an annotation from the user
		// and we attach it to the trace.
		// glue.On(core.EV_NEW_TRACE, func(_ string, data any) error {
		// 	note, ok := glue.Stack.CurrentGroup().Get(ABOUT_CACHE_KEY)

		// 	if !ok || len(note) == 0 {
		// 		return nil
		// 	}

		// 	trace, ok := data.(*core.Trace)

		// 	if ok {
		// 		trace.About = note
		// 	}

		// 	return nil
		// })

		glue.Plug("note", core.FUNCTION).
			Brief("Annotate the current group with some details").
			Arg("brief", STRING, "short explanation of the next step").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				// s := args.EnsureString(0)
				// glue.Stack.CurrentGroup().Set(ABOUT_CACHE_KEY, s.String())
				return nil, nil
			})

		return nil
	})
}
