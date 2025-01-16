package modules

import (
	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
	"github.com/patrixr/glue/pkg/shell"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug().
			Name("sh").
			Short("Run a shell command").
			Long("Run a shell command").
			Arg("cmd", STRING, "the shell command to run").
			Example("sh('ls')").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				if !glue.Verbose {
					glue.Log.Quiet()
				}

				defer glue.Log.Loud()

				cmd := args.EnsureString(0).String()

				return nil, shell.Run(cmd, glue.Log.Stdout, glue.Log.Stderr)
			})

		return nil
	})
}
