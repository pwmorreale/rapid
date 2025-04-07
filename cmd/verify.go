//
// Copyright © 2025 Peter W. Morreale
//

// Package cmd contains the commands
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a scenario configuration",
	Long:  `Verify syntax and elements of a yaml senario configuration`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("verify called", scenarioFile)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
