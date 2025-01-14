//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package service defines a sequence of RAPID operations
package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/data"
)

// Service defines the interface for managing requests and responses
//
//go:generate counterfeiter -o ../../test/mocks/fake_service.go . Service
type Service interface {
	CreateRequest(*config.Request) (*http.Request, error)
	CreateClient(*config.Request) (*http.Client, error)
	Send(*http.Client, *http.Request, *config.Request) (*http.Response, error)
	ValidateResponse(*http.Client, *http.Response, *config.Request) error
}

// Context defines a scenario context.
type Context struct {
	datum data.Data
}

func checkStatus(resp *http.Response, r *config.Request) error {

	for i := range r.Response.Status {
		if resp.Status == r.Response.Status[i] {
			return nil
		}
	}

	return fmt.Errorf("response status: %s not in expected status: %v", resp.Status, r.Response.Status)
}

// New returns a new context.
func New(d data.Data) *Context {
	return &Context{
		datum: d,
	}
}

// CreateRequest creates a http request.
func (s *Context) CreateRequest(r *config.Request) (*http.Request, error) {

	ctx, cancel := context.WithTimeout(context.Background(), r.TimeLimit)
	defer cancel()

	// Perform any substitutions on the url.
	url := s.datum.Replace(r.URL)

	rdr := s.getContentReader(r)
	request, err := http.NewRequestWithContext(ctx, r.Method, url, rdr)

	if rdr != nil {
		request.Header.Add("Content-Type", r.ContentType)
		request.Header.Add("Content-Length", strconv.Itoa(len(r.Content)))
	}

	if err != nil {
		return nil, err
	}

	// Add extra headers...
	for i := range r.ExtraHeaders {

		// Perform any substitutions on any extra headers.
		hv := s.datum.Replace(r.ExtraHeaders[i].Value)

		request.Header.Add(r.ExtraHeaders[i].Name, hv)
	}

	// Cookies...
	cookies := s.cookieEncode(r)
	if cookies != "" {
		request.Header.Add("Cookie", cookies)
	}

	return request, nil
}

// CreateClient creates a new http client
func (s *Context) CreateClient(_ *config.Request) (*http.Client, error) {

	client := &http.Client{}

	return client, nil
}

// ValidateResponse checks the response of a service request.
func (s *Context) ValidateResponse(client *http.Client, response *http.Response, request *config.Request) error {

	err := checkStatus(response, request)
	if err != nil {
		return err
	}

	return nil
}

// Send compares the response against the expected results.
func (s *Context) Send(client *http.Client, request *http.Request, _ *config.Request) (*http.Response, error) {

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Context) getContentReader(r *config.Request) io.Reader {

	// Perform any substitutions on cookie values.
	content := s.datum.Replace(r.Content)

	rdr := strings.NewReader(content)
	if rdr.Size() != 0 {
		return rdr
	}
	return nil
}

func (s *Context) cookieEncode(r *config.Request) string {

	var buf bytes.Buffer

	for k, v := range r.Cookies {

		// Perform any substitutions on cookie values.
		vv := s.datum.Replace(v)

		buf.WriteString(fmt.Sprintf("%s=%s; ", k, vv))
	}
	return buf.String()
}
