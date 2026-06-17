package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/denislemire/tesla-wall-connector-exporter/internal/client"
	"github.com/denislemire/tesla-wall-connector-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		twcAddress    = flag.String("twc.address", envOr("TWC_ADDRESS", ""), "Tesla Wall Connector hostname or IP (required)")
		listenAddress = flag.String("web.listen-address", envOr("WEB_LISTEN_ADDRESS", ":9859"), "Address to listen on for HTTP requests")
		metricsPath   = flag.String("web.metrics-path", envOr("WEB_METRICS_PATH", "/metrics"), "Path to expose metrics on")
		scrapeTimeout = flag.Duration("twc.timeout", durationEnv("TWC_TIMEOUT", 5*time.Second), "HTTP timeout when polling the wall connector")
	)
	flag.Parse()

	if *twcAddress == "" {
		slog.Error("twc.address is required (flag or TWC_ADDRESS env)")
		os.Exit(2)
	}

	slog.Info("starting exporter", "twc_address", *twcAddress, "listen", *listenAddress)

	twcClient := client.New(*twcAddress, *scrapeTimeout)
	reg := prometheus.NewRegistry()
	reg.MustRegister(collector.New(twcClient))

	mux := http.NewServeMux()
	mux.Handle(*metricsPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<html><body><h1>Tesla Wall Connector Exporter</h1><p><a href="` + *metricsPath + `">Metrics</a></p></body></html>`))
	})

	if err := http.ListenAndServe(*listenAddress, mux); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return fallback
}
