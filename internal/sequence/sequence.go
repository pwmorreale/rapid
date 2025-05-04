//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package sequence defines a sequence of RAPID operations
package sequence

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/logger"
	"github.com/pwmorreale/rapid/internal/rest"
)

// Sequence defines interfaces for executing scenarios
//
//go:generate go tool counterfeiter -o ../../test/mocks/fake_sequence.go . Sequence
type Sequence interface {
	Run(string) error
}

// Context defines a sequence
type Context struct {
	rest rest.Rest

	// Stats
	iterations   int64
	requests     int64
	restCalls    int64
	restCallTime int64 // So we can use atomics directly...
	elaspedTime  time.Duration
}

// New creates a new context instance
func New(r rest.Rest) *Context {

	return &Context{
		rest: r,
	}
}

// Run executes the sequence.
func (s *Context) Run(scenarioFile string) error {

	c := config.New()
	sc, err := c.ParseFile(scenarioFile)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sc.Sequence.Limit))
	defer cancel()

	start := time.Now()

	for i := 0; i < sc.Sequence.Iterations; i++ {

		s.ExecuteSequence(ctx, sc)

		err := ctx.Err()
		if err != nil {
			logger.Error(nil, nil, "sequence %v on iteration: %d", err, i)
			break
		}

		atomic.AddInt64(&s.iterations, 1)

	}

	s.elaspedTime = time.Since(start)

	return nil
}

// ExecuteSequence runs the sequence of requests
func (s *Context) ExecuteSequence(ctx context.Context, sc *config.Scenario) {

	for i := range sc.Sequence.Requests {

		s.ExecuteRequest(ctx, &sc.Sequence.Requests[i])

		// Did we timeout?
		select {
		case <-ctx.Done():
			return
		default:
		}

		atomic.AddInt64(&s.requests, 1)

	}

}

// ExecuteRequest ezxecutes a request
func (s *Context) ExecuteRequest(ctx context.Context, request *config.Request) {

	// Create a timeout context for the thundering herd...
	timeout, cancel := context.WithTimeout(context.Background(), time.Duration(request.ConcurrentDuration))
	defer cancel()

	wp := workerpool.New(request.ConcurrentCalls)

	napTime := time.Millisecond * 100

	for i := 0; i < request.ConcurrentCalls; i++ {

		// Did we timeout?
		select {
		case <-ctx.Done(): // Iteration timeout
			return
		case <-timeout.Done(): // 'Thundering herd' timeout.
			logger.Error(request, nil, "concurrent timeout: %s", timeout.Err())
			return
		default:
		}

		// We only want one (give or take, racy...) waiting task.
		if wp.WaitingQueueSize() == 0 {
			wp.Submit(func() {
				atomic.AddInt64(&s.requests, 1)
				start := time.Now()
				s.rest.Execute(request)
				atomic.AddInt64(&s.restCallTime, int64(time.Since(start)))
			})
		} else {
			time.Sleep(napTime)
		}
	}

	// Wait for everybody to complete.
	wp.Stop()

}
