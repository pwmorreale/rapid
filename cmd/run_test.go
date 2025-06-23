//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package cmd contains the rapid commands.
package cmd

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	URLKeyword = "TestURL"
)

func createTempScenario(t *testing.T, url string) *os.File {

	blob, err := os.ReadFile("../test/configs/run-test.yaml")
	assert.Nil(t, err)

	// Replace the URL keyword...
	re := regexp.MustCompile(URLKeyword)

	s := re.ReplaceAllLiteralString(string(blob), url)

	f, err := os.CreateTemp("", "rapid-run-*.yaml")
	assert.Nil(t, err)

	n, err := f.Write([]byte(s))
	assert.Nil(t, err)
	assert.Equal(t, n, len(s))

	f.Close()

	return f
}

func TestScenarioReWrite(t *testing.T) {

	myURL := "http://some/url"

	f := createTempScenario(t, myURL)
	assert.NotNil(t, f)

	blob, err := os.ReadFile(f.Name())
	assert.Nil(t, err)

	re, err := regexp.Compile(myURL)
	assert.Nil(t, err)

	b := re.Find(blob)
	assert.Equal(t, myURL, string(b))
}

func TestRun(t *testing.T) {

	logFormat = "text"
	logLevel = "info"
	logFilename = "/tmp/rapidRunTestFile"
	scenarioFile = ""

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Random sleep...
		d := time.Duration(rand.Int63n(int64(time.Millisecond * 100)))
		time.Sleep(d)

		// Verify sent header...
		v := r.Header.Get("MyHeader")
		assert.Equal(t, v, "some header value with fooBar")

		// Send back expected data.
		w.Header().Set("serverHeader", "something from the server")

		cookie := http.Cookie{
			Name:  "z",
			Value: "b",
		}

		http.SetCookie(w, &cookie)

		w.WriteHeader(200)
	}))
	assert.NotNil(t, ts)

	f := createTempScenario(t, ts.URL)
	assert.NotNil(t, f)
	defer os.Remove(f.Name())
	defer ts.Close()

	scenarioFile = f.Name()

	err := RunScenario(nil, []string{})
	assert.Nil(t, err)

	blob, err := os.ReadFile(logFilename)
	assert.Nil(t, err)

	assert.Contains(t, string(blob), "count=10")
	assert.Contains(t, string(blob), "errors=0")

	os.Remove(logFilename)
}
