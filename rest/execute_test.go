//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package rest executes the REST calls
package rest_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/data"
	"github.com/pwmorreale/rapid/logger"
	"github.com/pwmorreale/rapid/rest"
	"github.com/stretchr/testify/assert"
)

func initTest(sc *config.Scenario) (*rest.Context, error) {

	d := data.New()
	for i := range sc.Replacements {
		err := d.AddReplacement(sc.Replacements[i].Regex, sc.Replacements[i].Value)
		if err != nil {
			return nil, err
		}

	}

	r := rest.New(sc, d)

	return r, nil
}

func initLogger(wr io.Writer) {

	opts := logger.Options{
		Handler:   "text",
		Timestamp: false,
		Level:     "Info",
		Writer:    wr,
	}

	logger.Init(&opts)
}

func TestRequestToTestServer(t *testing.T) {

	initLogger(io.Discard)

	for _, test := range []struct {
		name string

		scenario  config.Scenario
		request   config.Request
		serverTLS config.TLSConfig

		tlsError string
		count    int64
		errors   int64
	}{
		{
			name:     "http",
			tlsError: "",
			count:    1,
			errors:   0,
			request: config.Request{
				Method: "POST",
			},
		},
		{
			name:     "https (TLS)",
			tlsError: "",
			count:    1,
			errors:   0,
			scenario: config.Scenario{
				TLS: config.TLSConfig{
					CertFilePath:   "../test/certs/dev.crt",
					KeyFilePath:    "../test/certs/dev.key",
					CACertFilePath: "../test/certs/devCA.pem",
				},
			},
			request: config.Request{
				Method: "POST",
			},
			serverTLS: config.TLSConfig{
				CertFilePath:   "../test/certs/dev.crt",
				KeyFilePath:    "../test/certs/dev.key",
				CACertFilePath: "../test/certs/devCA.pem",
			},
		},
		{
			name:   "https no CA",
			count:  0,
			errors: 1,
			scenario: config.Scenario{
				TLS: config.TLSConfig{
					CertFilePath: "../test/certs/dev.crt",
					KeyFilePath:  "../test/certs/dev.key",
				},
			},
			request: config.Request{
				Method: "POST",
			},
			serverTLS: config.TLSConfig{
				CertFilePath:   "../test/certs/dev.crt",
				KeyFilePath:    "../test/certs/dev.key",
				CACertFilePath: "../test/certs/devCA.pem",
			},
		},
		{
			name:   "https no CA, w/skip verify",
			count:  1,
			errors: 0,
			scenario: config.Scenario{
				TLS: config.TLSConfig{
					CertFilePath:       "../test/certs/dev.crt",
					KeyFilePath:        "../test/certs/dev.key",
					InsecureSkipVerify: true,
				},
			},
			request: config.Request{
				Method: "POST",
			},
		},
		{
			name:     "https bad client key",
			tlsError: "",
			count:    0,
			errors:   1,
			scenario: config.Scenario{
				TLS: config.TLSConfig{
					CertFilePath:   "../test/certs/dev.crt",
					KeyFilePath:    "../test/certs/bad.key",
					CACertFilePath: "../test/certs/devCA.pem",
				},
			},
			request: config.Request{
				Method: "POST",
			},
			serverTLS: config.TLSConfig{
				CertFilePath:   "../test/certs/dev.crt",
				KeyFilePath:    "../test/certs/dev.key",
				CACertFilePath: "../test/certs/devCA.pem",
			},
		},
	} {
		// Always... to start fresh on each test.
		r, err := initTest(&test.scenario)
		assert.NotNil(t, r)
		assert.Nil(t, err)

		ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "Hello")
		}))
		assert.NotNil(t, ts, test.name)

		ts.TLS, err = r.CreateTLSConfig(test.serverTLS.CertFilePath, test.serverTLS.KeyFilePath,
			test.serverTLS.CACertFilePath, test.serverTLS.InsecureSkipVerify)
		if err != nil {
			assert.Equal(t, test.tlsError, err.Error())
			return
		}

		if ts.TLS != nil {
			ts.StartTLS()
		} else {
			ts.Start()
		}

		ctx := context.Background()

		test.request.URL = ts.URL
		r.Execute(ctx, &test.request)

		assert.Equal(t, test.count, test.request.Stats.GetCount(), test.name)
		assert.Equal(t, test.errors, test.request.Stats.GetErrors(), test.name)

		ts.Close()

	}
}
