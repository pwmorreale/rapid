//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package data contins data substituion routines.
package data

import "regexp"

// Data defines interfaces for executing scenarios
type Data interface {
	AddReplacement(string, string) error
	AddMatcher(string) error
	Replace(string) string
	Match(string) bool
}

// Replacement defines a compiled regex and its associated replacement string
type Replacement struct {
	Regx *regexp.Regexp
	Repl string
}

// Context defines a sequence
type Context struct {
	All      []Replacement
	Matchers []*regexp.Regexp
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

// Add creates a new regex replacement.
func (d *Context) AddMatcher(rs string) error {

	r, err := regexp.Compile(rs)
	if err == nil {
		d.Matchers = append(d.Matchers, r)
	}
	return err
}

// Replace replaces any matches and returns a new string
func (d *Context) Replace(s string) string {

	for i := range d.All {
		s = d.All[i].Regx.ReplaceAllLiteralString(s, d.All[i].Repl)
	}
	return s
}

// Match finds matches in the given string.
func (d *Context) Match(s string) bool {

	for i := range d.Matchers {
		if d.Matchers[i].MatchString(s) == true {
			return true
		}
	}
	return false
}
