// Package dblog provides structured logging to the database with stack traces.
package dblog

import (
	"log/slog"
	"runtime/debug"
)

// Logger is the application-wide slog instance.
// Set it once in main.go after creating the multi-handler logger.
var Logger *slog.Logger

// Error logs an error with its message AND full stack trace.
// Use this for ALL application errors — never use bare slog.Error.
func Error(msg string, err error, attrs ...slog.Attr) {
	if Logger == nil {
		return
	}
	stack := string(debug.Stack())
	allAttrs := append([]slog.Attr{
		slog.String("error", err.Error()),
		slog.String("stack", stack),
	}, attrs...)
	Logger.LogAttrs(nil, slog.LevelError, msg, allAttrs...)
}
