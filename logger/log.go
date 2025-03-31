package syslogger

import (
	"fmt"
	"net/http"
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

func InitialLogger() *zap.Logger{
	encoderCfg := zap.NewProductionEncoderConfig()
    encoderCfg.TimeKey = "timestamp"
    encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	filePath:= getLogFilePath()


	config:= zap.Config{
		Level:  zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development: true,	
		DisableCaller: true,
		DisableStacktrace: true,
		Encoding: "json",
		EncoderConfig:    encoderCfg,
		OutputPaths: []string{"stdout",filePath},
		ErrorOutputPaths: []string{"stderr",filePath},
	}

	return zap.Must(config.Build())


}



func LoggerMiddleware(next http.Handler, logger *zap.Logger) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start:= time.Now()

		next.ServeHTTP(w,r)

		elapsed := time.Since(start)
		logger.Info("http response sent",zap.Duration("time elapsed", elapsed))
	})
}




