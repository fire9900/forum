package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func TestLoggerInit() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func InitLogger(main bool) error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if main {
		config.OutputPaths = []string{"stdout", "./forum.log"}
		config.ErrorOutputPaths = []string{"stderr", "./forum-error.log"}
	} else {
		config.OutputPaths = []string{"stdout", "./test.log"}
		config.ErrorOutputPaths = []string{"stderr", "./test-error.log"}
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}
	return nil
}
