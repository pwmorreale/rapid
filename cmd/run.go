//
// Copyright © 2025 Peter W. Morreale
//

// Package cmd defines the commands.
package cmd

import (
	"os"

	"github.com/pwmorreale/rapid/internal/logger"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the scenario",
	Long:  `The run command executes the specified scenario file.`,

	RunE: RunScenario,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// RunScenario executes the scenario.
func RunScenario(_ *cobra.Command, _ []string) error {
	opts := logger.Options{
		Writer:    os.Stdout,
		Handler:   logFormat,
		Level:     logLevel,
		Timestamp: logTimestamp,
	}

	if logFilename != "" {
		file, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		opts.Writer = file
	}

	err := logger.Init(&opts)
	if err != nil {
		return err
	}

	return nil
}
