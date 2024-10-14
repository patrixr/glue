package docs

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/patrixr/glue/pkg/core"
)

//go:embed templates/*
var tmplFS embed.FS

var templates = template.Must(
	template.ParseFS(tmplFS, "templates/*.md"),
)

func GenerateModuleDoc(mod *core.GlueModule) string {
	var buf bytes.Buffer
	templates.ExecuteTemplate(&buf, "module.md", mod)
	return buf.String()
}

func GenerateCLIDocumentation(glue *core.Glue) string {
	var builder strings.Builder

	builder.WriteString("# Glue modules\n\n")
	builder.WriteString("The following Lua modules are available in Glue:\n")

	for _, mod := range glue.Modules {
		builder.WriteString(fmt.Sprintf("- `%s`: %s\n", mod.Name, mod.Short))
	}

	return builder.String()
}
