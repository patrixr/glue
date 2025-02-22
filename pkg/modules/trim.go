package modules

import (
	"strings"

	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
	"github.com/patrixr/q"
)

func init() {
	// @auteur("Helpers/Trim")
	//
	// # Trim
	//
	// The trim method removes leading and trailing whitespaces from a string.
	// Itcan also be used to clean up multi-line strings in Lua scripts.
	//
	// Here's an example of how you might use it in a Lua script:
	//
	// ```lua
	// local text = [[
	//     This is a multi-line string
	//        that has inconsistent indentation.
	// ]]
	//
	// local trimmed_text = trim(text)
	// ```
	// This would result in the following text:
	//
	// ```text
	// This is a multi-line string
	//    that has inconsistent indentation.
	// ```
	Registry.RegisterModule(
		func(glue *core.Glue) error {

			glue.Plug("trim", core.FUNCTION).
				Brief("Trims the extra indentation of a multi-line string").
				Arg("txt", STRING, "the text to trim").
				Return(STRING, "the trimmed text").
				Do(func(R Runtime, args *Arguments) (RTValue, error) {
					s := args.EnsureString(0).String()
					return R.String(strings.TrimSpace(q.TrimIndent(s))), nil
				})

			return nil
		})
}
