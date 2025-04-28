//
// Copyright © 2025 Peter W. Morreale
//

// Package cmd contains the commands
package cmd

import (
	"os"

	"github.com/pwmorreale/rapid/internal/logger"
	"github.com/pwmorreale/rapid/internal/verify"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a scenario configuration",
	Long:  `Verify syntax and elements of a yaml senario configuration`,
	RunE:  DoVerify,
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}

// StartVerify starts the verify command.
func DoVerify(_ *cobra.Command, _ []string) error {

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

	return verify.Check(scenarioFile)
}
