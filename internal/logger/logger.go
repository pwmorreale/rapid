//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package logger defines the log facility
package logger

import (
	"io"
	"log/slog"

	"github.com/lmittmann/tint"
)

// NewSanity creates a logger for sanity checks.
func NewSanity(w io.Writer) *slog.Logger {

	return slog.New(tint.NewHandler(w, &tint.Options{
		Level:   slog.LevelDebug,
		NoColor: false,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
}
