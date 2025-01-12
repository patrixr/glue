package modules

import (
	"strings"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/luatools"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug().
			Name("capitalize").
			Short("Uppercase the first letter of a string").
			Arg("txt", "string", "the text to capitalize").
			Return("string", "the text with capitalized first letter").
			Example("capitalize(text)").
			Mode(core.NONE).
			Bypass().
			Do(luatools.StrInStrOutFunc(func(s string) (string, error) {
				return strings.ToUpper(s[:1]) + s[1:], nil
			}))

		return nil
	})
}
