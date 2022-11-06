package logger

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

var logger *zap.Logger = nil
var Cfg *Config

func init() {
	var err error

	if Cfg != nil {
		err := applyConfig(*Cfg)
		if err != nil {
			log.Fatalln(err)
		}
	}

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

		err = applyConfig(cfg)
		if err != nil {
			log.Fatalln(err)
		}
		break
	case "disabled":
		logger = zap.NewNop()
		break
	default:
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

func applyConfig(config Config) error {
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

	*logger = *newLogger

	logger.Sugar().Infow("initialized logging with custom config", "cfg", config)
	return nil
}
