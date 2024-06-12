package slog

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/skit-ai/vcore/env"
	"github.com/skit-ai/vcore/instruments"
)

type loggerWrapper struct {
	logger    log.Logger
	mutex     sync.Mutex
	sensitive bool // if sensitive is true, dont log the values
}

const defaultMsgKey = "msg"
const defaultErrKey = "error"

var (
	defaultLoggerWrapper *loggerWrapper
	logLevel             string
	logSensitive         bool
	callerDepth          int
)

func init() {
	logLevel = env.String("LOG_LEVEL", "info")
	logSensitive = false
	callerDepth = env.Int("LOG_CALLER_DEPTH", 4)
	defaultLoggerWrapper = newloggerWrapper(logLevel, logSensitive)
}

// NewLogger returns a new instance of Logger.
func NewLogger() Logger {
	return newloggerWrapper(logLevel, logSensitive)
}

func newloggerWrapper(logLevel string, sensitive bool) *loggerWrapper {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, levelFilter(logLevel))
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	logger = log.With(logger, "caller", log.Caller(callerDepth))

	return &loggerWrapper{
		logger:    logger,
		sensitive: sensitive,
	}
}

func levelFilter(logLevel string) level.Option {
	switch logLevel {
	case "debug":
		return level.AllowDebug()
	case "info":
		return level.AllowInfo()
	case "warn":
		return level.AllowWarn()
	case "error":
		return level.AllowError()
	// Invalid or no logLevel means all levels are allowed to be logged
	default:
		return level.AllowAll()
	}
}

func mapToSlice(m map[string]any) []any {
	var args []any
	for k, v := range m {
		args = append(args, k, v)
	}
	return args
}

// Info logs a line with level info using the loggerWrapper instance.
func (l *loggerWrapper) Info(msg string, args ...any) {
	if l.sensitive {
		args = make([]any, 0)
	}
	level.Info(log.With(l.logger, defaultMsgKey, msg)).Log(args...)
}

// Warn logs a line with level warn using the loggerWrapper instance.
func (l *loggerWrapper) Warn(msg string, args ...any) {
	if l.sensitive {
		args = make([]any, 0)
	}

	level.Warn(log.With(l.logger, defaultMsgKey, msg)).Log(args...)
}

// Debug logs a line with level debug using a loggerWrapper instance.
func (l *loggerWrapper) Debug(msg string, args ...any) {
	if l.sensitive {
		args = make([]any, 0)
	}

	level.Debug(log.With(l.logger, defaultMsgKey, msg)).Log(args...)
}

// Error logs a line with level error using a loggerWrapper instance.
// If err is not nil it adds only the msg string or vice-versa. Otherwise adds both.
func (l *loggerWrapper) Error(err error, msg string, args ...any) {
	if l.sensitive {
		args = make([]any, 0)
	}

	if err == nil {
		level.Error(log.With(l.logger, defaultMsgKey, msg)).Log(args...)
		return
	}

	if msg == "" {
		level.Error(log.With(l.logger, defaultErrKey, err.Error())).Log(args...)
		return
	}

	level.Error(log.With(l.logger, defaultMsgKey, msg, defaultErrKey, err.Error())).Log(args...)
}

// Infof logs a format line with level info using the loggerWrapper instance.
func (l *loggerWrapper) Infof(format string, args ...any) {
	if l.sensitive {
		args = make([]any, 0)
	}

	level.Info(l.logger).Log(defaultMsgKey, fmt.Sprintf(format, args...))
}

// Warnf logs a format line with level warn using the loggerWrapper instance.
func (l *loggerWrapper) Warnf(format string, args ...any) {
	if l.sensitive {
		args = make([]any, 0)
	}

	level.Warn(l.logger).Log(defaultMsgKey, fmt.Sprintf(format, args...))
}

// Debugf logs a format line with level debug using the loggerWrapper instance.
func (l *loggerWrapper) Debugf(format string, args ...any) {
	if l.sensitive {
		args = make([]any, 0)
	}

	level.Debug(l.logger).Log(defaultMsgKey, fmt.Sprintf(format, args...))
}

// Errorf logs a format line with level error using a loggerWrapper instance.
// If err is not nil it adds only the msg string or vice-versa. Otherwise adds both.
func (l *loggerWrapper) Errorf(err error, format string, args ...any) {
	if l.sensitive {
		args = make([]any, 0)
	}

	if err == nil {
		level.Error(l.logger).Log(defaultMsgKey, fmt.Sprintf(format, args...))
		return
	}

	if format == "" {
		level.Error(l.logger).Log(defaultErrKey, err.Error())
		return
	}

	level.Error(l.logger).Log(defaultMsgKey, fmt.Sprintf(format, args...), defaultErrKey, err.Error())
}

