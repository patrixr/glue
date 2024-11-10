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

func run(cmd *cobra.Command, selector string) {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	path, _ := cmd.Flags().GetString("path")

	glue := core.NewGlueWithOptions(core.GlueOptions{
		DryRun:   dryRun,
		Selector: selector,
	})

	defer glue.Close()

	if len(selector) > 0 && !core.ValidSelectorString(selector) {
		glue.Log.Error(fmt.Sprintf("Invalid selector '%s'. Bad format or invalid characters", selector))
		os.Exit(1)
	}

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
}

var rootCmd = &cobra.Command{
	Use:   "glue",
	Short: "Machine configuration tool",
	Long:  `Glue is a machine configuration tool that allows you to use Lua to easily streamline your system setup`,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, "")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	xCmd := &cobra.Command{
		Use:   "x",
		Short: "Execute Glue with a selector",
		Long: `Glue allows nested groups of execution blocks.
The x command allows us to only run a subset of these groups using a selector argument`,
		Args: cobra.MatchAll(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args[0])
		},
	}

	cmds := []*cobra.Command{
		xCmd,
		rootCmd,
	}

	for _, cmd := range cmds {
		cmd.Flags().StringP("path", "p", "", "Directory or file to look for glue.lua")
		cmd.Flags().BoolP(
			"dry-run", "d", false, "See the execution flow without running anything")
	}

	rootCmd.AddCommand(xCmd)
}
