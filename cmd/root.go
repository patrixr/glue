/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	. "github.com/patrixr/glue/pkg/runner"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "glue",
	Short: "Machine configuration tool",
	Long:  `Glue is a machine configuration tool that allows you to use Lua to easily streamline your system setup`,
	Run: func(cmd *cobra.Command, args []string) {
		planOnly, _ := cmd.Flags().GetBool("plan")
		verbose, _ := cmd.Flags().GetBool("verbose")
		path, _ := cmd.Flags().GetString("path")

		RunGlue(RunOptions{
			PlanOnly: planOnly,
			Verbose:  verbose,
			Path:     path,
		})
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("path", "", "Directory or file to look for glue.lua")
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose mode")
	rootCmd.Flags().Bool("plan", false, "See the execution blueprints without applying anything")
}
