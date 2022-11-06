package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"testing"
)

func TestServeMetrics(t *testing.T) {
	c := promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	c.Inc()
	c.Inc()
	c.Inc()
	c.Inc()

	err := serveMetrics("2112")
	if err != nil {
		t.Fatal(err)
	}
}
