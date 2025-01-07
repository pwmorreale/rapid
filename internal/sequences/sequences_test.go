//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

package sequences_test

import (
	"errors"
	"testing"
	"time"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/sequences"
	"github.com/pwmorreale/rapid/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	rpt := &mocks.FakeReport{}
	srv := &mocks.FakeService{}

	seq := sequences.New(srv, rpt)
	assert.NotNil(t, seq)
}

func TestRun(t *testing.T) {
	rpt := &mocks.FakeReport{}
	srv := &mocks.FakeService{}

	seq := sequences.New(srv, rpt)
	assert.NotNil(t, seq)

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	err = seq.Run(sc)
	assert.Nil(t, err)

	// Two requests in the test sequence...
	assert.Equal(t, sc.Sequence.Iterations*2, srv.CreateRequestCallCount())
	assert.Equal(t, sc.Sequence.Iterations*2, srv.SendCallCount())
	assert.Equal(t, sc.Sequence.Iterations*2, srv.ValidateResponseCallCount())
}

func TestRunFailSecondRequest(t *testing.T) {
	rpt := &mocks.FakeReport{}
	srv := &mocks.FakeService{}

	seq := sequences.New(srv, rpt)
	assert.NotNil(t, seq)

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	srv.CreateRequestReturnsOnCall(1, nil, errors.New("blowing chunks"))

	err = seq.Run(sc)

	assert.Equal(t, err, errors.New("blowing chunks"))
	assert.Equal(t, 2, srv.CreateRequestCallCount())
	assert.Equal(t, 1, srv.CreateClientCallCount())
	assert.Equal(t, 1, srv.SendCallCount())
	assert.Equal(t, 1, srv.ValidateResponseCallCount())
}

func TestRunFailClient(t *testing.T) {
	rpt := &mocks.FakeReport{}
	srv := &mocks.FakeService{}

	seq := sequences.New(srv, rpt)
	assert.NotNil(t, seq)

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	srv.CreateClientReturnsOnCall(0, nil, errors.New("blowing chunks"))

	err = seq.Run(sc)

	assert.Equal(t, err, errors.New("blowing chunks"))
	assert.Equal(t, 1, srv.CreateRequestCallCount())
	assert.Equal(t, 1, srv.CreateClientCallCount())
	assert.Equal(t, 0, srv.SendCallCount())
	assert.Equal(t, 0, srv.ValidateResponseCallCount())
}

func TestRunFailSend(t *testing.T) {
	rpt := &mocks.FakeReport{}
	srv := &mocks.FakeService{}

	seq := sequences.New(srv, rpt)
	assert.NotNil(t, seq)

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	srv.SendReturnsOnCall(0, nil, errors.New("blowing chunks"))

	err = seq.Run(sc)

	assert.Equal(t, err, errors.New("blowing chunks"))
	assert.Equal(t, 1, srv.CreateRequestCallCount())
	assert.Equal(t, 1, srv.CreateClientCallCount())
	assert.Equal(t, 1, srv.SendCallCount())
	assert.Equal(t, 0, srv.ValidateResponseCallCount())
}

func TestRunFailValidate(t *testing.T) {
	rpt := &mocks.FakeReport{}
	srv := &mocks.FakeService{}

	seq := sequences.New(srv, rpt)
	assert.NotNil(t, seq)

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	srv.ValidateResponseReturnsOnCall(0, errors.New("blowing chunks"))

	err = seq.Run(sc)

	assert.Equal(t, err, errors.New("blowing chunks"))
	assert.Equal(t, 1, srv.CreateRequestCallCount())
	assert.Equal(t, 1, srv.SendCallCount())
	assert.Equal(t, 1, srv.ValidateResponseCallCount())
}

func TestRunAbortOnError(t *testing.T) {
	rpt := &mocks.FakeReport{}
	srv := &mocks.FakeService{}

	seq := sequences.New(srv, rpt)
	assert.NotNil(t, seq)

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	sc.Sequence.AbortOnError = false

	srv.ValidateResponseReturnsOnCall(0, errors.New("blowing chunks"))

	err = seq.Run(sc)
	assert.Nil(t, err)

	assert.Equal(t, 20, srv.CreateRequestCallCount())
	assert.Equal(t, 20, srv.CreateClientCallCount())
	assert.Equal(t, 20, srv.SendCallCount())
	assert.Equal(t, 20, srv.ValidateResponseCallCount())
}
func TestRunExceedTimeLimit(t *testing.T) {
	rpt := &mocks.FakeReport{}
	srv := &mocks.FakeService{}

	seq := sequences.New(srv, rpt)
	assert.NotNil(t, seq)

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	// Ensure we exceed the time limit
	sc.Sequence.Delay, _ = time.ParseDuration("100ms")
	sc.Sequence.Limit, _ = time.ParseDuration("10ms")

	err = seq.Run(sc)
	assert.Nil(t, err)

	assert.Equal(t, 1, srv.CreateRequestCallCount())
	assert.Equal(t, 1, srv.SendCallCount())
	assert.Equal(t, 1, srv.ValidateResponseCallCount())
}
