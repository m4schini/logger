package logger

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	err := Init(Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		OutputPaths:       []string{"fluentd://"},
		ErrorOutputPaths:  []string{"fluentd://", "stderr"},
		Fluentd: FluentdConfig{
			Address: "http://localhost:9880",
			Name:    "TestLogger",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	log := Named("test").Sugar()

	for {
		log.Info("test")
		time.Sleep(5 * time.Second)
	}

}
