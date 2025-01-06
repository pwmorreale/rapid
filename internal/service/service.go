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
	"net/url"
	"strconv"
	"strings"

	"github.com/pwmorreale/rapid/internal/config"
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
	savedContent map[string]string
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
func New() *Context {
	return &Context{}
}

// Create creates a http request.
func (s *Context) CreateRequest(r *config.Request) (*http.Request, error) {

	ctx, cancel := context.WithTimeout(context.Background(), r.TimeLimit)
	defer cancel()

	u := s.BuildURL(r)

	rdr := s.getContentReader(r)
	request, err := http.NewRequestWithContext(ctx, r.Method, u, rdr)

	if rdr != nil {
		request.Header.Add("Content-Type", r.ContentType)
		request.Header.Add("Content-Length", strconv.Itoa(len(r.Content)))
	}

	if err != nil {
		return nil, err
	}

	// Add extra headers...
	for i := range r.ExtraHeaders {
		request.Header.Add(r.ExtraHeaders[i].Name, r.ExtraHeaders[i].Value)
	}

	// Cookies...
	cookies := cookieEncode(r)
	if cookies != "" {
		request.Header.Add("Cookie", cookies)
	}

	return request, nil
}

// CreateClient creates a new http client
func (s *Context) CreateClient(req *config.Request) (*http.Client, error) {

	client := &http.Client{}

	return client, nil
}

// ValidateResponse checks the response of a service request.
func (s *Context) ValidateResponse(*http.Client, *http.Response, *config.Request) error {
	return nil
}

// Send compares the response against the expected results.
func (s *Context) Send(client *http.Client, request *http.Request, r *config.Request) (*http.Response, error) {

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// BuildURL creates the URL from the configuration
func (s *Context) BuildURL(r *config.Request) string {

	var u url.URL

	u.Scheme = r.Scheme
	u.Path = r.Path
	u.Host = r.Host
	u.User = getUserinfo(r)
	u.Fragment = r.Fragment

	q := u.Query()
	for k, v := range r.Query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	us := u.String()

	return us
}

func (s *Context) getContentReader(r *config.Request) io.Reader {

	d := r.Content

	//
	// If first byte is a '$', then assume a lookup in the
	// saved content.
	//
	if strings.IndexByte(r.Content, '$') == 0 {
		d = s.savedContent[strings.TrimLeft(r.Content, "$")]
	}

	rdr := strings.NewReader(d)
	if rdr.Size() != 0 {
		return rdr
	}
	return nil
}

func cookieEncode(r *config.Request) string {

	var buf bytes.Buffer

	for k, v := range r.Cookies {
		buf.WriteString(fmt.Sprintf("%s=%s; ", k, v))
	}
	return buf.String()
}

func getUserinfo(r *config.Request) *url.Userinfo {

	if r.Username != "" && r.Password != "" {
		return url.UserPassword(r.Username, r.Password)
	} else if r.Username != "" {
		return url.User(r.Username)
	}

	return nil
}
