/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/docs"
	"github.com/patrixr/glue/pkg/luatools"
	"github.com/patrixr/glue/pkg/modules"
	"github.com/patrixr/glue/pkg/shell"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes glue on your system",
	Long: `Initializes glue in the ~/.config folder of your system
by creating a script file for you to configure your system
`,
	Run: func(cmd *cobra.Command, args []string) {
		var glueFolder string
		var err error

		glue := core.NewGlueWithOptions(core.GlueOptions{
			DryRun: true,
		})

		// Helpers
		boom := func(message string) {
			glue.Log.Error("Unable to initialize glue in folder")
			glue.Log.Error(message)
			os.Exit(1)
		}

		mkdirp := func(path string) {
			cmd := fmt.Sprintf("mkdir -p %s", path)

			if err := shell.Run(cmd, os.Stdout, os.Stderr); err != nil {
				boom(err.Error())
			}
		}

		assert := func(err error, message string) {
			if err != nil {
				glue.Log.Error(message)
				boom(err.Error())
			}
		}

		// Install modules for documentation
		assert(modules.Registry.InstallModules(glue), "Failed to initialize glue modules")

		// Check glue home location
		if len(args) > 0 {
			glueFolder = args[0]
			stat, err := os.Stat(glueFolder)

			assert(err, "Failed to read glue folder")

			if stat.IsDir() == false {
				boom("Expected a directory as argument")
			}
		} else {
			glueFolder, err = core.DefaultGlueFolder()

			if err != nil {
				boom("Unable to initialize glue in folder")
			}
		}

		glueScript := filepath.Join(glueFolder, "glue.lua")

		// Initialize glue script
		if _, err := core.TryFindGlueFile(glueFolder); err == nil {
			glue.Log.Info("A glue script already exists at ~/.config. Skipping")
		} else {
			mkdirp(glueFolder)
			assert(os.WriteFile(glueScript, []byte("-- Input code below"), 0644), "Failed to created glue file")
			glue.Log.Info("Glue script initialized at " + glueScript)
		}

		// Create lib with metadata
		luarcFile := filepath.Join(glueFolder, ".luarc.json")
		libFolder := filepath.Join(glueFolder, "lib")
		libFile := filepath.Join(libFolder, "glue_lib.lua")

		mkdirp(libFolder)

		assert(os.WriteFile(libFile, []byte(docs.GenerateLuaDocumentation(glue)), 0644), "Failed to create glue meta file")

		// Create luarc
		luarc, err := luatools.InitLuaRC(glueFolder)

		assert(err, "Failed to initialize .luarc")

		luarc.AddLibrary("./lib")
		json, err := luarc.ToJSON()

		assert(err, "Failed to marshal luarc into JSON")

		assert(os.WriteFile(luarcFile, []byte(json), 0644), "Failed to save luarc file")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
