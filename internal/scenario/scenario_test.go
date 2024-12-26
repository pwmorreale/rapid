// Copyright © 2024 Peter W. Morreale. All Rights Reserved.

package scenario_test

import (
	"os"
	"testing"

	"github.com/pwmorreale/rapid/internal/scenario"
	"github.com/stretchr/testify/assert"
)

func TestReadInConfig(t *testing.T) {

	s := scenario.New()

	err := s.ParseFile("../../test/configs/scenario_name.yaml")
	assert.NotNil(t, s)
	assert.Nil(t, err)
}

func TestReadInConfigBadExt(t *testing.T) {

	s := scenario.New()
	err := s.ParseFile("../../test/configs/scenario_name.bad_ext")
	assert.NotNil(t, err)

}

func TestReadInConfigBad(t *testing.T) {

	s := scenario.New()
	err := s.ParseFile("../../test/configs/scenario_no_name.yaml")
	assert.NotNil(t, err)
}

func TestReadInConfigTee(t *testing.T) {

	s := scenario.New()
	err := s.ParseFile("../../test/configs/test_scenario.yaml")
	assert.NotNil(t, s)
	assert.Nil(t, err)

	b, err := os.ReadFile("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)
	assert.Equal(t, s.Config(), string(b[:]))
}
