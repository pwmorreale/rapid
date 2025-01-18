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

// Various constants...
const (
	DefaultContentLimit = 4096

	TypeRegex = "regex"
)

// Configuration defines the interface for managing the scenario
//
//go:generate counterfeiter -o ../../test/mocks/fake_config.go . Configuration
type Configuration interface {
	ParseFile(f string) (*Scenario, error)
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
type ExtractData struct {
	Type      string `mapstructure:"type"`
	Path      string `mapstructure:"path"`
	DastaName string `mapstructure:"data_name"`
}

// HeaderData contains user defined headers for inclusion with the request.
type HeaderData struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

// ContentData cdefines expected response data.
type ContentData struct {
	MediaType string        `mapstructure:"content_type"`
	MaxSize   int           `mapstructure:"max_content"`
	Contains  []string      `mapstructure:"contains"`
	Extract   []ExtractData `mapstructure:"extract"`
}

// CookieData defines a cookie string
type CookieData struct {
	Value string `mapstructure:"value"`
}

// Response defines a REST response
type Response struct {
	StatusCode int          `mapstructure:"status_code"`
	Headers    []HeaderData `mapstructure:"headers"`
	Cookies    []CookieData `mapstructure:"cookies"`
	Content    ContentData  `mapstructure:"content"`
}

// Request defines the a request/response
type Request struct {
	Name         string        `mapstructure:"name"`
	Method       string        `mapstructure:"method"`
	URL          string        `mapstructure:"url"`
	ExtraHeaders []HeaderData  `mapstructure:"extra_headers"`
	Cookies      []CookieData  `mapstructure:"cookies"`
	Content      string        `mapstructure:"content"`
	ContentType  string        `mapstructure:"content_type"`
	TimeLimit    time.Duration `mapstructure:"time_limit"`
	Responses    []Response    `mapstructure:"responses"`
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
