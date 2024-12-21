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
	cfgFile      string
	scenarioFile string
	rootCmd      = &cobra.Command{
		Use:   "rapid",
		Short: "rapid is an REST API diagnostic tool.",
		Long:  `rapid enables testing of REST APIs.`,
		RunE:  RunRoot,
	}
)

// rootCmd represents the base command when called without any subcommands

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&scenarioFile, "scenario", "s", "", "Path to scenario file.")
	rootCmd.MarkFlagRequired("scenario")
	rootCmd.MarkFlagFilename("scenario")
}

// RunRoot executes the CLI interface.
func RunRoot(_ *cobra.Command, _ []string) error {

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
