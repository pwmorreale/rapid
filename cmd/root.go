//
// Copyright © 2025 Peter W. Morreale
//

// Package cmd contains the commands for the application
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var scenarioFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rapid",
	Short: "REST API diagnostic tool",
	Long:  `Rapid can be used to verify conformance of your REST APIS to their design specs.  You can also use it to measure performance and/or throughput of your REST servers.`,
}

// Execute is called by main.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&scenarioFile, "scenario", "s", "", "Path to scenario file")
	rootCmd.MarkFlagRequired("scenario")
	rootCmd.MarkFlagFilename("scenario")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
