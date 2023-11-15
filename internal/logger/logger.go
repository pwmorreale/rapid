package logger

import (
	"os"
	"time"

	"github.com/pwmorreale/rapid/internal/config"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func GetLogger(label string) zerolog.Logger {

	var level zerolog.Level

	p := config.LogLevel + "." + label

	v := viper.GetString(p)
	level, err := zerolog.ParseLevel(v)
	if err == nil {
		level = zerolog.InfoLevel
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(level).
		With().
		Timestamp().
		Logger()

	return logger
}
