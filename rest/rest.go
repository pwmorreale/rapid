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
	"strings"

	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/data"
	"github.com/pwmorreale/rapid/logger"
)

// Rest  defines the interface for managing requests and responses
//
//go:generate go tool counterfeiter -o ../test/mocks/fake_rest.go . Rest
type Rest interface {
	Execute(context.Context, *config.Request)
}

// Context defines a scenario context.
type Context struct {
	datum      data.Data
	sc         *config.Scenario
	httpClient *http.Client

	// For unit tests to set a mock roundtripper...
	mockRoundTripper http.RoundTripper
}

// New creates a new instance.
func New(sc *config.Scenario, d data.Data) *Context {
	return &Context{
		datum: d,
		sc:    sc,
	}
}

func (r *Context) createTLSConfig() (*tls.Config, error) {

	// No TLS config...
	if r.sc.TLS.CertFilePath == "" && r.sc.TLS.KeyFilePath == "" {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(r.sc.TLS.CertFilePath, r.sc.TLS.KeyFilePath)
	if err != nil {
		return nil, err
	}

	// If we have a CA cert path, use it and create a private pool,
	// otherwise the system pool will be used.
	var caCertPool *x509.CertPool
	if r.sc.TLS.CACertFilePath != "" {

		caCert, err := os.ReadFile(r.sc.TLS.CACertFilePath)
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
		InsecureSkipVerify: r.sc.TLS.InsecureSkipVerify, // Set to true only for testing purposes
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
		req.ContentLength = rdr.Size()
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

	// If we are configured for a single client for all requests,
	// use the cached client.
	if r.sc.UseSingleClient && r.httpClient != nil {
		return r.httpClient, nil
	}

	// If we are testing, then use the mock round tripper...
	if r.mockRoundTripper != nil {
		client.Transport = r.mockRoundTripper
		r.httpClient = client
		return client, nil
	}

	// If we have a TLS config, create and use it.
	tlsConfig, err := r.createTLSConfig()
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	// If we are using a single client, cache it...
	if r.sc.UseSingleClient {
		r.httpClient = client
	}

	return client, nil
}

// Execute creates and executes the request then validates the response.
func (r *Context) Execute(ctx context.Context, request *config.Request) {

	req, err := r.createRequest(ctx, request)
	if err != nil {
		logger.Error(request, nil, "createRequest: %v", err)
		return
	}

	client, err := r.createClient()
	if err != nil {
		logger.Error(request, nil, "createClient: %v", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(request, nil, "client.Do: %v", err)
		return
	}
	defer resp.Body.Close()

	err = r.validateResponse(resp, request)
	if err != nil {
		logger.Error(request, nil, "validateResponse: %v", err)
	}
}
