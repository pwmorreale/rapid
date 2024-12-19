//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package reporter implements reporting of a sequence run
package reporter

import (
	"io"

	"github.com/pwmorreale/rapid/internal/scenario"
)

// Report defines interfaces for executing scenarios
//
//go:generate counterfeiter -o ../../test/mocks/fake_reporter.go . Report
type Report interface {
	Generate(io.Writer) error
}

// Context defines a sequence
type Context struct {
	sc scenario.Scenario
}

// New creates a new context instance
func New(sc scenario.Scenario) *Context {
	return &Context{
		sc: sc,
	}
}

// Generate creates and sends the report to the specified writer
func (ctx *Context) Generate(_ io.Writer) error {
	return nil
}
