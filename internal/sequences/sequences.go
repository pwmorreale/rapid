//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package sequences defines a sequence of RAPID operations
package sequences

import (
	"fmt"
	"time"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/reporter"
	"github.com/pwmorreale/rapid/internal/service"
)

// Sequence defines interfaces for executing scenarios
//
//go:generate counterfeiter -o ../../test/mocks/fake_sequence.go . Sequence
type Sequence interface {
	Run(*config.Scenario) error
}

// Context defines a sequence
type Context struct {
	reporter reporter.Report
	service  service.Service
}

// New creates a new context instance
func New(srv service.Service, rpt reporter.Report) *Context {
	return &Context{
		reporter: rpt,
		service:  srv,
	}
}

func (ctx *Context) handleRequest(r *config.Request) error {

	request, err := ctx.service.CreateRequest(r)
	if err != nil {
		return err
	}

	client, err := ctx.service.CreateClient(r)
	if err != nil {
		return err
	}

	response, err := ctx.service.Send(client, request, r)
	if err != nil {
		return err
	}

	err = ctx.service.ValidateResponse(client, response, r)
	if err != nil {
		return err
	}

	return nil
}

// Run executes the sequence.
func (ctx *Context) Run(sc *config.Scenario) error {

	startTime := time.Now()

	for i := 0; i < sc.Sequence.Iterations; i++ {

		for _, r := range sc.Sequence.Requests {

			err := ctx.handleRequest(&r)
			if err != nil && sc.Sequence.AbortOnError {
				return err
			}

			time.Sleep(sc.Sequence.Delay)

			if sc.Sequence.Limit > 0 && time.Since(startTime) > sc.Sequence.Limit {
				// Log something.
				return fmt.Errorf("execution exceeded time limit of: %v", sc.Sequence.Limit)
			}
		}

	}

	return nil
}
