package modules

import (
	"strings"

	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug("capitalize", core.FUNCTION).
			Brief("Uppercase the first letter of a string").
			Arg("txt", STRING, "the text to capitalize").
			Return(STRING, "the text with capitalized first letter").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				s := args.EnsureString(0).String()
				return R.String(strings.ToUpper(s[:1]) + s[1:]), nil
			})

		return nil
	})
}
