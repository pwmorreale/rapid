//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

package logger_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/pwmorreale/rapid/internal/config"
	"github.com/pwmorreale/rapid/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {

	var b bytes.Buffer

	for _, test := range []struct {
		errText string
		opts    logger.Options
	}{
		{
			errText: "",
			opts:    logger.Options{Handler: "text", Timestamp: false, Level: "info", Writer: &b},
		},
		{
			errText: "",
			opts:    logger.Options{Handler: "json", Timestamp: false, Level: "info", Writer: &b},
		},
		{
			errText: "Unknown handler type: foo",
			opts:    logger.Options{Handler: "foo", Timestamp: false, Level: "info", Writer: &b},
		},
		{
			errText: "Missing Logger writer",
			opts:    logger.Options{Handler: "text", Timestamp: false, Level: "info", Writer: nil},
		},
		{
			errText: "slog: level string \"goo\": unknown name",
			opts:    logger.Options{Handler: "text", Timestamp: false, Level: "goo", Writer: &b},
		},
	} {
		err := logger.Init(&test.opts)
		if test.errText == "" {
			assert.Nil(t, err)
		} else {
			assert.ErrorContains(t, err, test.errText)
		}
	}
}

func TestTextlog(t *testing.T) {

	var b bytes.Buffer

	opts := logger.Options{
		Handler:   "text",
		Timestamp: false,
		Level:     "info",
		Writer:    &b,
	}

	err := logger.Init(&opts)
	assert.Nil(t, err)

	logger.Debug(nil, nil, "foo")
	actual, err := io.ReadAll(&b)
	assert.Nil(t, err)

	expected := []byte{}
	assert.Equal(t, expected, actual)

	b.Reset()
	logger.Info(nil, nil, "foo")
	actual, err = io.ReadAll(&b)
	assert.Nil(t, err)

	expected = []byte{0x1b, 0x5b, 0x39, 0x32, 0x6d, 0x49, 0x4e, 0x46, 0x1b, 0x5b, 0x30, 0x6d, 0x20, 0x66, 0x6f, 0x6f, 0xa}
	assert.Equal(t, expected, actual)

	b.Reset()
	logger.Warn(nil, nil, "foo")
	actual, err = io.ReadAll(&b)
	assert.Nil(t, err)

	expected = []byte{0x1b, 0x5b, 0x39, 0x33, 0x6d, 0x57, 0x52, 0x4e, 0x1b, 0x5b, 0x30, 0x6d, 0x20, 0x66, 0x6f, 0x6f, 0xa}
	assert.Equal(t, expected, actual)

	b.Reset()
	logger.Error(nil, nil, "foo")
	actual, err = io.ReadAll(&b)
	assert.Nil(t, err)

	expected = []byte{0x1b, 0x5b, 0x39, 0x31, 0x6d, 0x45, 0x52, 0x52, 0x1b, 0x5b, 0x30, 0x6d, 0x20, 0x66, 0x6f, 0x6f, 0xa}
	assert.Equal(t, expected, actual)

	req := &config.Request{
		Name:   "Goober",
		Method: "get",
	}

	rsp := &config.Response{
		Name:       "Goober-Response",
		StatusCode: 500,
	}

	b.Reset()
	logger.Error(req, rsp, "foo")
	actual, err = io.ReadAll(&b)
	assert.Nil(t, err)

	expected = []byte{0x1b, 0x5b, 0x39, 0x31, 0x6d, 0x45, 0x52, 0x52, 0x1b, 0x5b, 0x30, 0x6d, 0x20, 0x66, 0x6f, 0x6f, 0x20, 0x1b, 0x5b, 0x32, 0x6d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x6e, 0x61, 0x6d, 0x65, 0x3d, 0x1b, 0x5b, 0x30, 0x6d, 0x47, 0x6f, 0x6f, 0x62, 0x65, 0x72, 0x20, 0x1b, 0x5b, 0x32, 0x6d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x3d, 0x1b, 0x5b, 0x30, 0x6d, 0x67, 0x65, 0x74, 0x20, 0x1b, 0x5b, 0x32, 0x6d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x6e, 0x61, 0x6d, 0x65, 0x3d, 0x1b, 0x5b, 0x30, 0x6d, 0x47, 0x6f, 0x6f, 0x62, 0x65, 0x72, 0x2d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x20, 0x1b, 0x5b, 0x32, 0x6d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x3d, 0x1b, 0x5b, 0x30, 0x6d, 0x35, 0x30, 0x30, 0xa}
	assert.Equal(t, expected, actual)

}

func TestExitCode(t *testing.T) {

	var b bytes.Buffer

	opts := logger.Options{
		Handler:   "text",
		Timestamp: false,
		Level:     "info",
		Writer:    &b,
	}

	logger.Init(&opts)
	assert.Equal(t, 0, logger.ErrorCount())
	assert.Equal(t, 0, logger.InfoCount())
	assert.Equal(t, 0, logger.WarnCount())
	assert.Equal(t, 0, logger.DebugCount())

	logger.Error(nil, nil, "foo")
	assert.Equal(t, 1, logger.ErrorCount())
	assert.Equal(t, 0, logger.InfoCount())
	assert.Equal(t, 0, logger.WarnCount())
	assert.Equal(t, 0, logger.DebugCount())

	logger.Error(nil, nil, "foo")
	assert.Equal(t, 2, logger.ErrorCount())
	assert.Equal(t, 0, logger.InfoCount())
	assert.Equal(t, 0, logger.WarnCount())
	assert.Equal(t, 0, logger.DebugCount())

	// N.B. already have 2 errors reported
	for i := 0; i < 138; i++ {
		logger.Error(nil, nil, "foobar")
	}

	assert.Equal(t, 140, logger.ErrorCount())
}
