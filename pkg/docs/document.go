package docs

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/patrixr/glue/pkg/core"
)

//go:embed templates/*
var tmplFS embed.FS

var templates = template.Must(
	template.New("").ParseFS(tmplFS, "*.md"),
)

func GenerateModuleDoc(mod *core.GlueModule) string {
	var buf bytes.Buffer
	templates.ExecuteTemplate(&buf, "module.md", mod)
	return buf.String()
}
