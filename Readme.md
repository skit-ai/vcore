# golang/vcore

## log

The log package is a basic wrapper on the standard log package  in Go's stdlib.
The log package logs to STDOUT and supports log levels(`int`).

The log levels currently supported are:

* 0 - ERROR
* 1 - WARN
* 2 - INFO
* 3 - DEBUG
* 4 - TRACE

To log using this package, one needs to use an instance of the `log.Logger` struct.

The `log.Logger` struct supports the following methods which can be used to log messages:
* `log.Trace(args ...interface{})`
* `log.Tracef(format string, args ...interface{})`
* `log.Debug(args ...interface{})`
* `log.Debugf(format string, args ...interface{})`
* `log.Info(args ...interface{})`
* `log.Infof(format string, args ...interface{})`
* `log.Warn(args ...interface{})`
* `log.Warnf(format string, args ...interface{})`
* `log.Error(err error, args ...interface{})`
* `log.Errorf(err error, format string, args ...interface{})`

Each of these methods are wrappers that correspond to a log level. This enforces the user to take cognizance of the log 
level of whatever they are attempting to log.

To set the log level on a `log.Logger` struct, make use of the `log.SetLevel(level int)` function.

#### Default Logger
To quickly start logging messages, make use of the default logger(default level `WARN`). This can be done by simply 
calling the functions stated above.

Eg. To add a trace log
```go
headers := make(map[string]string)
for k := range req.Header {
			headers[strings.ToLower(k)] = req.Header.Get(k)
}
log.SetLevel(log.DEBUG)
log.Tracef("Headers: %s", headers)
``` 
 
 Here, we directly make use of the default logger. Please note, since the log level is set to DEBUG here, this trace 
 message will not be logged.
 
 #### Custom Logger
 
```go
customLogger := log.Logger{log.DEBUG}
customLogger.Debug("This is a debug message")
```
