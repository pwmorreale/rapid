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

// Loop defines the max iterations and a time limit to complete.
type Loop struct {
	Iterations int       `mapstructure:"iterations"`
	Duration   time.Time `mapstructure:"time_limit"`
}

// Context defines a sequence
type Context struct {
	sc  scenario.Scenario
	rpt reporter.Report
}

// New creates a new context instance
func New(sc scenario.Scenario, rpt reporter.Report) *Context {
	return &Context{
		sc:  sc,
		rpt: rpt,
	}
}

// Run executes the sequence.
func (ctx *Context) Run() error {
	return nil
}
