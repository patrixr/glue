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

func GenerateMarkdownDocumentation(glue *core.Glue) string {
	var builder strings.Builder

	builder.WriteString("# Glue modules\n\n")
	builder.WriteString("The following Lua modules are available in Glue:\n")

	for _, mod := range glue.Modules {
		builder.WriteString(fmt.Sprintf("- `%s`: %s\n", mod.Name, mod.Short))
	}

	return builder.String()
}

func GenerateLuaDocumentation(glue *core.Glue) string {
	var builder strings.Builder
	// var globals *q.Node[string] = q.Tree[string]("root")

	// builder.WriteString("---@meta\n\n")

	// for _, mod := range glue.Modules {
	// 	nestedKeys := strings.Split(mod.Name, ".")

	// 	if len(nestedKeys) <= 1 {
	// 		continue
	// 	}

	// 	node := globals
	// 	for _, key := range nestedKeys[:len(nestedKeys)] {
	// 		child := node.FindChild(func(s string) bool {
	// 			return s == key
	// 		})

	// 		if child == nil {
	// 			child = node.Add(key)
	// 		}

	// 		node = child
	// 	}
	// }

	// var traverse func(ref *q.Node[string], level int)

	// traverse = func(ref *q.Node[string], level int) {
	// 	builder.WriteString(fmt.Sprintf("%s%s = {\n", strings.Repeat("  ", level), ref.Data))
	// 	for _, child := range ref.Children {
	// 		traverse(child, level+1)
	// 	}
	// 	builder.WriteString(fmt.Sprintf("%s}\n", strings.Repeat("  ", level)))
	// }

	// for _, global := range globals.Children {
	// 	traverse(global, 0)
	// }

	builder.WriteString(glue.Annotations.Render())
	builder.WriteString("\n")

	return builder.String()
}
