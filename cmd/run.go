//
// Copyright © 2025 Peter W. Morreale
//

// Package cmd defines the commands.
package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/data"
	"github.com/pwmorreale/rapid/logger"
	"github.com/pwmorreale/rapid/rest"
	"github.com/pwmorreale/rapid/sequence"
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

	opts := logger.Options{
		Writer:    os.Stdout,
		Handler:   logFormat,
		Level:     logLevel,
		Timestamp: logTimestamp,
	}

	var file *os.File
	if logFilename != "" {
		var err error
		file, err = os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}
		opts.Writer = file
	}

	err := logger.Init(&opts)
	if err != nil {
		if file != nil {
			file.Close()
		}
		return nil, err
	}

	return file, nil
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	file, err := initLogger()
	if err != nil {
		return err
	}
	if file != nil {
		defer file.Close()
	}

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

	err = s.Run(ctx, sc)
	if err != nil {
		return err
	}

	LogResults(sc)

	return nil
}

// LogResults prints out the statistics from the run.
func LogResults(sc *config.Scenario) {

	for i := range sc.Sequence.Requests {
		request := sc.Sequence.Requests[i]

		str := request.Stats.String()
		logger.Info(&request, nil, "%s", str)

		for j := range request.Responses {
			response := request.Responses[j]

			str := response.Stats.String()
			logger.Info(&request, response, "%s", str)
		}
	}
}
