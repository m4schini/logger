package fluentd_http

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type FluentdHttpSink struct {
	FluentdAddress url.URL
	zap.Sink
}

func (f *FluentdHttpSink) Write(p []byte) (n int, err error) {
	r, err := http.Post(f.FluentdAddress.String(),
		"application/json",
		strings.NewReader(string(p)))
	if err != nil {
		return 0, err
	}
	if r.StatusCode == http.StatusOK {
		return len(p), nil
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("%v", string(b))
}

func (f *FluentdHttpSink) Sync() error {
	return nil
}

func (f *FluentdHttpSink) Close() error {
	return nil
}
