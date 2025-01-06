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
	Name     string            `mapstructure:"name"`
	Version  string            `mapstructure:"version"`
	Sequence Sequence          `mapstructure:"sequence"`
	Data     map[string]string `mapstructure:"data"`
}

// Sequence contains the sequence configuration.
type Sequence struct {
	Iterations   int           `mapstructure:"iterations"`
	Limit        time.Duration `mapstructure:"time_limit"`
	Delay        time.Duration `mapstructure:"delay"`
	AbortOnError bool          `mapstructure:"abort_on_error"`
	IgnoreDups   bool          `mapstructure:"ignore_duplicate_errors"`
	Requests     []Request     `mapstructure:"requests"`
}

// Extract defines response data extraction.
type Extract struct {
	Path   string `mapstructure:"path"`
	SaveAs string `mapstructure:"save_as"`
}

// Headers contains user defined headers for inclusion with the request.
type Headers struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

// ContentData cdefines expected response data.
type ContentData struct {
	Type     string    `mapstructure:"type"`
	Contains []string  `mapstructure:"contains"`
	Extract  []Extract `mapstructure:"extract"`
}

// Response defines a REST reqponse
type Response struct {
	Content ContentData `mapstructure:"content"`
	Status  []string    `mapstructure:"expected_status"`
}

// Request defines the a request/response
type Request struct {
	Name         string            `mapstructure:"name"`
	Method       string            `mapstructure:"method"`
	URL          string            `mapstructure:"url"`
	ExtraHeaders []Headers         `mapstructure:"extra_headers"`
	Cookies      map[string]string `mapstructure:"cookies"`
	Content      string            `mapstructure:"content"`
	ContentType  string            `mapstructure:"content_type"`
	TimeLimit    time.Duration     `mapstructure:"time_limit"`
	Response     Response          `mapstructure:"response"`
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
