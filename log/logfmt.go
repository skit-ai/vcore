package log

import (
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// InitLogger initialises the global gokit logger and overrides the
// default logger for the server.
func InitLogger() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	// add support for level based logging
	// logger = level.NewFilter(logger, LevelFilter(cfg.LogLevel.String()))
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
