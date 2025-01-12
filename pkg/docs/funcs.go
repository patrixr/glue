package docs

import (
	"strconv"
	"text/template"

	"github.com/patrixr/q"
)

var funcMap = template.FuncMap{
	"add": func(a, b int) string {
		return strconv.Itoa(a + b)
	},
	"ellipsis": func(txt string) string {
		return q.Ellipsis(txt, 25)
	},
	"len": func(array []any) int {
		return len(array)
	},
	"errorstr": func(e error) string {
		return e.Error()
	},
	"gt": func(a, b int) bool {
		return a > b
	},
}
