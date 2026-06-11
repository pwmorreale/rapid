//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package sequence_test provides blackbox tests.
package sequence_test

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/logger"
	"github.com/pwmorreale/rapid/sequence"
	"github.com/pwmorreale/rapid/testdata/mocks"
	"github.com/stretchr/testify/assert"
)

var RequestDuration = time.Millisecond * 100

func fakeExecuteStub(_ context.Context, _ int, _ *config.Request, _ *sync.Map) bool {
	time.Sleep(RequestDuration)
	return false
}

func initLogger(wr io.Writer) {

	opts := logger.Options{
		Handler:   "text",
		Timestamp: false,
		Level:     "Info",
		Writer:    wr,
	}

	logger.Init(&opts)
}

func initConfig(flnm string) (*config.Scenario, error) {
	c := config.New()
	return c.ParseFile(flnm)
}

func TestOnceOnly(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	s := sequence.New(r)
	assert.NotNil(t, s)

	request := config.Request{}

	request.OnceOnly = true
	request.Executed = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	s.ExecuteRequest(ctx, 1, &request, false)
	assert.Equal(t, 0, logger.ErrorCount())
	assert.Equal(t, 0, logger.WarnCount())
	assert.Equal(t, 0, logger.DebugCount())
	assert.Equal(t, 1, logger.InfoCount())

}

func TestExecuteRequest(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = fakeExecuteStub

	s := sequence.New(r)
	assert.NotNil(t, s)

	request := config.Request{}
	request.ThunderingHerd.Max = 1
	request.ThunderingHerd.Size = 1

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	start := time.Now()
	expected := start.Add(RequestDuration)

	s.ExecuteRequest(ctx, 1, &request, false)

	actual := time.Now()
	assert.WithinDuration(t, expected, actual, time.Millisecond*10)
}

func TestThunderingHerdTimeout(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = fakeExecuteStub

	s := sequence.New(r)
	assert.NotNil(t, s)

	TestDuration := time.Second * 2

	request := config.Request{}
	request.ThunderingHerd.Size = 2000
	request.ThunderingHerd.TimeLimit = TestDuration

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	start := time.Now()
	expected := start.Add(TestDuration)

	s.ExecuteRequest(ctx, 1, &request, false)

	actual := time.Now()

	// Require finish within 10% of the duration
	assert.WithinDuration(t, expected, actual, TestDuration/10)

	// Verify work was actually performed
	assert.Greater(t, r.ExecuteCallCount(), 0)
}

func TestThunderingHerdCount(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = fakeExecuteStub

	s := sequence.New(r)
	assert.NotNil(t, s)

	request := config.Request{}
	request.ThunderingHerd.Size = 2000
	request.ThunderingHerd.Max = 1000

	TestDuration := RequestDuration

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	start := time.Now()
	expected := start.Add(TestDuration)

	s.ExecuteRequest(ctx, 1, &request, false)
	actual := time.Now()

	// Require finish within 50ms of the expected duration (race detector adds overhead)
	assert.WithinDuration(t, expected, actual, 50*time.Millisecond)

	count := request.ThunderingHerd.Max

	// Require count within 10% too...
	assert.GreaterOrEqual(t, r.ExecuteCallCount(), count)
}

func TestExecuteRequestSequenceTimeout(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = fakeExecuteStub

	s := sequence.New(r)
	assert.NotNil(t, s)

	request := config.Request{}
	request.ThunderingHerd.Max = 100
	request.ThunderingHerd.Size = 1

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	s.ExecuteRequest(ctx, 1, &request, false)

	assert.NotNil(t, ctx.Err())

}

func TestExecuteRequestDefaults(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = fakeExecuteStub

	s := sequence.New(r)
	assert.NotNil(t, s)

	request := config.Request{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*500)
	defer cancel()

	s.ExecuteRequest(ctx, 1, &request, false)

	assert.Nil(t, ctx.Err())

	assert.Equal(t, 1, r.ExecuteCallCount())

}

func TestRun(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = fakeExecuteStub

	s := sequence.New(r)
	assert.NotNil(t, s)

	sc, err := initConfig("../testdata/configs/test_scenario.yaml")
	assert.Nil(t, err)

	err = s.Run(context.Background(), sc)
	assert.Nil(t, err)

	// N.B. OnceOnly is set for the first request,
	// and the max calls per request is 1.
	assert.Equal(t, 11, r.ExecuteCallCount())

}

