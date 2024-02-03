package log

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	Log_Path := os.Getenv("LOG_PATH")
	logger, _ := zap.NewProduction()
	// fmt.Println(Log_Path)
	config := fmt.Sprintf(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "%s"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`, Log_Path)
	rawJSON := []byte(config)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		log.Fatalf("Error init log: %s", err)
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Error build log: %s", err)
	}
	// 官方文档 没有检测这个方法返回的错误
	defer logger.Sync()
	Logger = logger
}
