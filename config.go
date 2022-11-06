package logger

import (
	"fmt"
	fluentd_http "github.com/m4schini/logger/sink/fluentd-http"
	"go.uber.org/zap"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	envPrefix = "LOGGING_"
)

var (
	Level           zap.AtomicLevel
	DevelopmentMode bool
	Encoding        string
	OutputPaths     []string = make([]string, 0)
	ErrOutputPaths  []string = make([]string, 0)
)

// init level
func init() {
	var err error
	Level, err = zap.ParseAtomicLevel(getEnvOrDefault(envPrefix+"LEVEL", "debug"))
	if err != nil {
		log.Fatalln(err)
	}
}

// init dev mode
func init() {
	DevelopmentMode = getEnvOrDefault(envPrefix+"DEV", "false") == "true"
}

// init encoding
func init() {
	Encoding = getEnvOrDefault(envPrefix+"ENCODING", "json")
}

// init fluentd
func init() {
	addr := os.Getenv(envPrefix + "FD_ADDR")
	name := os.Getenv(envPrefix + "FD_NAME")
	if addr != "" && name != "" {
		err := RegisterFluentdSink(addr, name)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// init output
func init() {
	outputPaths := getEnvOrDefault(envPrefix+"OUTPUT", "stderr")
	for _, s := range strings.Split(outputPaths, ",") {
		OutputPaths = append(OutputPaths, strings.TrimSpace(s))
	}
}

// init error output
func init() {
	outputPaths := getEnvOrDefault(envPrefix+"ERR_OUTPUT", "stderr")
	for _, s := range strings.Split(outputPaths, ",") {
		ErrOutputPaths = append(OutputPaths, strings.TrimSpace(s))
	}
}

func getEnvOrDefault(key, def string) string {
	e := os.Getenv(key)
	if e == "" {
		return def
	} else {
		return e
	}
}

func RegisterFluentdSink(address, name string) error {
	return zap.RegisterSink("fluentd", func(url *url.URL) (zap.Sink, error) {
		u, err := url.Parse(address)
		if err != nil {
			return nil, err
		}
		if !(u.Scheme == "http" || u.Scheme == "https") {
			return nil, fmt.Errorf("url scheme has to be http or https")
		}

		u.Path = name
		return &fluentd_http.FluentdHttpSink{
			FluentdAddress: *u,
		}, nil
	})
}
