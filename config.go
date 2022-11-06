package logger

import (
	"fmt"
	fluentd_http "github.com/m4schini/logger/sink/fluentd-http"
	"go.uber.org/zap"
	"net/url"
)

type FluentdConfig struct {
	Address string `json:"address" yaml:"address"`
	Name    string `json:"name" yaml:"name"`
}

type Config struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level zap.AtomicLevel `json:"level" yaml:"level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `json:"development" yaml:"development"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`
	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	Sampling *zap.SamplingConfig `json:"sampling" yaml:"sampling"`
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoding string `json:"encoding" yaml:"encoding"`
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`
	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`

	// ===> CUSTOM <===
	Fluentd FluentdConfig `json:"fluentd" yaml:"fluentd"`
}

func From(c zap.Config) Config {
	return Config{
		Level:             c.Level,
		Development:       c.Development,
		DisableCaller:     c.DisableCaller,
		DisableStacktrace: c.DisableStacktrace,
		Sampling:          c.Sampling,
		Encoding:          c.Encoding,
		OutputPaths:       c.OutputPaths,
		ErrorOutputPaths:  c.ErrorOutputPaths,
	}
}

func Update(c *zap.Config, config Config) *zap.Config {
	c.Level = config.Level
	c.Development = config.Development
	c.DisableCaller = config.DisableCaller
	c.DisableStacktrace = config.DisableStacktrace
	c.Sampling = config.Sampling
	c.Encoding = config.Encoding
	c.OutputPaths = config.OutputPaths
	c.ErrorOutputPaths = config.ErrorOutputPaths

	return c
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
