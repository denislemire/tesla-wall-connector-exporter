package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	pathVitals     = "/api/1/vitals"
	pathLifetime   = "/api/1/lifetime"
	pathWiFiStatus = "/api/1/wifi_status"
	pathVersion    = "/api/1/version"
)

// Client polls the Tesla Wall Connector Gen3 local HTTP API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New creates a client for the given address (hostname or IP, no scheme).
func New(address string, timeout time.Duration) *Client {
	addr := strings.TrimSpace(address)
	addr = strings.TrimPrefix(addr, "http://")
	addr = strings.TrimPrefix(addr, "https://")
	addr = strings.TrimSuffix(addr, "/")

	return &Client{
		baseURL: "http://" + addr,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) get(path string) ([]byte, float64, error) {
	start := time.Now()
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return nil, time.Since(start).Seconds(), err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	elapsed := time.Since(start).Seconds()
	if err != nil {
		return nil, elapsed, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, elapsed, fmt.Errorf("GET %s: HTTP %d", path, resp.StatusCode)
	}
	return body, elapsed, nil
}

// SanitizeJSON replaces invalid JSON tokens returned by some firmware (e.g. :nan).
func SanitizeJSON(body []byte) []byte {
	// Lifetime may contain "avg_startup_temp":nan
	out := bytes.ReplaceAll(body, []byte(":nan"), []byte(":null"))
	return out
}

func (c *Client) FetchVitals() (Vitals, float64, error) {
	body, elapsed, err := c.get(pathVitals)
	if err != nil {
		return Vitals{}, elapsed, err
	}
	var v Vitals
	if err := json.Unmarshal(body, &v); err != nil {
		return Vitals{}, elapsed, fmt.Errorf("decode vitals: %w", err)
	}
	return v, elapsed, nil
}

func (c *Client) FetchLifetime() (Lifetime, float64, error) {
	body, elapsed, err := c.get(pathLifetime)
	if err != nil {
		return Lifetime{}, elapsed, err
	}
	body = SanitizeJSON(body)
	var lt Lifetime
	if err := json.Unmarshal(body, &lt); err != nil {
		return Lifetime{}, elapsed, fmt.Errorf("decode lifetime: %w", err)
	}
	return lt, elapsed, nil
}

func (c *Client) FetchWiFiStatus() (WiFiStatus, float64, error) {
	body, elapsed, err := c.get(pathWiFiStatus)
	if err != nil {
		return WiFiStatus{}, elapsed, err
	}
	var w WiFiStatus
	if err := json.Unmarshal(body, &w); err != nil {
		return WiFiStatus{}, elapsed, fmt.Errorf("decode wifi_status: %w", err)
	}
	return w, elapsed, nil
}

func (c *Client) FetchVersion() (Version, float64, error) {
	body, elapsed, err := c.get(pathVersion)
	if err != nil {
		return Version{}, elapsed, err
	}
	var v Version
	if err := json.Unmarshal(body, &v); err != nil {
		return Version{}, elapsed, fmt.Errorf("decode version: %w", err)
	}
	return v, elapsed, nil
}
