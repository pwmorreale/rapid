//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package rest executes the REST calls
package rest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/pwmorreale/rapid/config"
	"github.com/pwmorreale/rapid/data"
	"github.com/pwmorreale/rapid/logger"
	"github.com/stretchr/testify/assert"
)

var testURL = "https://bob_ross.com/happy_little_trees"
var testCookie = `id=bob_ross; Max-Age=42; SameSite=Strict; id=betsy_ross; Expires="Thu, 21 Oct 2080 07:28:00 GMT"; SameSite=Strict`

var json = `{
  "foo": "barhoo",
  "goo": {
    "moo": {
      "boo": "doo"
    }
  }
}`

var xml = `<?xml version="1.0" encoding="UTF-8"?>
<note>
  <from>Bob Ross</from>
  <to>All Painters</to>
  <message>Remember to paint this weekend!</message>
</note>`

// Testing transport...

type TestingTransport struct {
	Response *http.Response
	Error    error
}

func (m *TestingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return m.Response, m.Error
}

func makeResponse(statusCode int, contentType string, content []byte, length int64, headers *[]config.HeaderData, cookies *[]config.CookieData) *http.Response {

	response := http.Response{}

	response.StatusCode = statusCode
	response.Status = fmt.Sprintf("%d", statusCode)

	response.Body = io.NopCloser(bytes.NewReader(content))

	response.Header = make(map[string][]string)

	if contentType != "" {
		response.Header.Add("Content-Type", contentType)
	}
	response.Header.Add("Content-Length", strconv.FormatInt(length, 10))

	response.ContentLength = length

	if headers != nil {
		for i := range *headers {
			response.Header.Add((*headers)[i].Name, (*headers)[i].Value)
		}
	}

	if cookies != nil {
		for i := range *cookies {
			response.Header.Add("Set-Cookie", (*cookies)[i].Value)
		}
	}

	return &response
}

func makeResponseFromResponse(response *config.Response, content []byte) *http.Response {

	return makeResponse(response.StatusCode, response.Content.MediaType, content, int64(len(content)), &response.Headers, &response.Cookies)
}

