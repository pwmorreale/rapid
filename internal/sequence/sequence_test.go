//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package sequence_test provides blackbox tests.
package sequence_test

import (
	"testing"
	"time"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/sequence"
	"github.com/pwmorreale/rapid/test/mocks"
	"github.com/stretchr/testify/assert"
)

var r = &mocks.FakeRest{}

func fakeExecuteStub(request *config.Request) {
	time.Sleep(1)
}

func TestBadConfig(t *testing.T) {

	s := sequence.New(r)
	assert.NotNil(t, s)

	err := s.Run("../../test/configs/bad.yaml")
	assert.Equal(t, "While parsing config: yaml: line 19: could not find expected ':'", err.Error())

}

func TestExecuteRequest(t *testing.T) {

	r.ExecuteStub = fakeExecuteStub

	s := sequence.New(r)
	assert.NotNil(t, s)

}
