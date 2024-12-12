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
	assert.NotNil(t, s)
	assert.NotEmpty(t, s.Created)
	assert.NotEmpty(t, s.ID)

	err := s.ReadInConfig("../../test/configs/scenario_name.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, s.Viper)
	assert.NotEmpty(t, s.Config)
	assert.NotEmpty(t, s.Name)

	_, b := scenario.AllScenarios.Load(s.Name)
	assert.True(t, b)
}

func TestReadInConfigBadExt(t *testing.T) {

	s := scenario.New()
	assert.NotNil(t, s)
	assert.NotEmpty(t, s.Created)
	assert.NotEmpty(t, s.ID)

	err := s.ReadInConfig("../../test/configs/scenario_name.bad_ext")
	assert.NotNil(t, err)
	assert.Nil(t, s.Viper)
	assert.Empty(t, s.Config)
	assert.Empty(t, s.Name)

}

func TestReadInConfigBad(t *testing.T) {

	s := scenario.New()
	assert.NotNil(t, s)
	assert.NotEmpty(t, s.Created)
	assert.NotEmpty(t, s.ID)

	err := s.ReadInConfig("../../test/configs/scenario_no_name.yaml")
	assert.NotNil(t, err)
	assert.Nil(t, s.Viper)
	assert.Empty(t, s.Config)
}

func TestReadInConfigTee(t *testing.T) {

	s := scenario.New()
	assert.NotNil(t, s)

	err := s.ReadInConfig("../../test/configs/scenario_name.yaml")
	assert.Nil(t, err)
	assert.Contains(t, s.Config, "Lurch")

	b, err := os.ReadFile("../../test/configs/scenario_name.yaml")
	assert.Nil(t, err)
	assert.Equal(t, s.Config, string(b[:]))
}

func TestGet(t *testing.T) {

	s := scenario.New()
	err := s.ReadInConfig("../../test/configs/scenario_name.yaml")
	assert.Nil(t, err)

	ss := s.Get(s.Name)
	assert.NotNil(t, ss)
}

func TestDelete(t *testing.T) {

	s := scenario.New()
	err := s.ReadInConfig("../../test/configs/scenario_name.yaml")
	assert.Nil(t, err)

	s.Delete(s.Name)

	ss := s.Get(s.Name)
	assert.Nil(t, ss)
}
