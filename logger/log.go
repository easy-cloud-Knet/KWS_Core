package syslogger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogFilePath() string {
	logDir := "/var/log/kws/"
	currentDate := time.Now().Format("20060102") // YYYYMMDD 포맷
	logFile := fmt.Sprintf("%s%s.log", logDir, currentDate)
	return logFile
}

func InitialLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	filePath := getLogFilePath()

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:       true,
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stdout", filePath},
		ErrorOutputPaths:  []string{"stderr", filePath},
	}

	return zap.Must(config.Build())

}
