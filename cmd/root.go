//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package cmd defines the Cobra root.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	scenarioFile string
	rootCmd      = &cobra.Command{
		Use:   "rapid",
		Short: "rapid is an REST API diagnostic tool.",
		Long:  `rapid enables testing of REST APIs.`,
		RunE:  RunCli,
	}
)

// Start starts the application.
func Start() error {

	rootCmd.Flags().StringVarP(&scenarioFile, "scenario", "s", "", "Path to scenario file.")
	rootCmd.MarkFlagRequired("scenario")
	rootCmd.MarkFlagFilename("scenario")

	rootCmd.AddCommand(sanityCmd)
	rootCmd.AddCommand(runCmd)

	return rootCmd.Execute()
}

// RunCli executes the CLI interface.
func RunCli(_ *cobra.Command, _ []string) error {

	fmt.Printf("running in ROOT")

	return nil
}
