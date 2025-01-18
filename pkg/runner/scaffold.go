package runner

import (
	"os"

	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/scaffold"
)

type ScaffoldOptions struct {
	Path string
}

func RunGlueScaffold(glueFolder string, scaffold Scaffold) {
	glue := InitializeGlue(core.GlueOptions{})

	defer glue.Close()

	boom := func(message string) {
		glue.Log.Error("Unable to initialize glue in folder")
		glue.Log.Error(message)
		os.Exit(1)
	}

	assert := func(err error, message string) {
		if err != nil {
			glue.Log.Error(message)
			boom(err.Error())
		}
	}

	if len(glueFolder) == 0 {
		home, err := core.GlueHome()
		if err != nil {
			boom("Unable to initialize glue in folder")
		}
		glueFolder = home
	}

	stat, err := os.Stat(glueFolder)

	if err != nil && !os.IsNotExist(err) {
		boom("Failed to read glue folder: " + err.Error())
	}

	if err == nil && !stat.IsDir() {
		boom("Scaffold target (" + glueFolder + ") is not a directory")
	}

	assert(scaffold.Setup(glueFolder), "Failed to scaffold glue")
}
