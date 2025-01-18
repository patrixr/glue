package runner

import (
	"os"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/modules"
)

func InitializeGlue(opts core.GlueOptions) *core.Glue {
	glue := core.NewGlueWithOptions(opts)

	if err := modules.Registry.InstallModules(glue); err != nil {
		glue.Log.Error(err)
		os.Exit(1)
	}

	return glue
}
