package service_test

import (
	"testing"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/service"
	"github.com/test-go/testify/assert"
)

var testURL = "https://bob-ross:happy-little-trees@google.com:80/blah/moo?key=value&key2=value2#fraggie"

func TestCreateRequest(t *testing.T) {

	s := service.New()

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/single_request.yaml")
	assert.Nil(t, err)

	request, err := s.CreateRequest(&sc.Seq.Reqs[0])
	assert.Nil(t, err)

	assert.Equal(t, testURL, request.URL.String())

	assert.Contains(t, request.Header, "Foo")
	assert.Contains(t, request.Header, "Content-Length")
}

func TestCreateClient(t *testing.T) {

	s := service.New()

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/single_request.yaml")
	assert.Nil(t, err)

	client, err := s.CreateClient(&sc.Seq.Reqs[0])
	assert.Nil(t, err)
	assert.NotNil(t, client)

}
