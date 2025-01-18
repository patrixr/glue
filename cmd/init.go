/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/runner"
	luascaffold "github.com/patrixr/glue/pkg/scaffold/lua"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes glue on your system",
	Long:  `Initializes glue in the ~/.config folder of your system by creating a script file for you to configure your system`,
	Run: func(cmd *cobra.Command, args []string) {
		targetFolder := ""

		// Check glue home location
		if len(args) > 0 {
			targetFolder = args[0]
		}

		runner.RunGlueScaffold(targetFolder, luascaffold.NewLuaScaffold(
			runner.InitializeGlue(core.GlueOptions{}),
		))
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
