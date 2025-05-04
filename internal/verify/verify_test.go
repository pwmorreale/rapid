//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package verify_test contains unit tests for the verify module.
package verify_test

import (
	"io"
	"testing"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/logger"
	"github.com/pwmorreale/rapid/internal/verify"
	"github.com/stretchr/testify/assert"
)

func initLogger(wr io.Writer) {

	opts := logger.Options{
		Handler:   "text",
		Timestamp: false,
		Level:     "Info",
		Writer:    wr,
	}

	logger.Init(&opts)
}

func TestCheckRequestCookiesGood(t *testing.T) {

	Cookie := config.CookieData{
		Value: "a=b; c=d; e=f",
	}

	request := &config.Request{
		Cookies: []config.CookieData{Cookie},
	}

	initLogger(io.Discard)

	verify.CheckCookies(request, nil, request.Cookies)
	assert.Equal(t, 0, logger.ErrorCount())
	assert.Equal(t, 4, logger.DebugCount())
}

func TestCheckRequestCookiesBad(t *testing.T) {

	Cookie := config.CookieData{
		Value: `a="a b";`,
	}

	request := &config.Request{
		Cookies: []config.CookieData{Cookie},
	}

	initLogger(io.Discard)

	verify.CheckCookies(request, nil, request.Cookies)
	assert.Equal(t, 1, logger.ErrorCount())
	assert.Equal(t, 0, logger.DebugCount())
}

func TestRequestContent(t *testing.T) {

	request := &config.Request{
		ContentType: "text/plain",
		Content:     "some text",
	}

	initLogger(io.Discard)

	verify.CheckRequestContent(request)
	assert.Equal(t, 0, logger.ErrorCount())
	assert.Equal(t, 0, logger.WarnCount())

}

func TestCheck(t *testing.T) {

	initLogger(io.Discard)

	err := verify.Check("../../test/configs/verify_test.yaml")
	assert.Nil(t, err)

	assert.Equal(t, 13, logger.ErrorCount())
	assert.Equal(t, 4, logger.WarnCount())
	assert.Equal(t, 11, logger.InfoCount())
	assert.Equal(t, 9, logger.DebugCount())
}
