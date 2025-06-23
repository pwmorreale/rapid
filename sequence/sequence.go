//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package sequence defines a sequence of RAPID operations
package sequence

import (
	"context"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/logger"
	"github.com/pwmorreale/rapid/rest"
)

// Sequence defines interfaces for executing scenarios
//
//go:generate go tool counterfeiter -o ../test/mocks/fake_sequence.go . Sequence
type Sequence interface {
	Run(*config.Scenario) error
}

// Context defines a sequence
type Context struct {
	rest rest.Rest
}

// New creates a new context instance
func New(r rest.Rest) *Context {
	return &Context{
		rest: r,
	}
}

// Run executes the sequence.
func (s *Context) Run(sc *config.Scenario) error {

	for i := 0; i < sc.Sequence.Iterations; i++ {
		s.ExecuteIteration(sc, i)
	}

	return nil
}

// ExecuteIteration exeutes a single iteration
func (s *Context) ExecuteIteration(sc *config.Scenario, iteration int) {

	ctx, cancel := context.WithTimeout(context.Background(), sc.Sequence.Limit)
	defer cancel()

	start := time.Now()

	s.ExecuteSequence(ctx, iteration, sc)

	select {
	case <-ctx.Done():
		sc.Sequence.Stats.Error(start)
		logger.Error(nil, nil, "sequence %v on iteration: %d", ctx.Err(), iteration)
		return
	default:
	}

	sc.Sequence.Stats.Success(start)

}

// ExecuteSequence runs the sequence of requests
func (s *Context) ExecuteSequence(ctx context.Context, iteration int, sc *config.Scenario) {

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
		s.ExecuteRequest(ctx, iteration, request)
		logger.Info(request, nil, "execution complete")

	}

}

// ExecuteRequest ezxecutes a request
func (s *Context) ExecuteRequest(ctx context.Context, iteration int, request *config.Request) {

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
	wp := workerpool.New(workerPoolSize)

	start := time.Now()

	i := 0

Loop:
	for {
		if wp.WaitingQueueSize() > 0 {
			time.Sleep(time.Millisecond * 10)
		}

		wp.Submit(func() {
			s.rest.Execute(ctx, iteration, request)
		})

		// Inter-request delay
		time.Sleep(request.ThunderingHerd.Delay)

		i++

		if request.ThunderingHerd.TimeLimit > 0 {
			if time.Since(start) >= request.ThunderingHerd.TimeLimit {
				break
			}
		} else if i >= request.ThunderingHerd.Max {
			break
		}

		select {
		case <-ctx.Done():
			break Loop
		default:
		}
	}

	// Wait for everybody to complete.
	wp.StopWait()

}
