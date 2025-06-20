package prom_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/common/expfmt"
	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/prom"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {

	var (
		lastMethod string
		lastBody   []byte
		lastPath   string
		lastHeader http.Header
	)

	// Fake a Pushgateway that responds with 202 to DELETE and with 200 in
	// all other cases.
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lastMethod = r.Method
			var err error
			lastBody, err = io.ReadAll(r.Body)
			lastHeader = r.Header
			if err != nil {
				t.Fatal(err)
			}
			lastPath = r.URL.EscapedPath()
			w.Header().Set("Content-Type", `text/plain; charset=utf-8`)
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusAccepted)
				return
			}

			w.WriteHeader(http.StatusOK)
		}),
	)
	defer ts.Close()

	c := config.New()
	sc, err := c.ParseFile("../test/configs/test_scenario.yaml")
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	sc.Prom.Headers = append(sc.Prom.Headers, config.HeaderData{Name: "foo", Value: "bar"})
	sc.Prom.JobName = "test"

	sc.Prom.PushURL = ts.URL

	pc := prom.New(sc)
	assert.NotNil(t, pc)

	pc.Requests(1, "req1", "resp1", "200")
	pc.Errors(0, "req1")

	start := time.Now()
	time.Sleep(time.Millisecond * 10)

	pc.Durations(start, 1, "req2", "GET", "resp3", "201")

	mfs, err := pc.Reg.Gather()
	assert.Nil(t, err)

	buf := &bytes.Buffer{}
	enc := expfmt.NewEncoder(buf, expfmt.NewFormat(expfmt.TypeProtoDelim))

	for _, mf := range mfs {
		err := enc.Encode(mf)
		assert.Nil(t, err)
	}
	wantBody := buf.Bytes()

	err = pc.Push()
	assert.Nil(t, err)

	assert.Equal(t, http.MethodPut, lastMethod)
	assert.Equal(t, "/metrics/job/test", lastPath)

	assert.True(t, bytes.Equal(wantBody, lastBody))

	assert.NotNil(t, lastHeader)
	assert.Contains(t, lastHeader, "Foo")
}
func TestTLSMetrics(t *testing.T) {

	var (
		lastMethod string
		lastBody   []byte
		lastPath   string
		lastHeader http.Header
	)

	// Fake a Pushgateway that responds with 202 to DELETE and with 200 in
	// all other cases.
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lastMethod = r.Method
			var err error
			lastBody, err = io.ReadAll(r.Body)
			lastHeader = r.Header
			if err != nil {
				t.Fatal(err)
			}
			lastPath = r.URL.EscapedPath()
			w.Header().Set("Content-Type", `text/plain; charset=utf-8`)
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusAccepted)
				return
			}

			w.WriteHeader(http.StatusOK)
		}),
	)
	defer ts.Close()

	c := config.New()
	sc, err := c.ParseFile("../test/configs/test_scenario.yaml")
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	sc.Prom.TLS.CertFilePath = "../test/certs/dev.crt"
	sc.Prom.TLS.KeyFilePath = "../test/certs/dev.key"
	sc.Prom.TLS.InsecureSkipVerify = true
	sc.Prom.Headers = append(sc.Prom.Headers, config.HeaderData{Name: "foo", Value: "bar"})
	sc.Prom.JobName = "test"

	sc.Prom.PushURL = ts.URL

	pc := prom.New(sc)
	assert.NotNil(t, pc)

	pc.Requests(1, "req1", "resp1", "200")
	pc.Errors(0, "req1")

	start := time.Now()
	time.Sleep(time.Millisecond * 10)

	pc.Durations(start, 1, "req2", "GET", "resp3", "201")

	mfs, err := pc.Reg.Gather()
	assert.Nil(t, err)

	buf := &bytes.Buffer{}
	enc := expfmt.NewEncoder(buf, expfmt.NewFormat(expfmt.TypeProtoDelim))

	for _, mf := range mfs {
		err := enc.Encode(mf)
		assert.Nil(t, err)
	}
	wantBody := buf.Bytes()

	err = pc.Push()
	assert.Nil(t, err)

	assert.Equal(t, http.MethodPut, lastMethod)
	assert.Equal(t, "/metrics/job/test", lastPath)

	assert.True(t, bytes.Equal(wantBody, lastBody))

	assert.NotNil(t, lastHeader)
	assert.Contains(t, lastHeader, "Foo")
}

func TestNoPrometheus(t *testing.T) {

	c := config.New()
	sc, err := c.ParseFile("../test/configs/test_scenario.yaml")
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	pc := prom.New(sc)
	assert.NotNil(t, pc)

	assert.Nil(t, pc.Reg)

	pc.Requests(1, "req1", "resp1", "200")
	pc.Errors(0, "req1")

	start := time.Now()
	time.Sleep(time.Millisecond * 10)

	pc.Durations(start, 1, "req2", "GET", "resp3", "201")

	err = pc.Push()
	assert.Nil(t, err)
}
