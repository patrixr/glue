package modules

import (
	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/runtime"
)

func init() {
	// @auteur("Modules/Sh")
	//
	// # Sh
	//
	// This module allows running a shell command as part of a blueprint.
	// Although we generally recommend using glue-native modules to achieve tasks, this is helpful
	// when the right tooling is not available.
	//
	// Example:
	//
	// ```lua
	// Sh("ls -la")
	// ```
	//
	Registry.RegisterModule(func(glue *core.Glue) error {
		glue.Plug("sh", core.MODULE).
			Brief("Run a shell command").
			Arg("cmd", STRING, "the shell command to run").
			Do(func(R Runtime, args *Arguments) (RTValue, error) {
				if !glue.Verbose {
					glue.Log.Quiet()
				}

				defer glue.Log.Loud()

				cmd := args.EnsureString(0).String()

				return nil, glue.Machine.Shell(cmd, glue.Log.Stdout, glue.Log.Stderr)
			})

		return nil
	})
}
