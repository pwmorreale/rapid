//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package rest executes REST calls
package rest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/data"
	"github.com/pwmorreale/rapid/logger"
	"github.com/pwmorreale/rapid/metrics"
)

// Rest  defines the interface for managing requests and responses
//
//go:generate go tool counterfeiter -o ../testdata/mocks/fake_rest.go . Rest
type Rest interface {
	Execute(context.Context, int, *config.Request, *sync.Map) bool
	Push() error
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
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("no valid certificates found in CA file: %s", caPath)
		}
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
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(request.Method), url, rdr)
	if err != nil {
		return nil, err
	}

	if request.ContentType != "" {
		req.Header.Add("Content-Type", request.ContentType)
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

	client := &http.Client{
		Timeout: r.sc.RequestTimeout,
	}

	// If we are testing, then use the mock round tripper...
	if r.mockRoundTripper != nil {
		client.Transport = r.mockRoundTripper
		return client, nil
	}

	tlsConfig, err := r.CreateTLSConfig(r.sc.TLS.CertFilePath, r.sc.TLS.KeyFilePath,
		r.sc.TLS.CACertFilePath, r.sc.TLS.InsecureSkipVerify)
	if err != nil {
		return nil, err
	}

	client.Transport = &http.Transport{
		DisableKeepAlives:   true, // Always, one request per connection.
		TLSClientConfig:     tlsConfig,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return client, nil
}

// Gestalt creates and executes the request then validates the response.
func (r *Context) Gestalt(ctx context.Context, request *config.Request) (*config.Response, error) {

	client, err := r.createClient()
	if err != nil {
		return nil, err
	}

	maxAttempts := request.Retry.MaxAttempts
	if maxAttempts <= 1 {
		maxAttempts = 1
	}

	var resp *http.Response
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err := r.createRequest(ctx, request)
		if err != nil {
			return nil, err
		}

		r.dumpRequest(request, req)

		resp, err = client.Do(req)
		if err != nil {
			if attempt == maxAttempts {
				return nil, err
			}
			logger.Debug(request, nil, "retry %d/%d after connection error: %v", attempt, maxAttempts, err)
		} else if len(request.Retry.StatusCodes) > 0 && shouldRetry(resp.StatusCode, request.Retry.StatusCodes) && attempt < maxAttempts {
			r.dumpResponse(request, resp)
			resp.Body.Close()
			logger.Debug(request, nil, "retry %d/%d after status %d", attempt, maxAttempts, resp.StatusCode)
		} else {
			break
		}

		delay := retryDelay(attempt, request.Retry.Delay, request.Retry.MaxDelay)
		if delay > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	r.dumpResponse(request, resp)
	defer resp.Body.Close()

	return r.validateResponse(resp, request)
}

func (r *Context) dumpRequest(request *config.Request, req *http.Request) {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		logger.Debug(request, nil, "dump request error: %v", err)
		return
	}
	logger.Debug(request, nil, ">>> REQUEST >>>\n%s", string(dump))
}

func (r *Context) dumpResponse(request *config.Request, resp *http.Response) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		logger.Debug(request, nil, "dump response error: %v", err)
		return
	}
	logger.Debug(request, nil, "<<< RESPONSE <<<\n%s", string(dump))
}

// Push sends collected metrics to the Prometheus push gateway.
func (r *Context) Push() error {
	return r.metrics.Push()
}

func shouldRetry(statusCode int, retryCodes []int) bool {
	for _, code := range retryCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}

func retryDelay(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	d := baseDelay
	for i := 1; i < attempt; i++ {
		d *= 2
	}
	if maxDelay > 0 && d > maxDelay {
		d = maxDelay
	}
	return d
}

// Execute creates and executes the request then validates the response.
// Returns true if an error occurred. When seenErrors is non-nil, duplicate
// error messages are suppressed from logging (but still counted in stats).
func (r *Context) Execute(ctx context.Context, iteration int, request *config.Request, seenErrors *sync.Map) bool {

	start := time.Now()

	response, err := r.Gestalt(ctx, request)
	if err != nil {
		errMsg := err.Error()
		duplicate := false
		if seenErrors != nil {
			_, duplicate = seenErrors.LoadOrStore(errMsg, true)
		}

		if !duplicate {
			logger.Error(request, nil, "%v", err)
		}

		if response == nil {
			r.metrics.Errors(iteration, request.Name, metrics.NoResponseName)
			request.Stats.Error(start)
		} else {
			r.metrics.Errors(iteration, request.Name, response.Name)
			response.Stats.Error(start)
		}
		return true
	}

	status := strconv.Itoa(response.StatusCode)

	r.metrics.Durations(start, iteration, request.Name, request.Method, response.Name, status)
	r.metrics.Requests(iteration, request.Name, response.Name, status)
	request.Stats.Success(start)
	response.Stats.Success(start)
	return false
}
