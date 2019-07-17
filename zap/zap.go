// Wrapper for uber-go/zap logger
package zap

import (
	"go.uber.org/zap"
	// "log"
	"encoding/json"
	"os"
)

var Logger *ZapLogger

type ZapLogger struct {
	logger *zap.Logger
}

// func New() *ZapLogger{
//   logger, _ := zap.NewProduction()
//   return &ZapLogger{"logger": logger}
// }

// func init(){
//   logger, _ := zap.NewProduction()
//   Logger = &ZapLogger{logger: logger}
// }

func init() {
	var encoding_type string
	var cfg zap.Config

	app_env := os.Getenv("APP_ENV")
	//source := os.Getenv("APP_NAME")
	if app_env == "production" {
		encoding_type = "json"
	} else {
		encoding_type = "console"
	}
	rawJSON := []byte(`{
      "level": "debug",
      "encoding": "` + encoding_type + `",
      "outputPaths": ["stdout"],
      "errorOutputPaths": ["stderr"],
      "encoderConfig": {
        "messageKey": "log",
        "levelKey": "level",
        "levelEncoder": "uppercase"
      }
    }`)

	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	//defer Logger.Sync()
	Logger = &ZapLogger{logger: logger}
}

// Proxy for reflect method
// func Reflect(key string, val interface{}) zap.Field{
//   return zap.Reflect(key, val)
// }

func (zaplogger *ZapLogger) Info(log string, args ...map[string]interface{}) {
	// Build the log metadata
	var log_message map[string]interface{} = map[string]interface{}{}
	if len(args) > 1 {
		log_message = args[1]
	}
	if len(args) > 0 {
		log_message["payload"] = args[0]
	}
	log_message["level"] = "INFO"
	log_message["event"] = "log"
	// tid, event, timestamp
	zaplogger.logger.Info(log, zap.Reflect("message", log_message))
}

// Debug, Warn, Error, Fatal, Panic
