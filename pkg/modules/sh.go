package modules

import (
	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/luatools"
	"github.com/patrixr/glue/pkg/shell"
)

func init() {
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug().
			Name("sh").
			Short("Run a shell command").
			Long("Run a shell command").
			Arg("cmd", "string", "the shell command to run").
			Example("sh('ls')").
			Do(luatools.StrFunc(func(cmd string) error {
				if !glue.Verbose {
					glue.Log.Quiet()
				}

				defer glue.Log.Loud()

				return shell.Run(cmd, glue.Log.Stdout, glue.Log.Stderr)
			}))

		return nil
	})
}
