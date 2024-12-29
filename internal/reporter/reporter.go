//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package reporter implements reporting of a sequence run
package reporter

import (
	"io"

	"github.com/pwmorreale/rapid/internal/config"
)

// Report defines interfaces for executing scenarios
//
//go:generate counterfeiter -o ../../test/mocks/fake_reporter.go . Report
type Report interface {
	Generate(*config.Scenario, io.Writer) error
}

// Context defines a sequence
type Context struct {
}

// New creates a new context instance
func New() *Context {
	return &Context{}
}

// Generate creates and sends the report to the specified writer
func (ctx *Context) Generate(_ *config.Scenario, _ io.Writer) error {
	return nil
}
