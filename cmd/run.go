//
// Copyright © 2025 Peter W. Morreale
//

// Package cmd defines the commands.
package cmd

import (
	"os"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/data"
	"github.com/pwmorreale/rapid/internal/logger"
	"github.com/pwmorreale/rapid/internal/rest"
	"github.com/pwmorreale/rapid/internal/sequence"
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

func initLogger() (*os.File, error) {

	var file os.File

	opts := logger.Options{
		Writer:    os.Stdout,
		Handler:   logFormat,
		Level:     logLevel,
		Timestamp: logTimestamp,
	}

	if logFilename != "" {
		file, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}
		opts.Writer = file
	}

	err := logger.Init(&opts)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &file, nil
}

func initData(sc *config.Scenario) (data.Data, error) {

	var err error

	d := data.New()
	for i := 0; i < len(sc.Replacements); i++ {
		r := sc.Replacements[i]
		err = d.AddReplacement(r.Regex, r.Value)
		if err != nil {
			break
		}
	}

	return d, err
}

// RunScenario executes the scenario.
func RunScenario(_ *cobra.Command, _ []string) error {

	file, err := initLogger()
	if err != nil {
		return nil
	}
	defer (*file).Close()

	c := config.New()
	sc, err := c.ParseFile(scenarioFile)
	if err != nil {
		return err
	}

	d, err := initData(sc)
	if err != nil {
		return err
	}

	r := rest.New(sc, d)
	s := sequence.New(r)

	return s.Run(sc)
}
