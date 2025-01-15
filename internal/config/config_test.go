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

	assert.Equal(t, len(s.Sequence.Requests[0].ExtraHeaders), 2)
	assert.Equal(t, s.Sequence.Requests[0].ExtraHeaders[0].Name, "X-Paintbrush-for-sky")
	assert.Equal(t, s.Sequence.Requests[0].Name, "request1")
	assert.Equal(t, len(s.Sequence.Requests[0].Response.Status), 1)

	assert.Equal(t, s.Sequence.Requests[1].Name, "request2")
	assert.Equal(t, s.Sequence.Requests[1].Response.Status[0], "200")
	assert.Equal(t, len(s.Sequence.Requests[0].Response.Status), 1)
}

func TestCompileExpressions(t *testing.T) {
	c := config.New()

	s, err := c.ParseFile("../../test/configs/single_request.yaml")
	assert.NotNil(t, s)
	assert.Nil(t, err)

	err = c.CompileExpressions(s)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(s.Sequence.Requests[0].Response.Content.Contains))
	assert.Equal(t, 2, len(s.Sequence.Requests[0].Response.Content.Regex))

	// Content.Type == "json"
	assert.Nil(t, s.Sequence.Requests[0].Response.Content.Extract[0].RegExp)
	assert.Nil(t, s.Sequence.Requests[0].Response.Content.Extract[1].RegExp)

	// Check Content.Type of Regex...
	s.Sequence.Requests[0].Response.Content.Type = config.TypeRegex

	err = c.CompileExpressions(s)
	assert.Nil(t, err)

	assert.NotNil(t, s.Sequence.Requests[0].Response.Content.Extract[0].RegExp)
	assert.NotNil(t, s.Sequence.Requests[0].Response.Content.Extract[1].RegExp)

}
