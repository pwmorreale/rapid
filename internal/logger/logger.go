//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package logger defines the log facility
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/lmittmann/tint"
	"github.com/pwmorreale/rapid/internal/config"
)

type Handler int

const (
	Text Handler = 1
	JSON Handler = 2
)

// Logger defines the interface for logging
//
//go:generate counterfeiter -o ../../test/mocks/fake_logger.go . Logger
type Logger interface {
	Debug(config.Request, *config.Response, string, ...any)
	Info(config.Request, *config.Response, string, ...any)
	Warn(config.Request, *config.Response, string, ...any)
	Error(config.Request, *config.Response, string, ...any)
}

type Options struct {
	handler      Handler
	defaultLevel slog.Level
	w            io.Writer
}

// COntext defines a context for the package.
type Context struct {
	nrErr     int
	logHandle *slog.Logger
}

func makeTextLogger(opts Options) *slog.Logger {
	hopts := &tint.Options{
		Level: opts.defaultLevel,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}

	return slog.New(tint.NewHandler(opts.w, hopts))
}

func makeJSONLogger(opts Options) *slog.Logger {
	hopts := &slog.HandlerOptions{
		Level: opts.defaultLevel,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}

	return slog.New(slog.NewJSONHandler(opts.w, hopts))
}

func New(opts Options) *Context {

	c := &Context{}

	switch opts.handler {
	case Text:
		c.logHandle = makeTextLogger(opts)
	case JSON:
		c.logHandle = makeJSONLogger(opts)
	default:
		return nil
	}

	return &Context{}
}

func (c *Context) handleLog(level slog.Level, reg config.Request, rsp *config.Response, format string, args ...any) {
	if !c.logHandle.Enabled(context.Background(), slog.LevelDebug) {
		return
	}

	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), 0)
	_ = c.logHandle.Handler().Handle(context.Background(), r)
}

// Debug writes a debug log message
func (c *Context) Debug(req config.Request, rsp *config.Response, format string, args ...any) {
	c.handleLog(slog.LevelDebug, req, rsp, format, args)
}

// Info writes an info log message
func (c *Context) Info(req config.Request, rsp *config.Response, format string, args ...any) {
	c.handleLog(slog.LevelInfo, req, rsp, format, args)
}

// Warn writes a warn log message
func (c *Context) Warn(req config.Request, rsp *config.Response, format string, args ...any) {
	c.handleLog(slog.LevelWarn, req, rsp, format, args)
}

// Error writes an error log message
func (c *Context) Error(req config.Request, rsp *config.Response, format string, args ...any) {
	c.handleLog(slog.LevelError, req, rsp, format, args)
}
