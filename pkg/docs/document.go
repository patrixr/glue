package docs

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/patrixr/glue/pkg/core"
)

//go:embed templates/*
var tmplFS embed.FS
var templates = template.Must(
	template.New("").Funcs(funcMap).ParseFS(tmplFS, "templates/*.tmpl"),
)

/*
@TODO: Implement module documentation generation

func GenerateModuleDoc(mod *core.GlueModule) string {
	var buf bytes.Buffer
	templates.ExecuteTemplate(&buf, "module.md", mod)
	return buf.String()
}
*/

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

func GenerateResultReport(glue *core.Glue) string {
	success, errorCount, traces := glue.Result()

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
		Traces            []core.Trace
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
		Traces:            traces,
		TraceCount:        len(traces),
		Success:           success,
		ErrorCount:        errorCount,
		IncludeTests:      glue.RunTests && len(tests) > 0,
		TestResults:       tests,
		TestLen:           len(tests),
		TestPassCount:     testPassCount,
		TestFailCount:     testFailCount,
		TestSkipCount:     0,
		SystemIsCompliant: success && errorCount == 0 && testFailCount == 0,
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