// WithTraceId returns a pointer to updated loggerWrapper with trace_id attached to the logger.
func (l *loggerWrapper) WithTraceId(ctx context.Context) Logger {
	traceId := instruments.ExtractTraceID(ctx)
	logger := log.With(l.logger, "trace_id", traceId)

	return &loggerWrapper{
		logger: logger,
	}
}

// WithFields returns a pointer to updated loggerWrapper with custom fields attached to the logger.
func (l *loggerWrapper) WithFields(fields map[string]any) Logger {
	fieldArgs := mapToSlice(fields)
	logger := log.With(l.logger, fieldArgs...)

	return &loggerWrapper{
		logger: logger,
	}
}

// WithSensitive returns a pointer to updated loggerWrapper with sensitive flag updated to the logger.
func (l *loggerWrapper) WithSensitive(sensitive bool) Logger {
	logger := log.With(l.logger)

	return &loggerWrapper{
		logger:    logger,
		sensitive: sensitive,
	}
}

func (l *loggerWrapper) SetSensitive(val bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.sensitive = val
}

// Info logs a line with level Info.
func Info(msg string, args ...any) {
	if defaultLoggerWrapper.sensitive {
		args = make([]any, 0)
	}
	level.Info(log.With(defaultLoggerWrapper.logger, defaultMsgKey, msg)).Log(args...)
}

// Warn logs a line with level warn.
func Warn(msg string, args ...any) {
	if defaultLoggerWrapper.sensitive {
		args = make([]any, 0)
	}
	level.Warn(log.With(defaultLoggerWrapper.logger, defaultMsgKey, msg)).Log(args...)
}

// Debug logs a line with level debug.
func Debug(msg string, args ...any) {
	if defaultLoggerWrapper.sensitive {
		args = make([]any, 0)
	}
	level.Debug(log.With(defaultLoggerWrapper.logger, defaultMsgKey, msg)).Log(args...)
}

// Error logs a line with level error.
// If err is not nil it adds only the msg string or vice-versa. Otherwise adds both.
func Error(err error, msg string, args ...any) {
	if defaultLoggerWrapper.sensitive {
		args = make([]any, 0)
	}

	if err == nil {
		level.Error(log.With(defaultLoggerWrapper.logger, defaultMsgKey, msg)).Log(args...)
		return
	}

	if msg == "" {
		level.Error(log.With(defaultLoggerWrapper.logger, defaultErrKey, err.Error())).Log(args...)
		return
	}

	level.Error(log.With(defaultLoggerWrapper.logger, defaultMsgKey, msg, defaultErrKey, err.Error())).Log(args...)
}

// Infof logs a format line with level info.
func Infof(format string, args ...any) {
	if defaultLoggerWrapper.sensitive {
		args = make([]any, 0)
	}

	level.Info(defaultLoggerWrapper.logger).Log(defaultMsgKey, fmt.Sprintf(format, args...))
}

// Warnf logs a format line with level warn.
func Warnf(format string, args ...any) {
	if defaultLoggerWrapper.sensitive {
		args = make([]any, 0)
	}

	level.Warn(defaultLoggerWrapper.logger).Log(defaultMsgKey, fmt.Sprintf(format, args...))
}

// Debugf logs a format line with level debug.
func Debugf(format string, args ...any) {
	if defaultLoggerWrapper.sensitive {
		args = make([]any, 0)
	}

	level.Debug(defaultLoggerWrapper.logger).Log(defaultMsgKey, fmt.Sprintf(format, args...))
}

// Errorf logs a format line with level error.
// If err is not nil it adds only the msg string or vice-versa. Otherwise adds both.
func Errorf(err error, format string, args ...any) {
	if defaultLoggerWrapper.sensitive {
		args = make([]any, 0)
	}

	if err == nil {
		level.Error(defaultLoggerWrapper.logger).Log(defaultMsgKey, fmt.Sprintf(format, args...))
		return
	}

	if format == "" {
		level.Error(defaultLoggerWrapper.logger).Log(defaultErrKey, err.Error())
		return
	}

	level.Error(defaultLoggerWrapper.logger).Log(defaultMsgKey, fmt.Sprintf(format, args...), defaultErrKey, err.Error())
}

// WithFields returns WithFields using the defaultLoggerWrapper.
func WithFields(fields map[string]any) Logger {
	return defaultLoggerWrapper.WithFields(fields)
}

// WithTraceId returns WithTraceId using the defaultLoggerWrapper.
func WithTraceId(ctx context.Context) Logger {
	return defaultLoggerWrapper.WithTraceId(ctx)
}

// DefaultLogger returns the default instance of logfmt logger
func DefaultLogger() log.Logger {
	return defaultLoggerWrapper.logger
}
