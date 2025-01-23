//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package service defines a sequence of RAPID operations
package service

import (
	"context"
	"fmt"
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
	Validate(*http.Client, *http.Response, *config.Request) error
	VerifyResponse(*http.Response, *config.Response) error
	FindResponse(*http.Response, *config.Request) *config.Response
}

// Context defines a scenario context.
type Context struct {
	datum data.Data
}

func (s *Context) addCookies(request *http.Request, r *config.Request) error {

	for i := range r.Cookies {

		// Perform any substitutions on cookies.
		ck := s.datum.Replace(r.Cookies[i].Value)

		cookies, err := http.ParseCookie(ck)
		if err != nil {
			return err
		}

		for n := range cookies {
			request.AddCookie(cookies[n])
		}
	}

	return nil
}

// GetContentReader returns a reader for an http request.
func (s *Context) GetContentReader(r *config.Request) *strings.Reader {

	// Perform any substitutions on cookie values.
	content := s.datum.Replace(r.Content)

	rdr := strings.NewReader(content)
	if rdr.Size() != 0 {
		return rdr
	}
	return nil
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

	rdr := s.GetContentReader(r)
	request, err := http.NewRequestWithContext(ctx, r.Method, url, rdr)

	if rdr != nil {
		request.Header.Add("Content-Type", r.ContentType)
		request.Header.Add("Content-Length", strconv.FormatInt(rdr.Size(), 10))
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

	err = s.addCookies(request, r)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// CreateClient creates a new http client
func (s *Context) CreateClient(_ *config.Request) (*http.Client, error) {

	client := &http.Client{}

	return client, nil
}

// FindResponse finds responses data based on returned status code.
func (s *Context) FindResponse(httpResponse *http.Response, request *config.Request) *config.Response {

	for i := range request.Responses {
		if httpResponse.StatusCode == request.Responses[i].StatusCode {
			return &request.Responses[i]
		}
	}
	return nil
}

// VerifyResponse compres response data to expected data.
func (s *Context) VerifyResponse(httpResponse *http.Response, response *config.Response) error {

	err := s.VerifyHeaders(httpResponse, response)
	if err != nil {
		return err
	}

	err = s.VerifyCookies(httpResponse, response)
	if err != nil {
		return err
	}

	err = s.VerifyContentAndExtract(httpResponse, response)
	if err != nil {
		return err
	}

	return nil
}

// Validate checks the response of a service request.
func (s *Context) Validate(client *http.Client, httpResponse *http.Response, request *config.Request) error {

	defer client.CloseIdleConnections()
	defer httpResponse.Body.Close()

	response := s.FindResponse(httpResponse, request)
	if response == nil {
		return fmt.Errorf("response not found for status code: %d for request: %s", response.StatusCode, request.Name)
	}

	err := s.VerifyResponse(httpResponse, response)

	return err
}

// Send compares the response against the expected results.
func (s *Context) Send(client *http.Client, request *http.Request, _ *config.Request) (*http.Response, error) {

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
