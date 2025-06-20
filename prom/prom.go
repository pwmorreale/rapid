//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package prom implements prometheus counters/etc.
package prom

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/pwmorreale/rapid/config"
)

const (
	namespace = "rapid"
)

// Context defines a context for the prom package.
type Context struct {
	// Expose the registry for tests...
	Reg *prometheus.Registry

	requests  *prometheus.CounterVec
	errors    *prometheus.CounterVec
	durations *prometheus.HistogramVec
	sc        *config.Scenario
}

// New creates a new instance
func New(sc *config.Scenario) *Context {

	// Ensure initial zero'ed state.
	ctx := new(Context)

	// Abort if no prometheus config...
	if sc.Prom.PushURL == "" {
		return ctx
	}

	ctx.sc = sc

	ctx.Reg = prometheus.NewRegistry()

	ctx.requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: sc.Name,
			Name:      "responses",
			Help:      "How many HTTP Responses processed, partitioned by iteration, request name, response name, and status code",
		},
		[]string{"iteration", "request", "response", "code"},
	)

	ctx.Reg.MustRegister(ctx.requests)

	ctx.errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: sc.Name,
			Name:      "errors",
			Help:      "How many HTTP client/transmission errors, partitioned by iteration and request name",
		},
		[]string{"iteration", "request"},
	)

	ctx.Reg.MustRegister(ctx.errors)

	minBucket := sc.Prom.Bucket.MinBucket
	if minBucket == 0 {
		minBucket = time.Duration(time.Nanosecond)
	}

	maxBucket := sc.Prom.Bucket.MaxBucket
	if maxBucket == 0 {
		maxBucket = time.Duration(time.Minute)
	}

	count := sc.Prom.Bucket.Count
	if count == 0 {
		count = 5
	}

	ctx.durations = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: sc.Name,
			Name:      "requests",
			Help:      "Time durations (in milliseconds) for HTTP Requests, partitioned by iteration, request name, and method",
			Buckets:   prometheus.ExponentialBucketsRange(float64(minBucket), float64(maxBucket), count),
		},
		[]string{"iteration", "request", "method", "response", "status"},
	)
	ctx.Reg.MustRegister(ctx.durations)

	return ctx
}

// Requests is the counter for requests made.
func (p *Context) Requests(iteration int, requestName, responseName, status string) {

	if p.Reg != nil {
		p.requests.WithLabelValues(strconv.Itoa(iteration), requestName, responseName, status).Add(1)
	}
}

// Errors is the counter for errors.
func (p *Context) Errors(iteration int, requestName string) {

	if p.Reg != nil {
		p.errors.WithLabelValues(strconv.Itoa(iteration), requestName).Add(1)
	}
}

// Durations records request durations in the histogram.
func (p *Context) Durations(start time.Time, iteration int, requestName, method, responseName, status string) {

	if p.Reg != nil {
		p.durations.WithLabelValues(strconv.Itoa(iteration), requestName, method, responseName, status).Observe(float64(time.Since(start).Milliseconds()))
	}
}

func (p *Context) createClient() (*http.Client, error) {

	client := &http.Client{}

	// Get the TLC config if present...
	tls, err := p.CreateTLS()

	// Probably should expose these in config...
	client.Transport = &http.Transport{
		DisableKeepAlives:   true, // Always, one request per connection.
		TLSClientConfig:     tls,  // May be nil...
		TLSHandshakeTimeout: 10 * time.Second,
		ForceAttemptHTTP2:   true,
	}

	return client, err
}

// CreateTLS creates a TLS config for an http client.
func (p *Context) CreateTLS() (*tls.Config, error) {

	// No TLS config...
	if p.sc.Prom.TLS.CertFilePath == "" && p.sc.Prom.TLS.KeyFilePath == "" {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(p.sc.Prom.TLS.CertFilePath, p.sc.Prom.TLS.KeyFilePath)
	if err != nil {
		return nil, err
	}

	// If we have a CA cert path, use it and create a private pool,
	// otherwise the system pool will be used.
	caCertPool := new(x509.CertPool)
	if p.sc.Prom.TLS.CACertFilePath != "" {

		caCert, err := os.ReadFile(p.sc.Prom.TLS.CACertFilePath)
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
		InsecureSkipVerify: p.sc.Prom.TLS.InsecureSkipVerify, // Set to true only for testing purposes
	}, nil
}

// Push pushes the metrics to the push gateway.
func (p *Context) Push() error {

	if p.Reg == nil {
		return nil // Nothing to do...
	}

	pusher := push.New(p.sc.Prom.PushURL, p.sc.Prom.JobName)

	if len(p.sc.Prom.Headers) > 0 {
		h := make(http.Header)
		for i := range p.sc.Prom.Headers {
			h.Add(p.sc.Prom.Headers[i].Name, p.sc.Prom.Headers[i].Value)
		}
		pusher.Header(h)
	}

	client, err := p.createClient()
	if err != nil {
		return err
	}

	pusher.Client(client)

	return pusher.Gatherer(p.Reg).Push()
}
