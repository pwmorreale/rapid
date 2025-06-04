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
	"time"

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

	type header struct {
		name  string
		value string
	}

	for _, test := range []struct {
		name string

		scenario  config.Scenario
		request   config.Request
		serverTLS config.TLSConfig

		serverHeaders []header

		responseStatus int

		requestError string

		tlsError string
	}{
		{
			name:     "http",
			tlsError: "",
			request: config.Request{
				Method: "POST",
			},
		},
		{
			name:     "https (TLS)",
			tlsError: "",
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
			name: "https no CA",
			scenario: config.Scenario{
				TLS: config.TLSConfig{
					CertFilePath: "../test/certs/dev.crt",
					KeyFilePath:  "../test/certs/dev.key",
				},
			},
			requestError: "tls: failed to verify certificate: x509: “RAPID Test”",
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
			name: "https no CA, w/skip verify",
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
			name:         "https bad client key",
			tlsError:     "",
			requestError: "tls: failed to parse private key",
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
		{
			name:     "server side headers",
			tlsError: "",
			serverHeaders: []header{
				{name: "MyHeader", value: "myvalue"},
			},
			responseStatus: http.StatusOK,
			request: config.Request{
				Method: "POST",
				ExtraHeaders: []config.HeaderData{
					{Name: "MyHeader", Value: "myvalue"},
				},
				Responses: []*config.Response{
					&config.Response{
						StatusCode: http.StatusOK,
						Headers: []config.HeaderData{
							{Name: "MyHeader", Value: "myvalue"},
						},
					},
				},
			},
		},
		{
			name:           "missing expected header",
			tlsError:       "",
			responseStatus: http.StatusOK,
			requestError:   "header: Myheader not found",
			request: config.Request{
				Method: "POST",
				ExtraHeaders: []config.HeaderData{
					{Name: "MyHeader", Value: "myvalue"},
				},
				Responses: []*config.Response{
					&config.Response{
						StatusCode: http.StatusOK,
						Headers: []config.HeaderData{
							{Name: "MyHeader", Value: "myvalue"},
						},
					},
				},
			},
		},
	} {
		// Always... to start fresh on each test.
		r, err := initTest(&test.scenario)
		assert.NotNil(t, r, test.name)
		assert.Nil(t, err, test.name)

		ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			for i := range test.serverHeaders {
				w.Header().Set(test.serverHeaders[i].name, test.serverHeaders[i].value)
			}

			if test.responseStatus != 0 {
				w.WriteHeader(test.responseStatus)
			}
		}))
		assert.NotNil(t, ts, test.name)

		ts.TLS, err = r.CreateTLSConfig(test.serverTLS.CertFilePath, test.serverTLS.KeyFilePath,
			test.serverTLS.CACertFilePath, test.serverTLS.InsecureSkipVerify)
		if err != nil {
			assert.Contains(t, err.Error(), test.tlsError, test.name)
			return
		}

		if ts.TLS != nil {
			ts.StartTLS()
		} else {
			ts.Start()
		}

		ctx := context.Background()

		test.request.URL = ts.URL
		err = r.Gestalt(ctx, &test.request, time.Now())
		if test.requestError != "" {
			assert.Contains(t, err.Error(), test.requestError, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}

		ts.Close()

	}
}
