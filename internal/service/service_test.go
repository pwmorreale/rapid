package service_test

import (
	"fmt"
	"testing"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/data"
	"github.com/pwmorreale/rapid/internal/service"
	"github.com/test-go/testify/assert"
)

var testURL = "https://bob.ross.com/happy_little_trees"

func TestCreateRequests(t *testing.T) {

	s := service.New(data.New())

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/single_request.yaml")
	fmt.Println(err)
	assert.Nil(t, err)

	request, err := s.CreateRequest(&sc.Sequence.Requests[0])
	assert.Nil(t, err)

	assert.Equal(t, testURL, request.URL.String())

	assert.Contains(t, request.Header, "Foo")
	assert.Contains(t, request.Header, "Content-Length")
}

func TestCreateClient(t *testing.T) {

	s := service.New(data.New())

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/single_request.yaml")
	assert.Nil(t, err)

	client, err := s.CreateClient(&sc.Sequence.Requests[0])
	assert.Nil(t, err)
	assert.NotNil(t, client)

}

func TestSend(t *testing.T) {

	err := service.checkContains()

}
