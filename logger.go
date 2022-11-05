package logger

import (
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

var logger *zap.Logger

func init() {
	var err error

	lvl := os.Getenv("LOGGING_LEVEL")
	switch strings.ToLower(lvl) {
	case "prod":
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatalln(err)
		}
		break
	case "debug":
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger, err = config.Build()
		if err != nil {
			log.Fatalln(err)
		}
		break
	case "dev":
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalln(err)
		}
		break
	case "debug-sink-redis":
		//err := zap.RegisterSink("redis", func(url *url.URL) (zap.Sink, error) {
		//	s := NewRedisSink("redis", "")
		//	s.Key = "log:web-scraper-module"
		//	return s, nil
		//})
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//config := zap.NewProductionConfig()
		//config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		//config.OutputPaths = []string{"redis"}
		//config.ErrorOutputPaths = []string{"redis", "stdout"}
		//logger, err = config.Build()
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//break
	default:
		logger = zap.NewNop()
		break
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

func UseLogger(newLogger *zap.Logger) {
	logger = newLogger
}
