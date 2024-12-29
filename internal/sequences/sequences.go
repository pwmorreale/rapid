//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package sequences defines a sequence of RAPID operations
package sequences

import (
	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/reporter"
)

// Sequence defines interfaces for executing scenarios
//
//go:generate counterfeiter -o ../../test/mocks/fake_sequence.go . Sequence
type Sequence interface {
	Run(*config.Scenario) error
}

// Context defines a sequence
type Context struct {
	rpt reporter.Report
}

// New creates a new context instance
func New(rpt reporter.Report) *Context {
	return &Context{
		rpt: rpt,
	}
}

// Run executes the sequence.
func (ctx *Context) Run(_ *config.Scenario) error {

	return nil
}
