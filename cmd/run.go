//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package cmd defines the Cobra run
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

// runCmd represents the server command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the scenario",
	Long:  "Execute the scenario",
	RunE:  runScenario,
}

// runScenario creates and runs a REST server instance
func runScenario(_ *cobra.Command, _ []string) error {
	c := config.New()

	scenario, err := c.ParseFile(scenarioFile)
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
