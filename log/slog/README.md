# Slog

Slog (short for "structured log") is a wrapper over [go-kit/log](https://github.com/go-kit/log) with a defined log structure and a set of helper methods.


## Configurations

| Environment Variable | Default Value   | Allowed Values   |
| --- | :---: | :---: |
| LOG_LEVEL | "info"   | "debug", "info", "warn", "error"   |
| LOG_CALLER_DEPTH | 4   | Z |


## Log Levels & Filtering

All log levels are allowed to be logged by default. Use the config "LOG_LEVEL" to filter out.

| Set Value | |
| ---   | --- |
| "debug" | debug + info + warn + error |
| "info" | info + warn + error |
| "warn" | warn + error |
| "error" | error |


## Usage

1. Log line without custom fields:
```
slog.Debug("the five boxing")
```
```
level=debug ts=2011-08-22T10:07:34.329181992Z caller=main.go:21 msg="the five boxing"
```

2. Log line with custom fields:
```
slog.Warn("wizards jump quickly", "unit_bytes", 8008)
```
```
level=warn ts=2012-07-18T11:27:32.616223846Z caller=main.go:70 msg="the five boxing" unit_bytes=8008
```

3. Create a logger with trace_id attached
```
slogger := slog.WithTraceId(ctx)
slogger.Info("the quick brown")
```

4. Create a logger with custom fields attached
```
logFields := map[string]interface{}{
	"integration_uuid": response.IntegrationUuid,
}
slogger := slog.WithFields(logFields)
slogger.Warn("fox jumps over")
```

5. Create a logger with trace_id and custom fields attached
```
logFields := map[string]interface{}{
	"integration_uuid": response.IntegrationUuid,
}
slogger := slog.WithFields(logFields).WithTraceId(ctx)
slogger.Debug("the lazy dog")
```


## Logging Methods
```
slog.Debug(msg string, args ...any)
slog.Debugf(format string, args ...any)
slog.Info(msg string, args ...any)
slog.Infof(format string, args ...any)
slog.Warn(msg string, args ...any)
slog.Warnf(format string, args ...any)
slog.Error(err error, msg string, args ...any)
slog.Errorf(err error, format string, args ...any)
```


## Helper Methods

```
slog.WithFields(fields map[string]any) Logger
slog.WithTraceId(ctx context.Context) Logger
slog.DefaultLogger() log.Logger
```
