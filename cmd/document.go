/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/patrixr/glue/pkg/core"
	"github.com/patrixr/glue/pkg/docs"
	"github.com/patrixr/glue/pkg/modules"
	"github.com/spf13/cobra"
)

// documentCmd represents the document command
var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Generates a documentation of Glue's internal functions",
	Long:  `Generates a documentation of Glue's internal functions`,
	Run: func(cmd *cobra.Command, args []string) {
		glue := core.NewGlueWithOptions(core.GlueOptions{
			DryRun: true,
		})

		if err := modules.Registry.InstallModules(glue); err != nil {
			glue.Log.Error(err)
			os.Exit(1)
		}

		fmt.Println(docs.GenerateCLIDocumentation(glue))
	},
}

func init() {
	rootCmd.AddCommand(documentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// documentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// documentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}