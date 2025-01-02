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

// Run executes the sequence.
func (ctx *Context) Run(sc *config.Scenario) error {

	startTime := time.Now()

	for i := 0; i < sc.Seq.Iterations; i++ {

		for _, r := range sc.Seq.Reqs {

			req, err := ctx.service.Create(&r)
			if err != nil {
				return err
			}

			rsp, err := ctx.service.Send(req, &r)
			if err != nil {
				return err
			}

			err = ctx.service.Validate(rsp, &r)
			if err != nil {
				return err
			}

			time.Sleep(sc.Seq.Delay)

			if sc.Seq.Limit > 0 && time.Since(startTime) > sc.Seq.Limit {
				// Log something.
				return fmt.Errorf("Execution exceeded time limit of: %v", sc.Seq.Limit)
			}
		}

	}

	return nil
}
