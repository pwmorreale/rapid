//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package data contins data substituion routines.
package data

import "regexp"

// Data defines interfaces for executing scenarios
type Data interface {
	Add(string, string) error
	Replace(string) string
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
func (d *Context) Add(name string, value string) error {

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
