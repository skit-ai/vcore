package log

import (
	"fmt"
	"log"
	"vcore/errors"
	"vcore/utils"
)

const (
	ERROR = iota
	WARN
	INFO
	DEBUG
	TRACE
)

type Logger struct {
	level int
}

var defaultLogger = Logger{WARN}

func levelPrefix(level int) string {
	switch level {
	case ERROR:
		return "[ERROR]"
	case WARN:
		return "[WARN]"
	case INFO:
		return "[INFO]"
	case DEBUG:
		return "[DEBUG]"
	case TRACE:
		return "[TRACE]"
	}
	return ""
}

func (logger *Logger) SetLevel(level int) {
	if level >= TRACE && level <= ERROR {
		defaultLogger.level = level
	} else {
		_format := "Cannot set log level to %d. Log levels allowed are %s. Default log level is %d(WARN)"
		logger.Warn(_format, level, utils.JoinInt(",", []int{TRACE, DEBUG, INFO, WARN, ERROR}), WARN)
	}
}

// Logs using stdlib logger based on the log level set
func (logger *Logger) log(LEVEL int, err error, format string, args ...interface{}) {
	if logger.level >= LEVEL {
		if err == nil {
			log.Printf("%s %s\n", levelPrefix(LEVEL), fmt.Sprintf(format, args...))
		} else {
			// Do not use log.Fatalf since it will call os.Exit and terminate the program
			log.Printf("%s %s:\n%s\n", levelPrefix(LEVEL), fmt.Sprintf(format, args...), errors.Stacktrace(err))
		}

	}
}

func (logger *Logger) Trace(format string, args ...interface{}) {
	logger.log(TRACE, nil, format, args...)
}

func (logger *Logger) Debug(format string, args ...interface{}) {
	logger.log(DEBUG, nil, format, args...)
}

func (logger *Logger) Info(format string, args ...interface{}) {
	logger.log(INFO, nil, format, args...)
}

func (logger *Logger) Warn(format string, args ...interface{}) {
	logger.log(WARN, nil, format, args...)
}

func (logger *Logger) Error(err error, format string, args ...interface{}) {
	logger.log(ERROR, err, format, args...)
}

func Trace(format string, args ...interface{}) {
	defaultLogger.Trace(format, args...)
}

func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

func Error(err error, format string, args ...interface{}) {
	defaultLogger.Error(err, format, args...)
}

func SetLevel(level int) {
	defaultLogger.SetLevel(level)
}
