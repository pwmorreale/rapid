//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package cmd defines the Cobra root.
package cmd

import (
	"os"

	"github.com/pwmorreale/rapid/internal/reporter"
	"github.com/pwmorreale/rapid/internal/scenario"
	"github.com/pwmorreale/rapid/internal/sequences"
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

	return rootCmd.Execute()
}

// RunCli executes the CLI interface.
func RunCli(_ *cobra.Command, _ []string) error {

	sc := scenario.New()

	err := sc.ParseFile(scenarioFile)
	if err != nil {
		return err
	}

	rpt := reporter.New(sc)

	seq := sequences.New(sc, rpt)

	// Run the sequence...
	err = seq.Run()
	if err != nil {
		return err
	}

	return rpt.Generate(os.Stdout)
}
