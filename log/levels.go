package log

import (
	"fmt"
	"log"
	"strings"

	"github.com/skit-ai/vcore/errors"
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

// Prefix based on the log level to be added to every log statement
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

// Logs using stdlib logger based on the log level set
func (logger *Logger) log(LEVEL int, err error, format string, args ...interface{}) {
	if logger.isLevel(LEVEL) {
		if err == nil {
			log.Printf("%s %s\n", levelPrefix(LEVEL), fmt.Sprintf(format, args...))
		} else {
			// Do not use log.Fatalf since it will call os.Exit and terminate the program
			log.Printf("%s %s:\n%s\n", levelPrefix(LEVEL), fmt.Sprintf(format, args...), errors.Stacktrace(err))
		}

	}
}

// Checks if the logger has the ability to log at a given log level
func (logger *Logger) isLevel(LEVEL int) bool {
	return logger.level >= LEVEL
}

// Set the level of the logger
func (logger *Logger) SetLevel(level int) {
	if level <= TRACE && level >= ERROR {
		defaultLogger.level = level
	} else {
		_format := "Cannot set log level to %d. Log levels allowed are %s. Default log level is %d(WARN)"
		logger.Warnf(_format, level, joinInt(",", []int{TRACE, DEBUG, INFO, WARN, ERROR}), WARN)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Methods to check the log level
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Methods to check the log level of a particular Logger struct instance

func (logger *Logger) IsTrace() bool {
	return logger.isLevel(TRACE)
}

func (logger *Logger) IsDebug() bool {
	return logger.isLevel(DEBUG)
}

func (logger *Logger) IsInfo() bool {
	return logger.isLevel(INFO)
}

func (logger *Logger) IsWarn() bool {
	return logger.isLevel(WARN)
}

func (logger *Logger) IsError() bool {
	return logger.isLevel(ERROR)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Methods to check default logger's log level

func IsTrace() bool {
	return defaultLogger.IsTrace()
}

func IsDebug() bool {
	return defaultLogger.IsDebug()
}

func IsInfo() bool {
	return defaultLogger.IsInfo()
}

func IsWarn() bool {
	return defaultLogger.IsWarn()
}

func IsError() bool {
	return defaultLogger.IsError()
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Methods to log a message using a Logger struct instance
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//Methods to log a message using formats

func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.log(TRACE, nil, format, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.log(DEBUG, nil, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.log(INFO, nil, format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.log(WARN, nil, format, args...)
}

func (logger *Logger) Errorf(err error, format string, args ...interface{}) {
	logger.log(ERROR, err, format, args...)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//Methods to log a message without using formats

func (logger *Logger) Trace(args ...interface{}) {
	logger.log(TRACE, nil, repeat("%v", len(args)), args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.log(DEBUG, nil, repeat("%v", len(args)), args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.log(INFO, nil, repeat("%v", len(args)), args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.log(WARN, nil, repeat("%v", len(args)), args...)
}

func (logger *Logger) Error(err error, args ...interface{}) {
	logger.log(ERROR, err, repeat("%v", len(args)), args...)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Methods to log a message using the default logger
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Methods to log a message using the default logger without a format

func Trace(args ...interface{}) {
	defaultLogger.Trace(args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Error(err error, args ...interface{}) {
	defaultLogger.Error(err, args...)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Methods to log messages using the default logger with a format

func Tracef(format string, args ...interface{}) {
	defaultLogger.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Errorf(err error, format string, args ...interface{}) {
	defaultLogger.Errorf(err, format, args...)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Set the level of the default logger
func SetLevel(level int) {
	defaultLogger.SetLevel(level)
}

// Wrapper for log.Fatal
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

// Wrapper for log.Fatalf
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// Repeat a string x times while building a string
func repeat(s string, x int) string {
	var builder strings.Builder
	for i := 0; i < x; i++ {
		if i == (x - 1) {
			builder.WriteString(s)
		} else {
			builder.WriteString(s)
			builder.WriteString(" ")
		}
	}
	return builder.String()
}

// Concatenates a variable slice of strings
func joinInt(delimiter string, slice []int) string {
	var builder strings.Builder
	for i, item := range slice {
		builder.WriteString(fmt.Sprintf("%v", item))
		if i != len(slice)-1 {
			builder.WriteString(delimiter)
		}
	}
	return builder.String()
}
