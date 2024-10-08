/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/patrixr/glue/pkg/core"
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
		homedir := os.Getenv("HOME")

		if len(homedir) == 0 {
			fmt.Println("Error initializing glue. $HOME env variable not found")
			os.Exit(1)
		}

		configFolder := filepath.Join(homedir, ".config")
		glueFolder := filepath.Join(configFolder, "glue")
		glueScript := filepath.Join(glueFolder, "glue.lua")

		if _, err := core.TryFindGlueFile(glueFolder); err == nil {
			fmt.Println("A glue script already exists at ~/.config\nSkipping")
			os.Exit(0)
		}

		fmt.Println("Initializing glue...")

		mkdirp := fmt.Sprintf("mkdir -p %s", glueFolder)

		if err := shell.Run(mkdirp, os.Stdout, os.Stderr); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		touch := fmt.Sprintf("touch %s", glueScript)

		if err := shell.Run(touch, os.Stdout, os.Stderr); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Glue script initialized at " + glueScript)
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
