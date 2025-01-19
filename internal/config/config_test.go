// Copyright © 2024 Peter W. Morreale. All Rights Reserved.

package config_test

import (
	"testing"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestReadInConfig(t *testing.T) {

	c := config.New()

	s, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	assert.NotNil(t, s)
	assert.Nil(t, err)

	assert.Equal(t, s.Name, "test-scenario")
	assert.Equal(t, len(s.Sequence.Requests), 2)

	assert.Equal(t, 2, len(s.Sequence.Requests[0].ExtraHeaders))
	assert.Equal(t, 2, len(s.Sequence.Requests[0].Cookies))
	assert.Equal(t, s.Sequence.Requests[0].ExtraHeaders[0].Name, "X-Paintbrush-for-sky")
	assert.Equal(t, s.Sequence.Requests[0].Name, "request1")

	assert.Equal(t, 2, len(s.Sequence.Requests[0].Responses))

	assert.Equal(t, 200, s.Sequence.Requests[0].Responses[0].StatusCode)
	assert.Equal(t, 500, s.Sequence.Requests[0].Responses[1].StatusCode)

	assert.Equal(t, 2, len(s.Sequence.Requests[0].Responses[0].Headers))
	assert.Equal(t, s.Sequence.Requests[0].Responses[0].Headers[0].Name, "header1")
	assert.Equal(t, s.Sequence.Requests[0].Responses[0].Headers[0].Value, "value1")

	assert.Equal(t, 2, len(s.Sequence.Requests[0].Responses[0].Cookies))

	assert.Equal(t, 500, s.Sequence.Requests[0].Responses[1].StatusCode)
	assert.Equal(t, 3, len(s.Sequence.Requests[0].Responses[1].Headers))

	assert.Equal(t, 0, len(s.Sequence.Requests[0].Responses[1].Cookies))

	assert.Equal(t, s.Sequence.Requests[1].Name, "request2")
}
