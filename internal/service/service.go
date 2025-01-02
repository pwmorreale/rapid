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
	Create(*config.Request) (*http.Request, error)
	Send(*http.Request, *config.Request) (*http.Response, error)
	Validate(*http.Response, *config.Request) error
}

// Context defines a scenario context.
type Context struct {
	savedContent map[string]string
}

// New returns a new context.
func New() *Context {
	return &Context{}
}

// Create creates a http request.
func (s *Context) Create(r *config.Request) (*http.Request, error) {

	ctx, cancel := context.WithTimeout(context.Background(), r.TimeLimit)
	defer cancel()

	u := buildURL(r)

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
	for h, v := range r.ExtraHeaders {
		request.Header.Add(h, v)
	}

	// Cookies...
	cookies := cookieEncode(r)
	if cookies != "" {
		request.Header.Add("Cookie", cookies)
	}

	return request, nil
}

// Send executes the request
func (s *Context) Send(hr *http.Request, _ *config.Request) (*http.Response, error) {

	httpClient := &http.Client{}

	return httpClient.Do(hr)
}

// Validate compares the response against the expected results.
func (s *Context) Validate(_ *http.Response, _ *config.Request) error {
	return nil
}

// buildURL creates the URL from the configuration
func buildURL(r *config.Request) string {

	var u url.URL

	u.Scheme = r.Scheme
	u.Path = r.Path
	u.Host = r.Host
	u.User = getUserinfo(r)

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
