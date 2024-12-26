//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package sequences defines a sequence of RAPID operations
package sequences

import (
	"time"

	"github.com/pwmorreale/rapid/internal/reporter"
	"github.com/pwmorreale/rapid/internal/scenario"
)

// Sequence defines interfaces for executing scenarios
//
//go:generate counterfeiter -o ../../test/mocks/fake_sequence.go . Sequence
type Sequence interface {
	Run() error
}

// Context defines a sequence
type Context struct {
	sc  scenario.Scenario
	rpt reporter.Report
}

// Seqs contains the sequence configuration.
type Seqs struct {
	Iterations int           `mapstructure:"iterations"`
	Duration   time.Duration `mapstructure:"time_limit"`
	ErrorAbort bool          `mapstructure:"abort_on_error"`
	IgnoreDups bool          `mapstructure:"ignore_duplicate_errors"`
	Reqs       []Request     `mapstructure:"requests"`
}

// ContentData cdefines expected response data.
type ContentData struct {
	Type     string   `mapstructure:"type"`
	Contains []string `mapstructure:"contains"`
	Extract  string   `mapstructure:"extract"`
	SaveAs   string   `mapstructure:"save_as"`
}

// Response defines a REST reqponse
type Response struct {
	Content ContentData `mapstructure:"content"`
	Status  []int       `mapstructure:"status"`
}

// Request defines the a request/response
type Request struct {
	Name         string            `mapstructure:"name"`
	Scheme       string            `mapstructure:"scheme"`
	Path         string            `mapstructure:"path"`
	Host         string            `mapstructure:"host"`
	Query        map[string]string `mapstructure:"query"`
	ExtraHeaders map[string]string `mapstructure:"extra_headers"`
	Content      string            `mapstructure:"content"`
	ContentType  string            `mapstructure:"content_type"`
	TimeLimit    time.Duration     `mapstructure:"time_limit"`
	Rsp          Response          `mapstructure:"response"`
}

// New creates a new context instance
func New(sc scenario.Scenario, rpt reporter.Report) *Context {
	return &Context{
		sc:  sc,
		rpt: rpt,
	}
}

func (ctx *Context) UnmarshalKey(key string) (*Seqs, error) {

	var s Seqs

	err := ctx.sc.Viper().UnmarshalKey(key, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Run executes the sequence.
func (ctx *Context) Run() error {

	return nil
}
