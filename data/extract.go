//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package data contains the data substitution module
package data

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/antchfx/xmlquery"
	"github.com/tidwall/gjson"
)

// ExtractXML saves a datum from the specified XPath.
func (d *Context) ExtractXML(path string, r io.Reader) (string, error) {

	doc, err := xmlquery.Parse(r)
	if err != nil {
		return "", err
	}

	n, err := xmlquery.Query(doc, path)
	if err != nil {
		return "", err
	}

	if n == nil {
		return "", fmt.Errorf("XML node not found for XPath: %s", path)
	}

	return n.InnerText(), nil
}

// ExtractJSON extracts a value using a gjson path.
func (d *Context) ExtractJSON(path string, r io.Reader) (string, error) {

	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(r)
	if err != nil {
		return "", fmt.Errorf("JSON: Read error: %s", err.Error())
	}

	result := gjson.GetBytes(buf.Bytes(), path)
	if !result.Exists() {
		return "", fmt.Errorf("JSON: Not found: %s", path)
	}

	return result.String(), nil
}

// ExtractRegex extracts a value using a regular expression.
func (d *Context) ExtractRegex(rs string, r io.Reader) (string, error) {

	re, err := regexp.Compile(rs)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(r)
	if err != nil {
		return "", fmt.Errorf("REGEX: Read error: %s", err.Error())
	}

	b := re.Find(buf.Bytes())
	if len(b) == 0 {
		return "", fmt.Errorf("REGEX: value not found for expression: %s", rs)
	}

	return string(b), nil
}
