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
	sc, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	fmt.Println(err)
	assert.Nil(t, err)

	request, err := s.CreateRequest(&sc.Sequence.Requests[0])
	assert.Nil(t, err)

	assert.Equal(t, testURL, request.URL.String())

	assert.Contains(t, request.Header, "X-Paintbrush-For-Sky")
	assert.Contains(t, request.Header, "Content-Length")

}

func TestCreateClient(t *testing.T) {

	s := service.New(data.New())

	c := config.New()
	sc, err := c.ParseFile("../../test/configs/test_scenario.yaml")
	assert.Nil(t, err)

	client, err := s.CreateClient(&sc.Sequence.Requests[0])
	assert.Nil(t, err)
	assert.NotNil(t, client)

}
