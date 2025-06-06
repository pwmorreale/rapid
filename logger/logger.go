//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package logger defines the log facility
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lmittmann/tint"
	"github.com/pwmorreale/rapid/config"
)

// Handler is used to define a text or json handler.
type Handler int

var logHandle *slog.Logger
var errCount atomic.Int32
var infoCount atomic.Int32
var warnCount atomic.Int32
var debugCount atomic.Int32

// Options defines options for the new logger instance.
type Options struct {
	Handler   string
	Timestamp bool
	Level     string
	Writer    io.Writer
}

func omitTimestamp(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}
	return a
}

// Init creates a new logger instance based on the options.
func Init(opts *Options) error {

	errCount.Store(0)
	infoCount.Store(0)
	warnCount.Store(0)
	debugCount.Store(0)

	replaceAttr := omitTimestamp
	if opts.Timestamp {
		replaceAttr = nil
	}

	if opts.Writer == nil {
		return fmt.Errorf("missing Logger writer")
	}

	var sl slog.Level
	err := sl.UnmarshalText([]byte(opts.Level))
	if err != nil {
		return err
	}

	switch strings.ToLower(opts.Handler) {
	case "text":
		topts := &tint.Options{
			Level:       sl,
			ReplaceAttr: replaceAttr,
		}
		logHandle = slog.New(tint.NewHandler(opts.Writer, topts))
	case "json":
		hopts := &slog.HandlerOptions{
			Level:       sl,
			ReplaceAttr: replaceAttr,
		}
		logHandle = slog.New(slog.NewJSONHandler(opts.Writer, hopts))
	default:
		return fmt.Errorf("unknown handler type: %v", opts.Handler)
	}

	return nil
}

func handleLog(level slog.Level, req *config.Request, rsp *config.Response, format string, args ...any) {

	if !logHandle.Enabled(context.Background(), level) {
		return
	}

	s := format
	if len(args) > 0 {
		s = fmt.Sprintf(format, args...)
	}

	r := slog.NewRecord(time.Now(), level, s, 0)

	if req != nil {
		r.Add("request", req)
	}
	if rsp != nil {
		r.Add("response", rsp)
	}

	logHandle.Handler().Handle(context.Background(), r)
}

// Debug writes a debug log message
func Debug(req *config.Request, rsp *config.Response, format string, args ...any) {
	handleLog(slog.LevelDebug, req, rsp, format, args...)
	debugCount.Add(1)
}

// Info writes an info log message
func Info(req *config.Request, rsp *config.Response, format string, args ...any) {
	handleLog(slog.LevelInfo, req, rsp, format, args...)
	infoCount.Add(1)
}

// Warn writes a warn log message
func Warn(req *config.Request, rsp *config.Response, format string, args ...any) {
	handleLog(slog.LevelWarn, req, rsp, format, args...)
	warnCount.Add(1)
}

// Error writes an error log message
func Error(req *config.Request, rsp *config.Response, format string, args ...any) {
	handleLog(slog.LevelError, req, rsp, format, args...)
	errCount.Add(1)
}

// DebugCount is the number of debug logs
func DebugCount() int {
	return int(debugCount.Load())
}

// InfoCount is the number of info logs
func InfoCount() int {
	return int(infoCount.Load())
}

// WarnCount is the number of warn logs
func WarnCount() int {
	return int(warnCount.Load())
}

// ErrorCount is the number of error logs
func ErrorCount() int {
	return int(errCount.Load())
}
