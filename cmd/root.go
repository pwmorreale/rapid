//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package cmd defines the Cobra root.
package cmd

import (
	"github.com/pwmorreale/rapid/internal/scenario"
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
	rootCmd.AddCommand(serverCmd)
}

// RunRoot executes the CLI interface.
func RunRoot(cmd *cobra.Command, args []string) error {

	sc := scenario.New()

	err := sc.ReadInConfig(scenarioFile)
	if err != nil {
		return err
	}

	return sc.Execute()
}
