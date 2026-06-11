//
// Copyright © 2025 Peter W. Morreale
//

// Package cmd defines the commands.
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/data"
	"github.com/pwmorreale/rapid/logger"
	"github.com/pwmorreale/rapid/report"
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

var reportFile string
var dumpFile string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&reportFile, "report", "", `Write structured report to file (format from extension: .json or .xml)`)

	runCmd.Flags().StringVar(&dumpFile, "dump", "", `Dump raw HTTP request/response traffic (to stdout if no file specified)`)
	runCmd.Flags().Lookup("dump").NoOptDefVal = "stdout"
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
		err = d.AddReplacement(r.Regex, os.ExpandEnv(r.Value))
		if err != nil {
			break
		}
	}

	return d, err
}

// RunScenario executes the scenario.
func RunScenario(_ *cobra.Command, _ []string) error {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
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

	if sc.Name != "" {
		logger.Info(nil, nil, "scenario: %s version: %s %s", sc.Name, sc.Version, sc.Comment)
	}

	dumpWriter, dumpCloser, err := initDump()
	if err != nil {
		return err
	}
	if dumpCloser != nil {
		defer dumpCloser.Close()
	}

	r := rest.New(sc, d, dumpWriter)
	s := sequence.New(r)

	err = s.Run(ctx, sc)
	if err != nil {
		return err
	}

	if err := r.Push(); err != nil {
		return err
	}

	LogResults(sc)

	if reportFile != "" {
		if err := writeReport(reportFile, sc); err != nil {
			return err
		}
	}

	if totalErrors(sc) > 0 {
		return fmt.Errorf("scenario completed with errors")
	}

	return nil
}

func writeReport(path string, sc *config.Scenario) error {
	switch filepath.Ext(path) {
	case ".json":
		return report.WriteJSON(path, sc)
	case ".xml":
		return report.WriteJUnit(path, sc)
	default:
		return fmt.Errorf("unsupported report format %q (use .json or .xml)", filepath.Ext(path))
	}
}

func initDump() (io.Writer, io.Closer, error) {
	if dumpFile == "" {
		return nil, nil, nil
	}
	if dumpFile == "stdout" {
		return os.Stdout, nil, nil
	}
	f, err := os.OpenFile(dumpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, nil, err
	}
	return f, f, nil
}

func totalErrors(sc *config.Scenario) int64 {
	var total int64
	for i := range sc.Sequence.Requests {
		total += sc.Sequence.Requests[i].Stats.GetErrors()
	}
	return total
}

// LogResults prints out the statistics from the run.
func LogResults(sc *config.Scenario) {

	for i := range sc.Sequence.Requests {
		request := &sc.Sequence.Requests[i]

		str := request.Stats.String()
		logger.Info(request, nil, "%s", str)

		for j := range request.Responses {
			response := request.Responses[j]

			str := response.Stats.String()
			logger.Info(request, response, "%s", str)
		}

		for j := range request.UnknownResponses {
			response := request.UnknownResponses[j]

			str := response.Stats.String()
			logger.Warn(request, response, "%s", str)
		}
	}
}
