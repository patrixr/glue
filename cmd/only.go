/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	. "github.com/patrixr/glue/pkg/runner"
	"github.com/spf13/cobra"
)

var onlyCmd = &cobra.Command{
	Use:   "only",
	Short: "Run Glue on a single part of the configuration using a selector",
	Long:  `Run Glue on a single part of the configuration using a selector`,
	Run: func(cmd *cobra.Command, args []string) {
		planOnly, _ := cmd.Flags().GetBool("plan")
		verbose, _ := cmd.Flags().GetBool("verbose")
		path, _ := cmd.Flags().GetString("path")

		RunGlue(RunOptions{
			PlanOnly: planOnly,
			Verbose:  verbose,
			Path:     path,
			Selector: args[0],
		})
	},
}

func init() {
	onlyCmd.Flags().String("path", "", "Directory or file to look for glue.lua")
	onlyCmd.Flags().BoolP("verbose", "v", false, "Enable verbose mode")
	onlyCmd.Flags().Bool("plan", false, "See the execution blueprints without applying anything")

	rootCmd.AddCommand(onlyCmd)
}
