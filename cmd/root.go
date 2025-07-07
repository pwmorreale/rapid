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
var logFormat string
var logFilename string
var logLevel string
var logTimestamp bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "rapid",
	Short:   "REST API diagnostic tool",
	Long:    `Rapid can be used to verify conformance of your REST APIS to their design specs.  You can also use it to measure performance and/or throughput of your REST servers.`,
	Version: "v0.1.0",
}

// Execute is called by main.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&logFormat, "log_format", "f", "text", `Log format: "text" or "json" `)
	rootCmd.PersistentFlags().BoolVarP(&logTimestamp, "log_timestamp", "t", false, `Add timeStamp to log entries `)
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log_level", "l", "info", `Log level, one of: "debug", "info", "warn", or "error"`)
	rootCmd.PersistentFlags().StringVarP(&logFilename, "log_file", "", "", `Log to "filename" instead of stdout`)

	rootCmd.PersistentFlags().StringVarP(&scenarioFile, "scenario", "s", "", "Path to scenario file")
	rootCmd.MarkPersistentFlagRequired("scenario")
	rootCmd.MarkPersistentFlagFilename("scenario", "yaml")
}
