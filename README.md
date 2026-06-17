# Tesla Wall Connector Gen3 Prometheus Exporter

Prometheus exporter for the **Tesla Wall Connector Gen3** (HWPC) local HTTP API. Polls all documented community endpoints and exposes metrics for [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack) or any Prometheus scraper.

> **Undocumented API:** Tesla does not publish this LAN API. It may change or be restricted in future firmware. Do not rely on this exporter for safety-critical control.

## Quick start

### Docker

```bash
docker run --rm -p 9859:9859 \
  ghcr.io/denislemire/tesla-wall-connector-exporter:v0.1.1 \
  --twc.address=192.168.1.100
```

Use a **hostname** (mDNS or DHCP reservation) when possible so IP changes do not break scraping.

### Binary

```bash
go run ./cmd/tesla-wall-connector-exporter --twc.address=192.168.1.100
curl -s localhost:9859/metrics | head
```

### Helm

```bash
helm install twc ./helm/tesla-wall-connector-exporter \
  --namespace monitoring \
  --set twc.address=192.168.1.100 \
  --set serviceMonitor.enabled=true \
  --set serviceMonitor.additionalLabels.release=monitoring
```

If the charger is only reachable from the host network (typical home LAN), set `hostNetwork: true` and pin to a node that can route to the charger:

```bash
helm upgrade --install twc ./helm/tesla-wall-connector-exporter \
  --namespace monitoring \
  --set twc.address=192.168.1.100 \
  --set hostNetwork=true \
  --set nodeSelector.kubernetes\\.io/hostname=worker-with-lan
```

## Configuration

| Flag / env | Default | Description |
|------------|---------|-------------|
| `--twc.address` / `TWC_ADDRESS` | *(required)* | Wall Connector hostname or IP |
| `--web.listen-address` / `WEB_LISTEN_ADDRESS` | `:9859` | Metrics HTTP listen address |
| `--web.metrics-path` / `WEB_METRICS_PATH` | `/metrics` | Metrics path |
| `--twc.timeout` / `TWC_TIMEOUT` | `5s` | HTTP timeout per API call |

One exporter instance maps to **one** wall connector. Run multiple releases for multiple chargers.

## API endpoints polled

| Path | Data |
|------|------|
| `/api/1/vitals` | Charging state, grid, per-phase V/I, temperatures, session energy |
| `/api/1/lifetime` | Total energy, charge starts, cycles, foldbacks |
| `/api/1/wifi_status` | Signal, RSSI, connectivity |
| `/api/1/version` | Firmware, serial, part number |

Lifetime JSON may contain invalid `"avg_startup_temp":nan`; the exporter sanitizes this before parsing.

## Metrics

All metrics use the prefix `tesla_wall_connector_`.

| Metric | Type | Description |
|--------|------|-------------|
| `scrape_success{endpoint}` | gauge | 1 on successful endpoint scrape |
| `scrape_duration_seconds{endpoint}` | gauge | Last scrape duration |
| `info{firmware_version,part_number,serial_number}` | gauge | Identity (value 1) |
| `wifi_info{ssid,mac,infra_ip}` | gauge | WiFi details (value 1) |
| `contactor_closed` | gauge | Contactor closed (0/1) |
| `vehicle_connected` | gauge | Vehicle connected (0/1) |
| `session_seconds` | gauge | Current session duration |
| `session_energy_watt_hours` | gauge | Energy this session (Wh) |
| `evse_state` | gauge | EVSE state code |
| `config_status` | gauge | Config status code |
| `current_alerts` | gauge | Count of active alerts |
| `alert_info{alert}` | gauge | One series per alert (value 1) |
| `grid_voltage_volts` | gauge | Grid voltage |
| `grid_frequency_hertz` | gauge | Grid frequency |
| `vehicle_current_amps` | gauge | Delivered current |
| `current_phase_*_amps` | gauge | Per-phase current |
| `voltage_phase_*_volts` | gauge | Per-phase voltage |
| `relay_coil_volts`, `pilot_*_volts`, `prox_volts` | gauge | Pilot / proximity |
| `*_temperature_celsius` | gauge | PCBA, handle, MCU temps |
| `input_thermopile_microvolts` | gauge | Thermopile reading |
| `uptime_seconds` | gauge | Uptime since restart |
| `lifetime_*_total` | counter | Lifetime counters |
| `avg_startup_temperature_celsius` | gauge | When present in lifetime API |
| `wifi_*` | gauge | Signal, RSSI, SNR, connected, internet |

## Development

```bash
make test
TWC_ADDRESS=192.168.1.100 make run
curl localhost:9859/metrics
```

### Release (local build)

```bash
make test
make docker-push VERSION=v0.1.1
make helm-package VERSION=v0.1.1
```

Push to a private registry by overriding `IMAGE`:

```bash
make docker-push VERSION=v0.1.1 IMAGE=registry.example.com/denislemire/tesla-wall-connector-exporter
```

## Troubleshooting

1. **No metrics / scrape_success 0** — From the same network as the exporter, run `curl http://<charger-ip>/api/1/vitals`. Guest WiFi or VLAN isolation blocks access.
2. **Timeout from Kubernetes** — Pod network may not reach the charger LAN; try `hostNetwork: true` on a node with LAN routing.
3. **Invalid JSON on lifetime** — Usually fixed automatically; report persistent firmware quirks in an issue.

## CI/CD

Automated build and deploy pipelines are **not enabled yet**. See [docs/CICD.md](docs/CICD.md) for the intended design and `ci/*.example` reference configs.

## License

Apache License 2.0 — see [LICENSE](LICENSE).
