package syslogger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)


func InitialLogger() *zap.SugaredLogger{

	baseLogger, err := zap.NewDevelopment()
	
	if err!=nil{ 
		panic("Error initializing logger")
	}
	logger:= baseLogger.Sugar()
	return logger
}



func LoggerMiddleware(next http.Handler, logger *zap.SugaredLogger) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start:= time.Now()

		next.ServeHTTP(w,r)

		elapsed := time.Since(start)
		logger.Infof("Handled http request time elapsed %d", elapsed)
	})
}