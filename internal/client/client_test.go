package client_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/denislemire/tesla-wall-connector-exporter/internal/client"
)

func TestSanitizeJSON(t *testing.T) {
	in := []byte(`{"avg_startup_temp":nan,"energy_wh":1}`)
	out := client.SanitizeJSON(in)
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal sanitized json: %v", err)
	}
	if m["avg_startup_temp"] != nil {
		t.Fatalf("expected null avg_startup_temp, got %v", m["avg_startup_temp"])
	}
}

func TestClientFetchAll(t *testing.T) {
	vitalsBody := mustRead(t, "vitals.json")
	lifetimeBody := mustRead(t, "lifetime_nan.json")
	wifiBody := mustRead(t, "wifi_status.json")
	versionBody := mustRead(t, "version.json")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/1/vitals":
			_, _ = w.Write(vitalsBody)
		case "/api/1/lifetime":
			_, _ = w.Write(lifetimeBody)
		case "/api/1/wifi_status":
			_, _ = w.Write(wifiBody)
		case "/api/1/version":
			_, _ = w.Write(versionBody)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	host := srv.Listener.Addr().String()
	c := client.New(host, 2*time.Second)

	v, _, err := c.FetchVitals()
	if err != nil {
		t.Fatalf("FetchVitals: %v", err)
	}
	if !v.VehicleConnected || v.GridV != 240.0 {
		t.Fatalf("unexpected vitals: %+v", v)
	}
	if len(v.CurrentAlerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(v.CurrentAlerts))
	}

	lt, _, err := c.FetchLifetime()
	if err != nil {
		t.Fatalf("FetchLifetime: %v", err)
	}
	if lt.EnergyWh != 125000 || lt.ConnectorCycles != 16 {
		t.Fatalf("unexpected lifetime: %+v", lt)
	}
	if lt.AvgStartupTemp != nil {
		t.Fatalf("expected nil avg_startup_temp after nan sanitize")
	}

	wifi, _, err := c.FetchWiFiStatus()
	if err != nil {
		t.Fatalf("FetchWiFiStatus: %v", err)
	}
	if wifi.WiFiSSID != "home-wifi" || !wifi.Internet {
		t.Fatalf("unexpected wifi: %+v", wifi)
	}

	ver, _, err := c.FetchVersion()
	if err != nil {
		t.Fatalf("FetchVersion: %v", err)
	}
	if ver.SerialNumber != "PGT123456789" {
		t.Fatalf("unexpected version: %+v", ver)
	}
}

func mustRead(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return b
}
