package modules

import (
	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/luatools"
)

const ABOUT_CACHE_KEY = "annotation:about"

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {

		// When a module is run, we check if there is an annotation from the user
		// and we attach it to the trace.
		glue.On(core.EV_NEW_TRACE, func(_ string, data any) error {
			note, ok := glue.Stack.CurrentGroup().Get(ABOUT_CACHE_KEY)

			if !ok || len(note) == 0 {
				return nil
			}

			trace, ok := data.(*core.Trace)

			if ok {
				trace.About = note
			}

			return nil
		})

		glue.Plug().
			Name("note").
			Short("Annotate the current group with some details").
			Arg("brief", "string", "short explanation of the next step").
			Mode(core.NONE).
			Bypass().
			Do(luatools.StrFunc(func(s string) error {
				glue.Stack.CurrentGroup().Set(ABOUT_CACHE_KEY, s)
				return nil
			}))

		return nil
	})
}
