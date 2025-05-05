//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package sequence_test provides blackbox tests.
package sequence_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/logger"
	"github.com/pwmorreale/rapid/internal/sequence"
	"github.com/pwmorreale/rapid/test/mocks"
	"github.com/stretchr/testify/assert"
)

func fakeExecuteStub(ctx context.Context, request *config.Request) {
	time.Sleep(time.Second)
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

	s.ExecuteRequest(ctx, &request)
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

	s.ExecuteRequest(ctx, &request)

	stats := s.GetStats()
	assert.Equal(t, int64(1), stats.Calls)

	expected := time.Now().Add(time.Second * 1)
	actual := time.Now().Add(stats.RestCallTime)
	assert.WithinDuration(t, expected, actual, time.Millisecond*100)
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

	s.ExecuteRequest(ctx, &request)

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

	s.ExecuteRequest(ctx, &request)

	assert.Nil(t, ctx.Err())

	stats := s.GetStats()
	assert.Equal(t, int64(1), stats.Calls)

}

func TestRun(t *testing.T) {

	initLogger(io.Discard)

	r := &mocks.FakeRest{}
	r.ExecuteStub = fakeExecuteStub

	s := sequence.New(r)
	assert.NotNil(t, s)

	sc, err := initConfig("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	err = s.Run(sc)
	assert.Nil(t, err)

	// N.B. OnceOnly is set for the first request,
	// and the max calls per request is 1.

	stats := s.GetStats()
	assert.Equal(t, int64(11), stats.Calls)
	assert.Equal(t, 11, r.ExecuteCallCount())

}
