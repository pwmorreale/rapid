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
	Run(*config.Scenario) error
}

// Statistics defines sequence execution statistics.
type Statistics struct {
	Iterations   int
	Calls        int64
	RestCallTime time.Duration
	ElaspedTime  time.Duration
}

// Context defines a sequence
type Context struct {
	rest rest.Rest

	// Stats
	iterations   int
	calls        int64
	restCallTime int64 // So we can use atomics directly...
	elaspedTime  time.Duration
}

// New creates a new context instance
func New(r rest.Rest) *Context {

	return &Context{
		rest: r,
	}
}

// GetStats returns current statistics.
func (s *Context) GetStats() *Statistics {
	return &Statistics{
		Iterations:   s.iterations,
		Calls:        s.calls,
		RestCallTime: time.Duration(s.restCallTime),
		ElaspedTime:  s.elaspedTime,
	}
}

// Run executes the sequence.
func (s *Context) Run(sc *config.Scenario) error {

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), sc.Sequence.Limit)
	defer cancel()

	for i := 0; i < sc.Sequence.Iterations; i++ {

		s.ExecuteSequence(ctx, sc)
		s.iterations++
	}

	s.elaspedTime = time.Since(start)

	err := ctx.Err()
	if err != nil {
		logger.Error(nil, nil, "sequence %v on iteration: %d", err, s.iterations)
	}

	return nil
}

// ExecuteSequence runs the sequence of requests
func (s *Context) ExecuteSequence(ctx context.Context, sc *config.Scenario) {

Loop:
	for i := range sc.Sequence.Requests {

		// Did we timeout?
		select {
		case <-ctx.Done():
			break Loop
		default:
		}

		request := &sc.Sequence.Requests[i]

		logger.Info(request, nil, "execution started")
		s.ExecuteRequest(ctx, request)
		logger.Info(request, nil, "execution complete")

	}

}

// ExecuteRequest ezxecutes a request
func (s *Context) ExecuteRequest(ctx context.Context, request *config.Request) {

	// Is this a once only request?
	if request.OnceOnly {
		if request.Executed {
			logger.Info(request, nil, "once_only request already executed, ignoring")
			return
		}
		request.Executed = true

	}

	// Default to one if not specified...
	workerPoolSize := request.ThunderingHerd.Size
	if workerPoolSize == 0 {
		workerPoolSize++
	}
	wp := workerpool.New(request.ThunderingHerd.Size)

	// Default to one if not specified...
	numRequests := request.ThunderingHerd.Max
	if numRequests == 0 {
		numRequests = 1
	}

Loop:
	for i := 0; i < numRequests; i++ {
		if wp.WaitingQueueSize() > 0 {
			time.Sleep(time.Millisecond * 100)
		}

		wp.Submit(func() {
			atomic.AddInt64(&s.calls, 1)
			start := time.Now()
			s.rest.Execute(ctx, request)
			atomic.AddInt64(&s.restCallTime, int64(time.Since(start)))
		})

		select {
		case <-ctx.Done():
			break Loop
		default:
		}
	}

	// Wait for everybody to complete.
	wp.Stop()

}
