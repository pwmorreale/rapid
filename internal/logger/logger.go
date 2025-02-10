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

// Handler is used to define a text or json handler.
type Handler int

const (

	// Text specifies using the tinted text handler.
	Text Handler = 1

	// JSON specifies using the slog JSON handler.
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

// Options defines options for the new logger instance.
type Options struct {
	Handler       Handler
	OmitTimestamp bool
	DefaultLevel  slog.Level
	Writer        io.Writer
}

// Context defines a context for the package.
type Context struct {
	logHandle *slog.Logger
}

func omitTimestamp(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}
	return a
}

// New creates a new logger instance based on the options.
func New(opts Options) *Context {

	c := &Context{}

	replaceAttr := omitTimestamp
	if !opts.OmitTimestamp {
		replaceAttr = nil
	}

	switch opts.Handler {
	case Text:
		topts := &tint.Options{
			Level:       opts.DefaultLevel,
			ReplaceAttr: replaceAttr,
		}
		c.logHandle = slog.New(tint.NewHandler(opts.Writer, topts))
	case JSON:
		hopts := &slog.HandlerOptions{
			Level:       opts.DefaultLevel,
			ReplaceAttr: replaceAttr,
		}
		c.logHandle = slog.New(slog.NewJSONHandler(opts.Writer, hopts))
	default:
		return nil
	}

	return c

}

func (c *Context) handleLog(level slog.Level, req *config.Request, rsp *config.Response, format string, args ...any) {

	if !c.logHandle.Enabled(context.Background(), level) {
		return
	}

	s := format
	if len(args) > 1 {
		s = fmt.Sprintf(format, args...)
	}

	r := slog.NewRecord(time.Now(), level, s, 0)

	if req != nil {
		r.Add("request", req)
	}
	if rsp != nil {
		r.Add("response", rsp)
	}

	c.logHandle.Handler().Handle(context.Background(), r)
}

// Debug writes a debug log message
func (c *Context) Debug(req *config.Request, rsp *config.Response, format string, args ...any) {
	c.handleLog(slog.LevelDebug, req, rsp, format, args)
}

// Info writes an info log message
func (c *Context) Info(req *config.Request, rsp *config.Response, format string, args ...any) {
	c.handleLog(slog.LevelInfo, req, rsp, format, args)
}

// Warn writes a warn log message
func (c *Context) Warn(req *config.Request, rsp *config.Response, format string, args ...any) {
	c.handleLog(slog.LevelWarn, req, rsp, format, args)
}

// Error writes an error log message
func (c *Context) Error(req *config.Request, rsp *config.Response, format string, args ...any) {
	c.handleLog(slog.LevelError, req, rsp, format, args)
}
