//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package config contains config variables.and utilities
package config

import (
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// Configuration defines the interface for managing the scenario
//
//go:generate counterfeiter -o ../../test/mocks/fake_config.go . Configuration
type Configuration interface {
	ParseFile(f string) error
}

// Context defines a scenario context.
type Context struct {
	v  *viper.Viper
	id string
}

// Scenario defines the entire configuration.
type Scenario struct {
	Name    string   `mapstructure:"name"`
	Version string   `mapstructure:"version"`
	Seq     Sequence `mapstructure:"sequence"`
}

// Sequences contains the sequence configuration.
type Sequence struct {
	Iterations int           `mapstructure:"iterations"`
	Duration   time.Duration `mapstructure:"time_limit"`
	ErrorAbort bool          `mapstructure:"abort_on_error"`
	IgnoreDups bool          `mapstructure:"ignore_duplicate_errors"`
	Reqs       []Request     `mapstructure:"requests"`
}

// ContentData cdefines expected response data.
type ContentData struct {
	Type     string   `mapstructure:"type"`
	Contains []string `mapstructure:"contains"`
	Extract  string   `mapstructure:"extract"`
	SaveAs   string   `mapstructure:"save_as"`
}

// Response defines a REST reqponse
type Response struct {
	Content ContentData `mapstructure:"content"`
	Status  []int       `mapstructure:"expected_status"`
}

// Request defines the a request/response
type Request struct {
	Name         string            `mapstructure:"name"`
	Scheme       string            `mapstructure:"scheme"`
	Path         string            `mapstructure:"path"`
	Host         string            `mapstructure:"host"`
	Query        map[string]string `mapstructure:"query"`
	ExtraHeaders map[string]string `mapstructure:"extra_headers"`
	Content      string            `mapstructure:"content"`
	ContentType  string            `mapstructure:"content_type"`
	TimeLimit    time.Duration     `mapstructure:"time_limit"`
	Rsp          Response          `mapstructure:"response"`
}

// New creates a new context instance
func New() *Context {
	return &Context{}
}

// ParseFile parse a scenario configuration
func (c *Context) ParseFile(flnm string) (*Scenario, error) {

	var s Scenario

	v := viper.New()
	v.SetConfigFile(flnm)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err

	}

	// Ensure text config matches what we expect.
	err = v.UnmarshalExact(&s)
	if err != nil {
		return nil, err
	}

	c.v = v
	c.id = uuid.New().String()

	return &s, nil
}
