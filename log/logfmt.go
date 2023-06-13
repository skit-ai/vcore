package log

import (
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// InitLogfmtLogger initialises the global gokit logger and overrides the
// default logger for the server.
func InitLogfmtLogger() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	// Using legacy logger's level for level filtering
	logger = level.NewFilter(logger, LevelFilter(defaultLogger.level))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	KitLogger = log.With(logger, "caller", log.Caller(4))
}

// CheckFatal prints an error and exits with error code 1 if err is non-nil.
func CheckFatal(location string, err error) {
	if err == nil {
		return
	}

	logger := level.Error(KitLogger)
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
