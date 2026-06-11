//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package sequence defines a sequence of RAPID operations
package sequence

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/logger"
	"github.com/pwmorreale/rapid/rest"
)

// Sequence defines interfaces for executing scenarios
//
//go:generate go tool counterfeiter -o ../testdata/mocks/fake_sequence.go . Sequence
type Sequence interface {
	Run(context.Context, *config.Scenario) error
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
func (s *Context) Run(ctx context.Context, sc *config.Scenario) error {

	if sc.Sequence.Iterations == 0 {
		logger.Warn(nil, nil, "iterations is 0, nothing to execute")
		return nil
	}

	for i := 0; i < sc.Sequence.Iterations; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		hadError := s.ExecuteIteration(ctx, sc, i)
		if hadError && sc.Sequence.AbortOnError {
			logger.Warn(nil, nil, "aborting on error at iteration %d", i)
			return nil
		}
	}

	return nil
}

// ExecuteIteration executes a single iteration. Returns true if any request had an error.
func (s *Context) ExecuteIteration(parent context.Context, sc *config.Scenario, iteration int) bool {

	var ctx context.Context
	var cancel context.CancelFunc

	if sc.Sequence.Limit > 0 {
		ctx, cancel = context.WithTimeout(parent, sc.Sequence.Limit)
	} else {
		ctx, cancel = context.WithCancel(parent)
	}
	defer cancel()

	start := time.Now()

	hadError := s.ExecuteSequence(ctx, iteration, sc)

	select {
	case <-ctx.Done():
		sc.Sequence.Stats.Error(start)
		logger.Error(nil, nil, "sequence %v on iteration: %d", ctx.Err(), iteration)
		return true
	default:
	}

	if hadError {
		sc.Sequence.Stats.Error(start)
	} else {
		sc.Sequence.Stats.Success(start)
	}

	return hadError
}

// ExecuteSequence runs the sequence of requests. Returns true if any request had an error.
func (s *Context) ExecuteSequence(ctx context.Context, iteration int, sc *config.Scenario) bool {

	hadError := false

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
		requestHadError := s.ExecuteRequest(ctx, iteration, request, sc.Sequence.IgnoreDups)
		logger.Info(request, nil, "execution complete")

		if requestHadError {
			hadError = true
			if sc.Sequence.AbortOnError {
				break Loop
			}
		}
	}

	return hadError
}

// ExecuteRequest executes a request. Returns true if any execution had an error.
func (s *Context) ExecuteRequest(ctx context.Context, iteration int, request *config.Request, ignoreDups bool) bool {

	// Is this a once only request?
	if request.OnceOnly {
		if request.Executed {
			logger.Info(request, nil, "once_only request already executed, ignoring")
			return false
		}
		request.Executed = true

	}

	// Default to one if not specified...
	workerPoolSize := request.ThunderingHerd.Size
	if workerPoolSize == 0 {
		workerPoolSize++
	}
	wp := workerpool.New(workerPoolSize)

	var hadError atomic.Bool
	var seenErrors *sync.Map
	if ignoreDups {
		seenErrors = &sync.Map{}
	}

	start := time.Now()

	i := 0

Loop:
	for {

		for wp.WaitingQueueSize() > 0 {
			time.Sleep(time.Millisecond * 10)
		}

		wp.Submit(func() {
			errored := s.rest.Execute(ctx, iteration, request, seenErrors)
			if errored {
				hadError.Store(true)
			}
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

	return hadError.Load()
}
