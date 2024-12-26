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
	ParseFile(f string) error
	Viper() *viper.Viper
	Config() string
}

// Context defines a scenario context.
type Context struct {
	name    string
	version string
	id      string
	v       *viper.Viper
	conf    string
}

// New returns a context.
func New() *Context {
	return &Context{}
}

// ParseFile parse a scenario configuration
func (c *Context) ParseFile(flnm string) error {

	err := c.fromFile(flnm)
	if err != nil {
		return err
	}

	c.id = uuid.New().String()

	return nil
}

// readReader parses config from a reader.
func (c *Context) fromReader(in io.Reader, contentType string) error {

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
	// Must have the name of this scenario.
	if n == "" {
		return errors.New("Missing scenario name")
	}

	c.name = n
	c.version = v.GetString(config.ScenarioVersion)
	c.v = v
	c.conf = b.String()

	return nil
}

// readFile creates a scenario from a config file.
func (c *Context) fromFile(flnm string) error {

	var contentType string

	in, err := os.Open(flnm)
	if err != nil {
		return err
	}

	defer in.Close()

	// Get the type of content from the ext...
	t := path.Ext(flnm)
	if t != "" {
		contentType = t[1:]
	}

	// Use the reader to complete.
	err = c.fromReader(in, contentType)
	if err != nil {
		return err
	}

	log.Info().Str("Using config file:", flnm)

	return nil
}

// Viper returns the viper handle
func (c *Context) Viper() *viper.Viper {
	return c.v
}

// Config returns the original configuration as a string
func (c *Context) Config() string {
	return c.conf
}
