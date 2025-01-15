//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package cmd defines the Cobra root.
package cmd

import (
	"os"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/data"
	"github.com/pwmorreale/rapid/internal/reporter"
	"github.com/pwmorreale/rapid/internal/sequences"
	"github.com/pwmorreale/rapid/internal/service"
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

	c := config.New()

	scenario, err := c.ParseFile(scenarioFile)
	if err != nil {
		return err
	}

	// precompile all regex's to find typos fast...
	err = c.CompileExpressions(scenario)
	if err != nil {
		return err
	}

	d := data.New()
	for k, v := range scenario.Data {
		err = d.AddReplacement(k, v)
		if err != nil {
			return err
		}
	}

	rpt := reporter.New()

	srv := service.New(d)

	seq := sequences.New(srv, rpt)

	// Run the sequence...
	err = seq.Run(scenario)
	if err != nil {
		return err
	}

	return rpt.Generate(scenario, os.Stdout)
}
