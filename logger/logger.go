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
	"sync"
	"sync/atomic"
	"time"

	"github.com/lmittmann/tint"
	"github.com/pwmorreale/rapid/config"
)

// Handler is used to define a text or json handler.
type Handler int

// instance holds all logger state so that Init creates a fresh isolated context.
type instance struct {
	handle     *slog.Logger
	errCount   atomic.Int32
	infoCount  atomic.Int32
	warnCount  atomic.Int32
	debugCount atomic.Int32
}

var (
	current *instance
	mu      sync.RWMutex
)

func getInst() *instance {
	mu.RLock()
	defer mu.RUnlock()
	return current
}

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

	if opts.Writer == nil {
		return fmt.Errorf("missing Logger writer")
	}

	replaceAttr := omitTimestamp
	if opts.Timestamp {
		replaceAttr = nil
	}

	var sl slog.Level
	err := sl.UnmarshalText([]byte(opts.Level))
	if err != nil {
		return err
	}

	inst := &instance{}

	switch strings.ToLower(opts.Handler) {
	case "text":
		topts := &tint.Options{
			Level:       sl,
			ReplaceAttr: replaceAttr,
		}
		inst.handle = slog.New(tint.NewHandler(opts.Writer, topts))
	case "json":
		hopts := &slog.HandlerOptions{
			Level:       sl,
			ReplaceAttr: replaceAttr,
		}
		inst.handle = slog.New(slog.NewJSONHandler(opts.Writer, hopts))
	default:
		return fmt.Errorf("unknown handler type: %v", opts.Handler)
	}

	mu.Lock()
	current = inst
	mu.Unlock()

	return nil
}

func handleLog(level slog.Level, req *config.Request, rsp *config.Response, format string, args ...any) {

	inst := getInst()

	if !inst.handle.Enabled(context.Background(), level) {
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

	inst.handle.Handler().Handle(context.Background(), r)
}

// Debug writes a debug log message
func Debug(req *config.Request, rsp *config.Response, format string, args ...any) {
	handleLog(slog.LevelDebug, req, rsp, format, args...)
	getInst().debugCount.Add(1)
}

// Info writes an info log message
func Info(req *config.Request, rsp *config.Response, format string, args ...any) {
	handleLog(slog.LevelInfo, req, rsp, format, args...)
	getInst().infoCount.Add(1)
}

// Warn writes a warn log message
func Warn(req *config.Request, rsp *config.Response, format string, args ...any) {
	handleLog(slog.LevelWarn, req, rsp, format, args...)
	getInst().warnCount.Add(1)
}

// Error writes an error log message
func Error(req *config.Request, rsp *config.Response, format string, args ...any) {
	handleLog(slog.LevelError, req, rsp, format, args...)
	getInst().errCount.Add(1)
}

// DebugCount is the number of debug logs
func DebugCount() int {
	return int(getInst().debugCount.Load())
}

// InfoCount is the number of info logs
func InfoCount() int {
	return int(getInst().infoCount.Load())
}

// WarnCount is the number of warn logs
func WarnCount() int {
	return int(getInst().warnCount.Load())
}

// ErrorCount is the number of error logs
func ErrorCount() int {
	return int(getInst().errCount.Load())
}
