package collector_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/denislemire/tesla-wall-connector-exporter/internal/client"
	"github.com/denislemire/tesla-wall-connector-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestCollectorExportsExpectedMetrics(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/api/1/")
		var file string
		switch name {
		case "vitals":
			file = "vitals.json"
		case "lifetime":
			file = "lifetime_nan.json"
		case "wifi_status":
			file = "wifi_status.json"
		case "version":
			file = "version.json"
		default:
			http.NotFound(w, r)
			return
		}
		b, err := os.ReadFile(filepath.Join("..", "..", "testdata", file))
		if err != nil {
			t.Fatalf("read fixture: %v", err)
		}
		_, _ = w.Write(b)
	}))
	defer srv.Close()

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector.New(client.New(srv.Listener.Addr().String(), 2*time.Second)))

	metricFamilies, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather: %v", err)
	}

	names := map[string]*dto.MetricFamily{}
	for _, mf := range metricFamilies {
		names[mf.GetName()] = mf
	}

	required := []string{
		"tesla_wall_connector_scrape_success",
		"tesla_wall_connector_grid_voltage_volts",
		"tesla_wall_connector_lifetime_energy_watt_hours_total",
		"tesla_wall_connector_wifi_connected",
		"tesla_wall_connector_info",
		"tesla_wall_connector_alert_info",
		"tesla_wall_connector_prox_volts",
		"tesla_wall_connector_lifetime_connector_cycles_total",
	}
	for _, r := range required {
		if _, ok := names[r]; !ok {
			t.Errorf("missing metric %s", r)
		}
	}

	if mf := names["tesla_wall_connector_scrape_success"]; mf != nil {
		for _, m := range mf.Metric {
			for _, lp := range m.Label {
				if lp.GetName() == "endpoint" && lp.GetValue() == "vitals" {
					if m.GetGauge().GetValue() != 1 {
						t.Errorf("vitals scrape_success want 1 got %v", m.GetGauge().GetValue())
					}
				}
			}
		}
	}
}
