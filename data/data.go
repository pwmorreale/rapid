//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package data contins data substitution routines.
package data

import (
	"io"
	"regexp"
	"sync"
)

// Data defines interfaces for executing scenarios
//
//go:generate go tool counterfeiter -o ../testdata/mocks/fake_data.go . Data
type Data interface {
	AddReplacement(string, string) error
	Replace(string) string
	Lookup(string) string
	Len() int
	ExtractJSON(string, io.Reader) (string, error)
	ExtractXML(string, io.Reader) (string, error)
	ExtractRegex(string, io.Reader) (string, error)
}

// Replacement defines a compiled regex and its associated replacement string
type Replacement struct {
	name  string
	regx  *regexp.Regexp
	value string
}

// Context defines a sequence
type Context struct {
	mu  sync.RWMutex
	all []Replacement
}

// New creates a new context instance
func New() *Context {
	return &Context{}
}

// AddReplacement creates a new regex replacement.
func (d *Context) AddReplacement(name string, value string) error {

	r := Replacement{}

	re, err := regexp.Compile(name)
	if err != nil {
		return err
	}

	r.name = name
	r.regx = re
	r.value = value

	d.mu.Lock()
	d.all = append(d.all, r)
	d.mu.Unlock()

	return nil
}

// Replace replaces any matches and returns a new string
func (d *Context) Replace(s string) string {

	d.mu.RLock()
	defer d.mu.RUnlock()

	for i := range d.all {
		s = d.all[i].regx.ReplaceAllLiteralString(s, d.all[i].value)
	}
	return s
}

// Lookup returns the replacement text (value) for a name
func (d *Context) Lookup(n string) string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for i := range d.all {
		if n == d.all[i].name {
			return d.all[i].value
		}
	}
	return ""
}

// Len returns the number of replacement elements.
func (d *Context) Len() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.all)
}
