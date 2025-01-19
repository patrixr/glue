package docs

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/patrixr/glue/pkg/blueprint"
	"github.com/patrixr/glue/pkg/core"
	luascaffold "github.com/patrixr/glue/pkg/scaffold/lua"
)

//go:embed templates/*
var tmplFS embed.FS
var templates = template.Must(
	template.New("").Funcs(funcMap).ParseFS(tmplFS, "templates/*.tmpl"),
)

func PrintBlueprintDetails(blueprint blueprint.Blueprint) string {
	return blueprint.PrettyPrint()
}

func PrintMarkdownDocumentation(glue *core.Glue) string {
	var builder strings.Builder

	builder.WriteString("# Glue modules\n\n")
	builder.WriteString("The following Lua modules are available in Glue:\n")

	for _, mod := range glue.Modules {
		builder.WriteString(fmt.Sprintf("- `%s`: %s\n", mod.Name, mod.Brief))
	}

	doc := builder.String()
	prettified, err := glamour.Render(doc, "dark")

	if err != nil {
		return doc
	}

	return prettified
}

func PrintLuaDocumentation(glue *core.Glue) string {
	return luascaffold.NewLuaScaffold(glue).Typegen()
}

func PrintResultReport(glue *core.Glue, results blueprint.Results) string {
	tests := glue.TestResults()
	testPassCount := 0

	for _, test := range tests {
		if test.Error == nil {
			testPassCount++
		}
	}

	testFailCount := len(tests) - testPassCount

	var buf bytes.Buffer

	err := templates.ExecuteTemplate(&buf, "report.md.tmpl", struct {
		Time              string
		Traces            []blueprint.Trace
		TraceCount        int
		Success           bool
		ErrorCount        int
		IncludeTests      bool
		TestResults       []core.TestResult
		TestLen           int
		TestPassCount     int
		TestFailCount     int
		TestSkipCount     int
		SystemIsCompliant bool
	}{
		Time:              time.Now().Format(time.RFC822),
		Traces:            results.Traces,
		TraceCount:        len(results.Traces),
		Success:           results.Success,
		ErrorCount:        results.ErrorCount,
		IncludeTests:      len(tests) > 0,
		TestResults:       tests,
		TestLen:           len(tests),
		TestPassCount:     testPassCount,
		TestFailCount:     testFailCount,
		TestSkipCount:     0,
		SystemIsCompliant: results.Success && results.ErrorCount == 0 && testFailCount == 0,
	})

	if err != nil {
		return err.Error()
	}

	markdown := buf.String()

	prettified, err := glamour.Render(markdown, "auto")

	if err != nil {
		return markdown
	}

	return prettified
}
