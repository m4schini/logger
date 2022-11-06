package logger

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
	"sync"
)

var initWG sync.WaitGroup
var logger *zap.Logger = nil

func init() {
	initWG.Add(1)
	var err error

	lvl := os.Getenv("LOGGING_LEVEL")
	switch strings.ToLower(lvl) {
	case "prod":
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatalln(err)
		}
		initWG.Done()
		break
	case "debug":
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger, err = config.Build()
		if err != nil {
			log.Fatalln(err)
		}
		initWG.Done()
		break
	case "dev":
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalln(err)
		}
		initWG.Done()
		break
	case "custom":
		content, err := os.ReadFile("logging.yml") // the file is inside the local directory
		if err != nil {
			log.Fatalln(err)
		}

		var cfg Config
		err = yaml.Unmarshal(content, &cfg)
		if err != nil {
			log.Fatalln(err)
		}

		err = Init(cfg)
		if err != nil {
			log.Fatalln(err)
		}
		initWG.Done()
		break
	case "disabled":
		logger = zap.NewNop()
		break
	default:
		break
	}
}

func Named(names ...string) *zap.Logger {
	initWG.Wait()

	subLogger := logger.Sugar()
	for _, name := range names {
		subLogger.Named(name)
	}

	return subLogger.Desugar()
}

func Sync() error {
	return logger.Sync()
}

func Init(config Config) error {
	if config.Fluentd.Address != "" {
		err := RegisterFluentdSink(config.Fluentd.Address, config.Fluentd.Name)
		if err != nil {
			return err
		}
	}

	prodCfg := zap.NewProductionConfig()
	Update(&prodCfg, config)
	newLogger, err := prodCfg.Build()
	if err != nil {
		return err
	}

	logger = newLogger
	initWG.Done()

	logger.Sugar().Infow("initialized logging with custom config", "cfg", config)
	return nil
}
