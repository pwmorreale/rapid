//
//  Copyright © 2023 Peter W. Morreale. All RIghts Reserved.
//

// Package scenario defines a complete testing scenario.
package scenario

import (
	"bytes"
	"io"
	"os"
	"path"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pwmorreale/rapid/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Processor defines interfaces for executing scenarios
//
//go:generate counterfeiter -o ../../test/mocks/fake_scenario.go . Processor
type Processor interface {
	ReadConfig(in io.Reader, contentType string) error
	ReadInConfig(flnm string) error
	Get(n string) *Scenario
	Delete(n string)
}

// Scenario defines a scenario instance.
type Scenario struct {
	Name    string
	ID      string
	Created time.Time
	LastRun time.Time
	Viper   *viper.Viper
	Config  string
}

// AllScenarios contains all current scenario instances.
var AllScenarios sync.Map

// New returns a new instance of a scenario.
func New() *Scenario {
	return &Scenario{
		Created: time.Now(),
		ID:      uuid.New().String(),
	}
}

// ReadConfig parses config from a reader.
func (s *Scenario) ReadConfig(in io.Reader, contentType string) error {

	var b bytes.Buffer

	// Verify viper can parse...
	if !slices.Contains(viper.SupportedExts, contentType) {
		return viper.UnsupportedConfigError(contentType)
	}

	// So we can save the original config for provenance if desired.
	tee := io.TeeReader(in, &b)

	v := viper.New()
	v.SetConfigType(contentType)
	err := v.ReadConfig(tee)
	if err != nil {
		return err
	}

	// Must have the name of this scenario.
	n := v.GetString(config.ScenarioName)
	if n == "" {
		return errors.New("Missing scenario name")
	}

	s.Name = n
	s.Viper = v
	s.Config = b.String()

	AllScenarios.Store(n, s)

	return nil
}

// ReadInConfig creates a scenario from a config file.
func (s *Scenario) ReadInConfig(flnm string) error {

	var contentType string

	in, err := os.Open(flnm)
	if err != nil {
		return err
	}

	defer in.Close()

	// Get the type of content...
	c := path.Ext(flnm)
	if c != "" {
		contentType = c[1:]
	}

	err = s.ReadConfig(in, contentType)
	if err != nil {
		return errors.Wrapf(err, "creating from file: %s", flnm)
	}

	log.Info().Str("Using config file:", flnm)

	return nil
}

// Get returns a scenario from the table.
func (s *Scenario) Get(n string) *Scenario {
	ss, b := AllScenarios.Load(n)
	if !b {
		return nil
	}
	return ss.(*Scenario)
}

// Delete deletes a scenario from the table.
func (s *Scenario) Delete(n string) {
	AllScenarios.Delete(n)
}
