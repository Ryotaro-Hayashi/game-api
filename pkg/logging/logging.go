package logging

import (
	"20dojo-online/pkg/dcontext"
	"log"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var accessLogger *zap.SugaredLogger

func AccessLogging(request *http.Request) {
	accessLogger.Infow("incoming request",
		zap.String("host", request.Host),
		zap.String("remoteAddress", request.RemoteAddr),
		zap.String("method", request.Method),
		zap.String("path", request.URL.Path),
		zap.String("requestID", dcontext.GetRequestIDFromContext(request.Context())))
}

func NewAccessLogger() {
	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"pkg/logging/log/access.log"},
		ErrorOutputPaths: []string{"pkg/logging/log/access.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			TimeKey:       "time",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			CallerKey:     "caller",
			EncodeCaller:  zapcore.ShortCallerEncoder,
			StacktraceKey: "stackTrace",
		},
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
