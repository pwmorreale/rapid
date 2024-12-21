// Copyright © 2024 Peter W. Morreale. All Rights Reserved.

package scenario_test

import (
	"os"
	"testing"

	"github.com/pwmorreale/rapid/internal/scenario"
	"github.com/test-go/testify/assert"
)

func TestReadInConfig(t *testing.T) {

	s := scenario.New()

	err := s.ParseFile("../../test/configs/scenario_name.yaml")
	assert.NotNil(t, s)
	assert.NotEmpty(t, s.Created)
	assert.NotEmpty(t, s.ID)
	assert.Nil(t, err)
	assert.NotNil(t, s.Viper)
	assert.NotEmpty(t, s.Config)
	assert.NotEmpty(t, s.Name)
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
	err := s.ParseFile("../../test/configs/scenario_name.yaml")
	assert.NotNil(t, s)
	assert.Nil(t, err)
	assert.Contains(t, s.Config, "Lurch")

	b, err := os.ReadFile("../../test/configs/scenario_name.yaml")
	assert.Nil(t, err)
	assert.Equal(t, s.Config, string(b[:]))
}
