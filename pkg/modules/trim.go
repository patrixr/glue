package modules

import (
	"strings"

	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
	"github.com/patrixr/q"
)

func init() {
	Registry.RegisterModule(
		func(glue *core.Glue) error {

			glue.Plug().
				Name("trim").
				Short("Trims the extra indentation of a multi-line string").
				Long("Trims the extra indentation of a multi-line string").
				Arg("txt", STRING, "the text to trim").
				Return("string", "the trimmed text").
				Example("trim(text)").
				Mode(core.NONE).
				Bypass().
				Do(func(R Runtime, args *Arguments) (RTValue, error) {
					s := args.EnsureString(0).String()
					return R.String(strings.TrimSpace(q.TrimIndent(s))), nil
				})

			return nil
		})
}
