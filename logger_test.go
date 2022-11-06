package logger

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	log := Named("test").Sugar()

	for {
		log.Info("test")
		time.Sleep(5 * time.Second)
	}

}
