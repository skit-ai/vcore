package slog

import "context"

// Logger represents the methods supported by slog
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
	Error(err error, msg string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Debugf(format string, args ...any)
	Errorf(err error, format string, args ...any)
	WithTraceId(ctx context.Context) Logger
	WithFields(fields map[string]any) Logger
}
