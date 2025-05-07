//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package rest executes REST calls
package rest

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/data"
	"github.com/pwmorreale/rapid/internal/logger"
)

// Rest  defines the interface for managing requests and responses
//
//go:generate go tool counterfeiter -o ../../test/mocks/fake_rest.go . Rest
type Rest interface {
	Execute(context.Context, *config.Request)
}

// Context defines a scenario context.
type Context struct {
	datum        data.Data
	roundTripper http.RoundTripper
}

// New creates a new instance.
func New(d data.Data) *Context {
	return &Context{
		datum:        d,
		roundTripper: nil, // N.B. Used by the test package to mock out the default
	}
}

func (r *Context) addCookies(req *http.Request, request *config.Request) error {

	for i := range request.Cookies {

		// Perform any substitutions on cookies.
		ck := r.datum.Replace(request.Cookies[i].Value)

		cookies, err := http.ParseCookie(ck)
		if err != nil {
			return err
		}

		for n := range cookies {
			req.AddCookie(cookies[n])
		}
	}

	return nil
}

func (r *Context) getContentReader(request *config.Request) *strings.Reader {

	// Perform any substitutions on cookie values.
	content := r.datum.Replace(request.Content)

	rdr := strings.NewReader(content)
	if rdr.Size() != 0 {
		return rdr
	}
	return nil
}

func (r *Context) createRequest(ctx context.Context, request *config.Request) (*http.Request, error) {

	// Perform any substitutions on the url.
	url := r.datum.Replace(request.URL)

	rdr := r.getContentReader(request)
	req, err := http.NewRequestWithContext(ctx, request.Method, url, rdr)

	if rdr != nil {
		req.Header.Add("Content-Type", request.ContentType)
		req.Header.Add("Content-Length", strconv.FormatInt(rdr.Size(), 10))
	}

	if err != nil {
		return nil, err
	}

	// Add extra headers...
	for i := range request.ExtraHeaders {

		// Perform any substitutions on any extra headers.
		hv := r.datum.Replace(request.ExtraHeaders[i].Value)

		req.Header.Add(request.ExtraHeaders[i].Name, hv)
	}

	err = r.addCookies(req, request)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Execute creates and executes the request then validates the response.
func (r *Context) Execute(ctx context.Context, request *config.Request) {

	req, err := r.createRequest(ctx, request)
	if err != nil {
		logger.Error(request, nil, "createRequest: %v", err)
		return
	}

	// N.B.  We specify the transport solely to enble testing for this
	// routine.  In the normal path r.roundTrippe will be nil, which implies the
	// default transport.  Tests will specify a mock transport.
	client := &http.Client{
		Transport: r.roundTripper,
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(request, nil, "client.Do: %v", err)
		return
	}
	defer client.CloseIdleConnections()
	defer resp.Body.Close()

	err = r.validateResponse(resp, request)
	if err != nil {
		logger.Error(request, nil, "validateResponse: %v", err)
	}
}