func initTestService(t *testing.T) (*Context, *config.Scenario, data.Data, error) {
	c := config.New()
	sc, err := c.ParseFile("../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	d := data.New()
	for i := range sc.Replacements {
		err := d.AddReplacement(sc.Replacements[i].Regex, sc.Replacements[i].Value)
		if err != nil {
			return nil, nil, nil, err
		}

	}

	r := New(sc, d)

	return r, sc, d, nil
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

func TestExecute(t *testing.T) {

	r, sc, _, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	initLogger(os.Stdout)

	httpResponse := makeResponseFromResponse(sc.Sequence.Requests[0].Responses[0], []byte(json))

	r.mockRoundTripper = &TestingTransport{
		Response: httpResponse,
		Error:    nil,
	}

	ctx := context.Background()

	r.Execute(ctx, &sc.Sequence.Requests[0])
}
func TestCreateRequest(t *testing.T) {

	r, sc, _, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	ctx := context.Background()

	request, err := r.createRequest(ctx, &sc.Sequence.Requests[0])
	assert.Nil(t, err)
	assert.Equal(t, testURL, request.URL.String())

	// N.B.  Canonical form for the key vs. original in yaml...
	assert.Contains(t, request.Header, "X-Paintbrush-For-Sky")

	v := request.Header.Get("X-Paintbrush-For-Sky")
	assert.Equal(t, "wide", v)

	v = request.Header.Get("Cookie")
	assert.Equal(t, testCookie, v)

	// Should have a substitution...
	ior, err := request.GetBody()
	assert.Nil(t, err)
	contents, err := io.ReadAll(ior)
	assert.Nil(t, err)
	assert.Equal(t, "various paint colors in blue", string(contents))

}

func TestHeaderMultipleValues(t *testing.T) {

	r, sc, _, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	headers := []config.HeaderData{{Name: "header1", Value: "value1"},
		{Name: "header2", Value: "value2"},
		{Name: "header3", Value: "value3"}}
	response := makeResponse(200, "", []byte{}, 0, &headers, nil)

	configResponse := sc.Sequence.Requests[0].Responses[0]

	err = r.verifyHeaders(response, configResponse)
	assert.Nil(t, err)
}

func TestHeaderMissingHeader(t *testing.T) {

	r, sc, _, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	headers := []config.HeaderData{}
	response := makeResponse(200, "", []byte{}, 0, &headers, nil)

	configResponse := sc.Sequence.Requests[0].Responses[0]

	err = r.verifyHeaders(response, configResponse)
	assert.Equal(t, "header: Header1 not found", err.Error())

}

func TestHeaderBadValue(t *testing.T) {

	r, sc, _, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	headers := []config.HeaderData{{Name: "header1", Value: "foobar"}}
	response := makeResponse(200, "", []byte{}, 0, &headers, nil)

	configResponse := sc.Sequence.Requests[0].Responses[0]

	err = r.verifyHeaders(response, configResponse)
	assert.Equal(t, "header: Header1, expected value (value1) not found", err.Error())

}

func TestHeaders(t *testing.T) {

	r, sc, _, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	response := makeResponse(200, "", []byte{}, 0, nil, nil)

	configResponse := sc.Sequence.Requests[0].Responses[0]

	// Missing...
	err = r.verifyHeaders(response, configResponse)
	assert.NotNil(t, err)

	response.Header = make(map[string][]string)
	response.Header.Add(configResponse.Headers[0].Name, configResponse.Headers[0].Value)
	response.Header.Add(configResponse.Headers[1].Name, configResponse.Headers[1].Value)

	err = r.verifyHeaders(response, configResponse)
	assert.Nil(t, err)

}

func TestVerifyCookies(t *testing.T) {

	r, sc, _, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	configResponse := sc.Sequence.Requests[0].Responses[0]

	response := &http.Response{}

	response.Header = make(map[string][]string)
	response.Header.Add("Set-Cookie", configResponse.Cookies[0].Value)

	err = r.verifyCookies(response, configResponse)
	assert.Equal(t, "cookie: id=marion_ross; Expires=Mon, 21 Oct 2080 07:28:00 GMT not found in response", err.Error())

	response.Header.Add("Set-Cookie", configResponse.Cookies[1].Value)
	err = r.verifyCookies(response, configResponse)
	assert.Nil(t, err)

	// Typo in config...
	configResponse.Cookies[0].Value = "foo"

	err = r.verifyCookies(response, configResponse)
	assert.Equal(t, "http: '=' not found in cookie", err.Error())

}

func TestVerifyNoContent(t *testing.T) {

	r, sc, _, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	response := makeResponse(200, "", []byte{}, 0, nil, nil)

	configResponse := sc.Sequence.Requests[0].Responses[1]

	// No content...
	err = r.verifyContentAndExtract(response, configResponse)
	assert.Nil(t, err)

}

func TestVerifyJSONContent(t *testing.T) {

	r, sc, d, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	response := makeResponse(200, "application/json", []byte(json), -1, nil, nil)
	configResponse := sc.Sequence.Requests[0].Responses[0]

	err = r.verifyContentAndExtract(response, configResponse)
	assert.Nil(t, err)

	assert.Equal(t, "doo", d.Lookup("foo"))
}

func TestVerifyXMLContent(t *testing.T) {

	r, sc, d, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	response := makeResponse(200, "text/xml", []byte(xml), int64(len(xml)), nil, nil)
	configResponse := sc.Sequence.Requests[0].Responses[2]

	err = r.verifyContentAndExtract(response, configResponse)
	assert.Nil(t, err)

	assert.Equal(t, "Bob Ross", d.Lookup("who"))
}

func TestFindResponse(t *testing.T) {

	r, sc, _, err := initTestService(t)
	assert.NotNil(t, r)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	httpResponse := makeResponse(501, "", []byte(""), 0, nil, nil)
	assert.NotNil(t, httpResponse)

	request := &sc.Sequence.Requests[1]
	assert.Empty(t, request.UnknownResponses)

	resp := r.findResponse(httpResponse, request)
	assert.NotNil(t, resp)

	assert.Equal(t, 1, len(request.UnknownResponses))
	assert.Equal(t, 501, request.UnknownResponses[0].StatusCode)

}
