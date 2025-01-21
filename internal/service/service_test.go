package service_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/data"
	"github.com/pwmorreale/rapid/internal/service"
	"github.com/test-go/testify/assert"
)

var testURL = "https://bob_ross.com/happy_little_trees"
var testCookie = `id=bob_ross; Max-Age=42; SameSite=Strict; id=betsy_ross; Expires="Thu, 21 Oct 2080 07:28:00 GMT"; SameSite=Strict`

func initTestService(t *testing.T) (*service.Context, *config.Scenario, error) {
	c := config.New()
	sc, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	d := data.New()
	for k, v := range sc.Data {
		err := d.AddReplacement(k, v)
		if err != nil {
			return nil, nil, err
		}

	}

	s := service.New(d)

	return s, sc, nil
}

func TestCreateRequest(t *testing.T) {

	s, sc, err := initTestService(t)
	assert.NotNil(t, s)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	request, err := s.CreateRequest(&sc.Sequence.Requests[0])
	assert.Nil(t, err)

	assert.Equal(t, testURL, request.URL.String())

	// N.B.  Canonical form for the key vs. original in yaml...
	assert.Contains(t, request.Header, "X-Paintbrush-For-Sky")
	assert.Contains(t, request.Header, "Content-Length")

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

func TestCreateClient(t *testing.T) {

	s, sc, err := initTestService(t)
	assert.NotNil(t, s)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	client, err := s.CreateClient(&sc.Sequence.Requests[0])
	assert.Nil(t, err)
	assert.NotNil(t, client)

}

func TestHeaderMultipleValues(t *testing.T) {

	s, sc, err := initTestService(t)
	assert.NotNil(t, s)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	response := &http.Response{}
	response.Header = make(map[string][]string)

	configResponse := &sc.Sequence.Requests[0].Responses[0]

	response.Header.Add("header1", "value1")
	response.Header.Add("header2", "value2")
	response.Header.Add("header3", "value3")

	err = s.VerifyHeaders(response, configResponse)
	assert.Nil(t, err)
}

func TestHeaderMissingHeader(t *testing.T) {

	s, sc, err := initTestService(t)
	assert.NotNil(t, s)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	response := &http.Response{}
	response.Header = make(map[string][]string)

	configResponse := &sc.Sequence.Requests[0].Responses[0]

	err = s.VerifyHeaders(response, configResponse)
	assert.Equal(t, "header: Header1 not found", err.Error())

}

func TestHeaderBadValue(t *testing.T) {

	s, sc, err := initTestService(t)
	assert.NotNil(t, s)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	response := &http.Response{}
	response.Header = make(map[string][]string)

	response.Header.Add("header1", "foobar")

	configResponse := &sc.Sequence.Requests[0].Responses[0]

	err = s.VerifyHeaders(response, configResponse)
	assert.Equal(t, "header: Header1, expected value (value1) not found", err.Error())

}

func TestHeaders(t *testing.T) {

	s, sc, err := initTestService(t)
	assert.NotNil(t, s)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	response := &http.Response{}

	configResponse := &sc.Sequence.Requests[0].Responses[0]

	// Missing...
	err = s.VerifyHeaders(response, configResponse)
	assert.NotNil(t, err)

	response.Header = make(map[string][]string)
	response.Header.Add(configResponse.Headers[0].Name, configResponse.Headers[0].Value)
	response.Header.Add(configResponse.Headers[1].Name, configResponse.Headers[1].Value)

	err = s.VerifyHeaders(response, configResponse)
	assert.Nil(t, err)

}

func TestVerifyCookies(t *testing.T) {

	s, sc, err := initTestService(t)
	assert.NotNil(t, s)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	response := &http.Response{}

	configResponse := &sc.Sequence.Requests[0].Responses[0]

	response.Header = make(map[string][]string)
	response.Header.Add("Set-Cookie", configResponse.Cookies[0].Value)

	err = s.VerifyCookies(response, configResponse)
	assert.Equal(t, "cookie: id=marion_ross; Expires=Mon, 21 Oct 2080 07:28:00 GMT not found in response", err.Error())

	response.Header.Add("Set-Cookie", configResponse.Cookies[1].Value)
	err = s.VerifyCookies(response, configResponse)
	assert.Nil(t, err)

	// Typo in config...
	configResponse.Cookies[0].Value = "foo"

	err = s.VerifyCookies(response, configResponse)
	assert.Equal(t, "http: '=' not found in cookie", err.Error())

}
