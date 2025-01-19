package runner

import (
	"fmt"
	"os"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/docs"
)

type RunOptions struct {
	Verbose  bool
	PlanOnly bool
	Path     string
	Selector string
}

func RunGlue(opts RunOptions) {
	glue := InitializeGlue(core.GlueOptions{
		Selector: opts.Selector,
		Verbose:  opts.Verbose,
	})

	defer glue.Close()

	if len(opts.Selector) > 0 && !core.ValidSelectorString(opts.Selector) {
		glue.Log.Error(fmt.Sprintf("Invalid selector '%s'. Bad format or invalid characters", opts.Selector))
		os.Exit(1)
	}

	var script string
	var err error
	if opts.Path != "" {
		script, err = core.TryFindGlueFile(opts.Path)
	} else {
		script, err = core.AutoDetectScriptFile()
	}

	if err != nil {
		glue.Log.Error(err)
		os.Exit(1)
	}

	plan, err := glue.CompilePlan(script)

	if err != nil {
		glue.Log.Error(err)
		os.Exit(1)
	}

	if opts.Verbose {
		fmt.Println(docs.PrintBlueprintDetails(plan))
	}

	if opts.PlanOnly {
		return
	}

	results := plan.Execute()

	glue.Test()

	fmt.Println(docs.PrintResultReport(glue, results))

	if !results.Success {
		os.Exit(1)
	}
}
