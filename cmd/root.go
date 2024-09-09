/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/modules"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "glue",
	Short: "Machine configuration tool",
	Long:  `Glue is a machine configuration tool that allows you to use Lua to easily streamline your system setup`,
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		unsafe, _ := cmd.Flags().GetBool("unsafe")
		path, _ := cmd.Flags().GetString("path")

		if dryRun && unsafe {
			fmt.Println("The options --dry-run and --unsafe cannot both be set")
			os.Exit(1)
		}

		glue := core.NewGlueWithOptions(core.GlueOptions{
			DryRun: dryRun,
			Unsafe: unsafe,
		})

		defer glue.Close()

		var script string
		var err error
		if path != "" {
			script, err = core.TryFindGlueFile(path)
		} else {
			script, err = core.AutoDetectScriptFile()
		}

		if err != nil {
			glue.Log.Error(err)
			os.Exit(1)
		}

		if err := modules.Registry.InstallModules(glue); err != nil {
			glue.Log.Error(err)
			os.Exit(1)
		}

		if err := glue.Execute(script); err != nil {
			glue.Log.Error(err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("path", "p", "", "Directory or file to look for glue.lua")
	rootCmd.Flags().BoolP(
		"dry-run", "d", false, "See the execution flow without running anything")
	rootCmd.Flags().BoolP(
		"unsafe", "u", false, "When enabled, allows native lua libraries to be loaded")
}
