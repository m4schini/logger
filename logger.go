package logger

import (
	"go.uber.org/zap"
	"log"
)

var logger *zap.Logger = nil

func init() {
	var err error

	config := zap.NewProductionConfig()
	config.Level = Level
	config.Development = DevelopmentMode
	config.Encoding = Encoding
	config.OutputPaths = OutputPaths
	config.ErrorOutputPaths = ErrOutputPaths

	logger, err = config.Build()
	if err != nil {
		log.Fatalln(err)
	}
}

func Named(names ...string) *zap.Logger {
	subLogger := logger.Sugar()
	for _, name := range names {
		subLogger.Named(name)
	}

	return subLogger.Desugar()
}

func Sync() error {
	return logger.Sync()
}
