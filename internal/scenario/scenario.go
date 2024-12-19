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

// Scenario defines interfaces for executing scenarios
//
//go:generate counterfeiter -o ../../test/mocks/fake_scenario.go . Scenario
type Scenario interface {
	Get(n string) *Context
	Delete(n string)
}

// Context defines a scenario context.
type Context struct {
	Name    string
	Version string
	ID      string
	Created time.Time
	LastRun time.Time
	Viper   *viper.Viper
	Config  string
}

// AllScenarios contains all current scenario instances.
var AllScenarios sync.Map

// NewFile returns a new instance of a scenario.
func NewFile(flnm string) (*Context, error) {

	sc := &Context{
		Created: time.Now(),
		ID:      uuid.New().String(),
	}

	err := sc.fromFile(flnm)
	if err != nil {
		return nil, err
	}

	return sc, nil
}

// readReader parses config from a reader.
func (s *Context) fromReader(in io.Reader, contentType string) error {

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

// readFile creates a scenario from a config file.
func (s *Context) fromFile(flnm string) error {

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

	// Use the reader to complete.
	err = s.fromReader(in, contentType)
	if err != nil {
		return errors.Wrapf(err, "creating from file: %s", flnm)
	}

	log.Info().Str("Using config file:", flnm)

	return nil
}

// Get returns a scenario from the table.
func (s *Context) Get(n string) *Context {
	ss, b := AllScenarios.Load(n)
	if !b {
		return nil
	}
	return ss.(*Context)
}

// Delete deletes a scenario from the table.
func (s *Context) Delete(n string) {
	AllScenarios.Delete(n)
}
