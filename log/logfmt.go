package log

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/skit-ai/vcore/instruments"
)

type logfmtWrapper struct {
    logger *log.Logger
}

// InitLogfmtLogger initialises a logfmt logger from go-kit
func initLogfmt() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	// Using legacy logger's level for level filtering
	logger = level.NewFilter(logger, LevelFilter(defaultLogger.level))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.Caller(5))
	logfmtLogger = logfmtWrapper{&logger}
}

// CheckFatal prints an error and exits with error code 1 if err is non-nil.
func CheckFatal(location string, err error) {
	if err == nil {
		return
	}

	logger := level.Error(*logfmtLogger.logger)
	if location != "" {
		logger = log.With(logger, "msg", "error "+location)
	}
	// %+v gets the stack trace from errors using github.com/pkg/errors
	errStr := fmt.Sprintf("%+v", err)
	fmt.Fprintln(os.Stderr, errStr)

	logger.Log("err", errStr)
	os.Exit(1)
}

func LevelFilter(l int) level.Option {
	switch l {
	case DEBUG:
		return level.AllowDebug()
	case INFO:
		return level.AllowInfo()
	case WARN:
		return level.AllowWarn()
	case ERROR:
		return level.AllowError()
	default:
		return level.AllowAll()
	}
}

// TODO: better naming
func mapToArgs(m map[string]interface{}) []interface{} {
    var args []interface{}
    for k, v := range m {
        args = append(args, k)
        args = append(args, v)
    }
    return args
}

func (lw *logfmtWrapper) Info(args ...interface{}) {
    level.Info(*lw.logger).Log("msg", args[0])
}

func (lw *logfmtWrapper) Warn(args ...interface{}) {
    level.Warn(*lw.logger).Log("msg", args[0])
}

func (lw *logfmtWrapper) Debug(args ...interface{}) {
    level.Debug(*lw.logger).Log("msg", args[0])
}

func (lw *logfmtWrapper) Error(err error, args ...interface{}) {
    level.Error(*lw.logger).Log("msg", args)
}

func (lw *logfmtWrapper) Infof(format string, args ...interface{}) {
    level.Info(*lw.logger).Log("msg", fmt.Sprintf(format, args...))
}

func (lw *logfmtWrapper) Warnf(format string, args ...interface{}) {
    level.Warn(*lw.logger).Log("msg", fmt.Sprintf(format, args...))
}

func (lw *logfmtWrapper) Debugf(format string, args ...interface{}) {
    level.Debug(*lw.logger).Log("msg", fmt.Sprintf(format, args...))
}

func (lw *logfmtWrapper) Errorf(err error, format string, args ...interface{}) {
    level.Error(*lw.logger).Log("msg", fmt.Sprintf(format, args...))
}

func (lw *logfmtWrapper) WithTraceId(ctx context.Context) *logfmtWrapper {
    traceId := instruments.ExtractTraceID(ctx)
    *lw.logger = log.With(*lw.logger, "trace_id", traceId)
    return lw
}

func (lw *logfmtWrapper) WithFields(fields map[string]interface{}) *logfmtWrapper {
    fieldArgs := mapToArgs(fields)
    *lw.logger = log.With(*lw.logger, fieldArgs...)
    return lw
}

func WithTraceId(ctx context.Context) *logfmtWrapper {
    return logfmtLogger.WithTraceId(ctx)
}

func WithFields(fields map[string]interface{}) *logfmtWrapper {
    return logfmtLogger.WithFields(fields)
}

func WithContext(ctx context.Context, fields map[string]interface{}) *logfmtWrapper {
    if ctx == nil {
        return logfmtLogger.WithFields(fields)
    }

    if fields == nil {
        return logfmtLogger.WithTraceId(ctx)
    }

    return logfmtLogger.WithFields(fields).WithTraceId(ctx)
}
