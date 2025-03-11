package syslogger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


func InitialLogger() *zap.Logger{
	encoderCfg := zap.NewProductionEncoderConfig()
    encoderCfg.TimeKey = "timestamp"
    encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder


	config:= zap.Config{
		Level:  zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development: true,	
		DisableCaller: true,
		DisableStacktrace: true,
		Encoding: "json",
		EncoderConfig:    encoderCfg,
		OutputPaths: []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return zap.Must(config.Build())

	// zap.Build(config)

	// baseLogger, err := zap.NewDevelopment()
	
	// if err!=nil{ 
	// 	panic("Error initializing logger")
	// }
	// logger:= baseLogger.Sugar()
	// return logger
}



func LoggerMiddleware(next http.Handler, logger *zap.Logger) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start:= time.Now()

		next.ServeHTTP(w,r)

		elapsed := time.Since(start)
		logger.Info("http response sent",zap.Duration("time elapsed", elapsed))
	})
}