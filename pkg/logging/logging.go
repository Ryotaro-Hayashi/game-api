package logging

import (
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var accessLogger *zap.SugaredLogger

func AccessLogging(request *http.Request) {
	accessLogger.Infow("incoming request",
		zap.String("Host", request.Host),
		zap.String("RemoteAddr", request.RemoteAddr),
		zap.String("Request Method", request.Method),
		zap.String("Path", request.URL.Path),
		zap.Any("RequestID", request.Context().Value("RequestID")),
		zap.Duration("elapsed", time.Second))
}

func NewAccessLogger() {
	config := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths: []string{"pkg/logging/log/access.log"},
	}
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	accessLogger = logger.Sugar()
}

func init() {
	NewAccessLogger()
}