func TestAbortOnError(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = func(_ context.Context, _ int, _ *config.Request, _ *sync.Map) bool {
		return true
	}

	s := sequence.New(r)
	assert.NotNil(t, s)

	sc := &config.Scenario{
		Sequence: config.Sequence{
			Iterations:   10,
			AbortOnError: true,
			Requests: []config.Request{
				{
					Name:   "fail-request",
					Method: "get",
					ThunderingHerd: config.Stampede{
						Max:  1,
						Size: 1,
					},
				},
			},
		},
	}

	err := s.Run(context.Background(), sc)
	assert.Nil(t, err)

	// With abort_on_error, should stop after the first iteration.
	assert.Equal(t, 1, r.ExecuteCallCount())
}

func TestAbortOnErrorDisabled(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = func(_ context.Context, _ int, _ *config.Request, _ *sync.Map) bool {
		return true
	}

	s := sequence.New(r)
	assert.NotNil(t, s)

	sc := &config.Scenario{
		Sequence: config.Sequence{
			Iterations:   5,
			AbortOnError: false,
			Requests: []config.Request{
				{
					Name:   "fail-request",
					Method: "get",
					ThunderingHerd: config.Stampede{
						Max:  1,
						Size: 1,
					},
				},
			},
		},
	}

	err := s.Run(context.Background(), sc)
	assert.Nil(t, err)

	// Without abort_on_error, all iterations should execute.
	assert.Equal(t, 5, r.ExecuteCallCount())
}

func TestAbortOnErrorStopsSequence(t *testing.T) {

	initLogger(io.Discard)

	callCount := 0
	r := &mocks.FakeRest{}
	r.ExecuteStub = func(_ context.Context, _ int, req *config.Request, _ *sync.Map) bool {
		callCount++
		// Only the second request errors.
		return req.Name == "second"
	}

	s := sequence.New(r)
	assert.NotNil(t, s)

	sc := &config.Scenario{
		Sequence: config.Sequence{
			Iterations:   3,
			AbortOnError: true,
			Requests: []config.Request{
				{
					Name:   "first",
					Method: "get",
					ThunderingHerd: config.Stampede{
						Max:  1,
						Size: 1,
					},
				},
				{
					Name:   "second",
					Method: "get",
					ThunderingHerd: config.Stampede{
						Max:  1,
						Size: 1,
					},
				},
				{
					Name:   "third",
					Method: "get",
					ThunderingHerd: config.Stampede{
						Max:  1,
						Size: 1,
					},
				},
			},
		},
	}

	err := s.Run(context.Background(), sc)
	assert.Nil(t, err)

	// First iteration: "first" succeeds, "second" errors, "third" is skipped.
	// abort_on_error stops further iterations.
	assert.Equal(t, 2, callCount)
}

func TestIgnoreDuplicateErrors(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = func(_ context.Context, _ int, _ *config.Request, seen *sync.Map) bool {
		// Verify that a seen-errors map was provided.
		assert.NotNil(t, seen)
		return true
	}

	s := sequence.New(r)
	assert.NotNil(t, s)

	sc := &config.Scenario{
		Sequence: config.Sequence{
			Iterations: 1,
			IgnoreDups: true,
			Requests: []config.Request{
				{
					Name:   "herd-request",
					Method: "get",
					ThunderingHerd: config.Stampede{
						Max:  5,
						Size: 5,
					},
				},
			},
		},
	}

	err := s.Run(context.Background(), sc)
	assert.Nil(t, err)

	assert.Equal(t, 5, r.ExecuteCallCount())
}

func TestIgnoreDuplicateErrorsDisabled(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = func(_ context.Context, _ int, _ *config.Request, seen *sync.Map) bool {
		// When ignore_duplicate_errors is false, no seen-errors map should be passed.
		assert.Nil(t, seen)
		return false
	}

	s := sequence.New(r)
	assert.NotNil(t, s)

	sc := &config.Scenario{
		Sequence: config.Sequence{
			Iterations: 1,
			IgnoreDups: false,
			Requests: []config.Request{
				{
					Name:   "normal-request",
					Method: "get",
					ThunderingHerd: config.Stampede{
						Max:  3,
						Size: 3,
					},
				},
			},
		},
	}

	err := s.Run(context.Background(), sc)
	assert.Nil(t, err)

	assert.Equal(t, 3, r.ExecuteCallCount())
}
