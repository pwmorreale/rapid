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
	assert.Equal(t, len(s.Seq.Reqs), 2)

	assert.Equal(t, len(s.Seq.Reqs[0].ExtraHeaders), 2)
	assert.Equal(t, s.Seq.Reqs[0].ExtraHeaders[0].Name, "Foo")
	assert.Equal(t, s.Seq.Reqs[0].Name, "request1")
	assert.Equal(t, len(s.Seq.Reqs[0].Rsp.Status), 3)

	assert.Equal(t, s.Seq.Reqs[1].Name, "request2")
	assert.Equal(t, s.Seq.Reqs[1].Rsp.Status[0], 500)
	assert.Equal(t, len(s.Seq.Reqs[0].Rsp.Status), 3)
}
