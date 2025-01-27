package logger_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/pwmorreale/rapid/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {

	var b bytes.Buffer

	l := logger.NewSanity(&b)
	assert.NotNil(t, l)

	expected := []byte{0x1b, 0x5b, 0x39, 0x32, 0x6d, 0x49, 0x4e, 0x46, 0x1b, 0x5b, 0x30, 0x6d, 0x20, 0x66, 0x6f, 0x6f, 0xa}

	l.Info("foo")
	actual, err := io.ReadAll(&b)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

}
