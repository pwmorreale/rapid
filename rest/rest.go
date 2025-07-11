//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package rest executes REST calls
package rest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/data"
	"github.com/pwmorreale/rapid/logger"
	"github.com/pwmorreale/rapid/metrics"
)

// Rest  defines the interface for managing requests and responses
//
//go:generate go tool counterfeiter -o ../test/mocks/fake_rest.go . Rest
type Rest interface {
	Execute(context.Context, int, *config.Request)
}

// Context defines a scenario context.
type Context struct {
	datum   data.Data
	sc      *config.Scenario
	metrics metrics.Metrics

	// For unit tests to set a mock roundtripper...
	mockRoundTripper http.RoundTripper
}

// New creates a new instance.
func New(sc *config.Scenario, d data.Data) *Context {
	return &Context{
		datum:   d,
		sc:      sc,
		metrics: metrics.New(sc),
	}
}

// CreateTLSConfig creates a TLS configuration
func (r *Context) CreateTLSConfig(certPath, keyPath, caPath string, enableInsecure bool) (*tls.Config, error) {

	// No TLS config...
	if certPath == "" && keyPath == "" {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	// If we have a CA cert path, use it and create a private pool,
	// otherwise the system pool will be used.
	var caCertPool *x509.CertPool
	if caPath != "" {

		caCert, err := os.ReadFile(caPath)
		if err != nil {
			return nil, err
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	// Configure TLS
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: enableInsecure, // Set to true only for testing purposes
	}, nil
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
	return strings.NewReader(content)
}

func (r *Context) createRequest(ctx context.Context, request *config.Request) (*http.Request, error) {

	// Perform any substitutions on the url.
	url := r.datum.Replace(request.URL)

	rdr := r.getContentReader(request)
	req, err := http.NewRequestWithContext(ctx, request.Method, url, rdr)

	if request.ContentType != "" {
		req.Header.Add("Content-Type", request.ContentType)
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

func (r *Context) createClient() (*http.Client, error) {

	client := &http.Client{}

	// If we are testing, then use the mock round tripper...
	if r.mockRoundTripper != nil {
		client.Transport = r.mockRoundTripper
		return client, nil
	}

	// Get the TLC config if present...
	tls, err := r.CreateTLSConfig(r.sc.TLS.CertFilePath, r.sc.TLS.KeyFilePath,
		r.sc.TLS.CACertFilePath, r.sc.TLS.InsecureSkipVerify)

	// Probably should expose these in config...
	client.Transport = &http.Transport{
		DisableKeepAlives:   true, // Always, one request per connection.
		TLSClientConfig:     tls,  // May be nil...
		TLSHandshakeTimeout: 10 * time.Second,
		ForceAttemptHTTP2:   true,
	}

	return client, err
}

// Gestalt creates and executes the request then validates the response.
func (r *Context) Gestalt(ctx context.Context, request *config.Request) (*config.Response, error) {

	req, err := r.createRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	client, err := r.createClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return r.validateResponse(resp, request)
}

// Execute creates and executes the request then validates the response.
func (r *Context) Execute(ctx context.Context, iteration int, request *config.Request) {

	start := time.Now()

	response, err := r.Gestalt(ctx, request)
	if err != nil {
		logger.Error(request, nil, "%v", err)

		if response == nil {
			r.metrics.Errors(iteration, request.Name, metrics.NoResponseName)
			request.Stats.Error(start)
		} else {
			r.metrics.Errors(iteration, request.Name, response.Name)
			response.Stats.Error(start)
		}
		return
	}

	status := strconv.Itoa(response.StatusCode)

	r.metrics.Durations(start, iteration, request.Name, request.Method, response.Name, status)
	r.metrics.Requests(iteration, request.Name, response.Name, status)
	request.Stats.Success(start)
}
