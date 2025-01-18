package service_test

import (
	"fmt"
	"io"
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
	fmt.Println(err)
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
