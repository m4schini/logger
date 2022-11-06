package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

var (
	metricsPort = os.Getenv("METRICS_PORT")
)

func ServeMetrics() error {
	if metricsPort == "" {
		log.Fatalln("no METRICS_PORT defined")
	}

	return serveMetrics(metricsPort)
}

func serveMetrics(port string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":"+port, nil)
}
