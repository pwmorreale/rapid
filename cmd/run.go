//
// Copyright © 2025 Peter W. Morreale
//

// Package cmd defines the commands.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the scenario",
	Long:  `The run command executes the specified scenario file.`,

	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("run called", scenarioFile)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
