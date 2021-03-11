package logging

import (
	"20dojo-online/pkg/dcontext"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var accessLogger *zap.SugaredLogger
var ApplicationLogger *zap.Logger

func AccessLogging(request *http.Request) {
	accessLogger.Infow("incoming request",
		zap.String("host", request.Host),
		zap.String("remoteAddress", request.RemoteAddr),
		zap.String("method", request.Method),
		zap.String("path", request.URL.Path),
		zap.String("requestID", dcontext.GetRequestIDFromContext(request.Context())))
}

func NewAccessLogger(zapCoreLevel zapcore.Level) {
	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapCoreLevel),
		OutputPaths:      []string{"pkg/logging/log/access.log"},
		ErrorOutputPaths: []string{"pkg/logging/log/access.log"},
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	accessLogger = logger.Sugar()
}

func NewApplicationLogging(zapCoreLevel zapcore.Level) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		MessageKey:   "msg",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	file, err := os.Create("pkg/logging/log/application.log")
	if err != nil {
		log.Fatal(err)
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	logCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		zapCoreLevel,
	)

	return zap.New(zapcore.NewTee(
		consoleCore,
		logCore,
	))
}

func init() {
	env := os.Getenv("ENV")
	var zapCoreLevel zapcore.Level
	if env == "production" {
		zapCoreLevel = zap.InfoLevel
	} else {
		zapCoreLevel = zap.DebugLevel
	}

	NewAccessLogger(zapCoreLevel)
	ApplicationLogger = NewApplicationLogging(zapCoreLevel)
	ApplicationLogger = ApplicationLogger.WithOptions(zap.AddCaller())

	defer accessLogger.Sync()
	defer ApplicationLogger.Sync()
}
