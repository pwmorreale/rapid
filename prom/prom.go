//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package prom implements prometheus counters/etc.
package prom

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/pwmorreale/rapid/config"
)

const (
	namespace = "rapid"
)

// Context defines a context for the prom package.
type Context struct {
	reg          *prometheus.Registry
	counter      *prometheus.CounterVec
	errorCounter *prometheus.CounterVec
	histogram    *prometheus.HistogramVec
}

// New creates a new instance
func New(sc *config.Scenario) *Context {

	// Ensure initial zero'ed state.
	ctx := new(Context)

	// Abort if no prometheus config...
	if sc.Prom.PushURL == "" {
		return ctx
	}

	ctx.reg = prometheus.NewRegistry()

	ctx.counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: sc.Name,
			Name:      "responses",
			Help:      "How many HTTP Responses processed, partitioned by iteration, request name, response name, and status code",
		},
		[]string{"iteration", "request", "response", "code"},
	)

	ctx.reg.MustRegister(ctx.counter)

	ctx.errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: sc.Name,
			Name:      "errors",
			Help:      "How many HTTP client/transmission errors, partitioned by iteration and request name",
		},
		[]string{"iteration", "request"},
	)

	ctx.reg.MustRegister(ctx.errorCounter)

	minBucket := sc.Prom.Bucket.MinBucket
	if minBucket == 0 {
		minBucket = time.Duration(1)
	}

	ctx.histogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: sc.Name,
			Name:      "requests",
			Help:      "Time durations (in milliseconds) for HTTP Requests, partitioned by request name, and method",
			Buckets:   prometheus.ExponentialBucketsRange(float64(minBucket), float64(sc.Prom.Bucket.MaxBucket), sc.Prom.Bucket.Count),
		},
		[]string{"iteration", "request", "method"},
	)
	ctx.reg.MustRegister(ctx.histogram)

	return ctx
}
