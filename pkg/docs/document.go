package docs

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/charmbracelet/glamour"
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

func GenerateMarkdownDocumentation(glue *core.Glue) string {
	var builder strings.Builder

	builder.WriteString("# Glue modules\n\n")
	builder.WriteString("The following Lua modules are available in Glue:\n")

	for _, mod := range glue.Modules {
		builder.WriteString(fmt.Sprintf("- `%s`: %s\n", mod.Name, mod.Short))
	}

	doc := builder.String()
	prettified, err := glamour.Render(doc, "dark")

	if err != nil {
		return doc
	}

	return prettified
}

func GenerateLuaDocumentation(glue *core.Glue) string {
	var builder strings.Builder

	builder.WriteString(glue.Annotations.Render())
	builder.WriteString("\n")

	return builder.String()
}
