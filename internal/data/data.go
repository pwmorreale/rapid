//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package data contins data substituion routines.
package data

import (
	"io"
	"regexp"
)

// Data defines interfaces for executing scenarios
//
//go:generate counterfeiter -o ../../test/mocks/fake_data.go . Data
type Data interface {
	AddReplacement(string, string) error
	Replace(string) string
	ExtractJSON(string, io.Reader) (string, error)
	ExtractXML(string, io.Reader) (string, error)
	ExtractRegex(string, io.Reader) (string, error)
}

// Replacement defines a compiled regex and its associated replacement string
type Replacement struct {
	Regx *regexp.Regexp
	Repl string
}

// Context defines a sequence
type Context struct {
	All []Replacement
}

// New creates a new context instance
func New() *Context {
	return &Context{}
}

// Add creates a new regex replacement.
func (d *Context) AddReplacement(name string, value string) error {

	r := Replacement{}

	re, err := regexp.Compile(`\$` + name)
	if err != nil {
		return err
	}

	r.Regx = re
	r.Repl = value

	d.All = append(d.All, r)

	return nil
}

// Replace replaces any matches and returns a new string
func (d *Context) Replace(s string) string {

	for i := range d.All {
		s = d.All[i].Regx.ReplaceAllLiteralString(s, d.All[i].Repl)
	}
	return s
}
