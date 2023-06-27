package slog

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/skit-ai/vcore/env"
	"github.com/skit-ai/vcore/instruments"
)

type loggerWrapper struct {
    logger log.Logger
}

var defaultLoggerWrapper *loggerWrapper

func init() {
	logLevel := env.String("LOG_LEVEL", "info")
	defaultLoggerWrapper = newloggerWrapper(logLevel)
}

func newloggerWrapper(logLevel string) *loggerWrapper {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, levelFilter(logLevel))
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	logger = log.With(logger, "caller", log.Caller(4))

	return &loggerWrapper{
		logger: logger,
	}
}

func (l *loggerWrapper) Info(msg string, args ...any) {
    level.Info(log.With(l.logger, "msg", msg)).Log(args...)
}

func (l *loggerWrapper) Warn(msg string, args ...any) {
    level.Warn(log.With(l.logger, "msg", msg)).Log(args...)
}

func (l *loggerWrapper) Debug(msg string, args ...any) {
    level.Debug(log.With(defaultLoggerWrapper.logger, "msg", msg)).Log(args...)
}

func (l *loggerWrapper) Error(err error, msg string, args ...any) {
	if err == nil {
		level.Error(log.With(l.logger, "msg", msg)).Log(args...)
		return
	}

	if msg == "" {
		level.Error(log.With(l.logger, "err", err)).Log(args...)
		return
	}

    level.Error(log.With(l.logger, "msg", msg, "err", err)).Log(args...)
}

func (l *loggerWrapper) Infof(format string, args ...any) {
    level.Info(l.logger).Log("msg", fmt.Sprintf(format, args...))
}

func (l *loggerWrapper) Warnf(format string, args ...any) {
    level.Warn(l.logger).Log("msg", fmt.Sprintf(format, args...))
}

func (l *loggerWrapper) Debugf(format string, args ...any) {
    level.Debug(l.logger).Log("msg", fmt.Sprintf(format, args...))
}

func (l *loggerWrapper) Errorf(err error, format string, args ...any) {
    level.Error(l.logger).Log("err", err.Error(), "msg", fmt.Sprintf(format, args...))
}

func (l *loggerWrapper) WithTraceId(ctx context.Context) *loggerWrapper {
    traceId := instruments.ExtractTraceID(ctx)
    logger := log.With(l.logger, "trace_id", traceId)

	return &loggerWrapper{
		logger: logger,
	}
}

func (l *loggerWrapper) WithFields(fields map[string]any) *loggerWrapper {
    fieldArgs := mapToSlice(fields)
    logger := log.With(l.logger, fieldArgs...)

	return &loggerWrapper{
		logger: logger,
	}
}

func Info(msg string, args ...any) {
    level.Info(log.With(defaultLoggerWrapper.logger, "msg", msg)).Log(args...)
}

func Warn(msg string, args ...any) {
    level.Warn(log.With(defaultLoggerWrapper.logger, "msg", msg)).Log(args...)
}

func Debug(msg string, args ...any) {
    level.Debug(log.With(defaultLoggerWrapper.logger, "msg", msg)).Log(args...)
}

func Error(err error, msg string, args ...any) {
	if err == nil {
		level.Error(log.With(defaultLoggerWrapper.logger, "msg", msg)).Log(args...)
		return
	}

	if msg == "" {
		level.Error(log.With(defaultLoggerWrapper.logger, "err", err)).Log(args...)
		return
	}

    level.Error(log.With(defaultLoggerWrapper.logger, "msg", msg, "err", err)).Log(args...)
}

func Infof(format string, args ...any) {
    level.Info(defaultLoggerWrapper.logger).Log("msg", fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...any) {
    level.Warn(defaultLoggerWrapper.logger).Log("msg", fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...any) {
    level.Debug(defaultLoggerWrapper.logger).Log("msg", fmt.Sprintf(format, args...))
}

func Errorf(err error, format string, args ...any) {
    level.Error(defaultLoggerWrapper.logger).Log("err", err.Error(), "msg", fmt.Sprintf(format, args...))
}

func WithFields(fields map[string]any) *loggerWrapper {
	return defaultLoggerWrapper.WithFields(fields)
}

func WithTraceId(ctx context.Context) *loggerWrapper {
	return defaultLoggerWrapper.WithTraceId(ctx)
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

func GetLogger() log.Logger {
	return defaultLoggerWrapper.logger
}
