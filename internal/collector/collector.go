package collector

import (
	"fmt"
	"strconv"

	"github.com/denislemire/tesla-wall-connector-exporter/internal/client"
	"github.com/denislemire/tesla-wall-connector-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector implements prometheus.Collector for one Wall Connector.
type Collector struct {
	client *client.Client
}

func New(c *client.Client) *Collector {
	return &Collector{client: c}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- metrics.ScrapeSuccess
	ch <- metrics.ScrapeDurationSeconds
	ch <- metrics.Info
	ch <- metrics.ContactorClosed
	ch <- metrics.VehicleConnected
	ch <- metrics.SessionSeconds
	ch <- metrics.SessionEnergyWh
	ch <- metrics.EvseState
	ch <- metrics.ConfigStatus
	ch <- metrics.CurrentAlerts
	ch <- metrics.AlertInfo
	ch <- metrics.GridVoltageV
	ch <- metrics.GridFrequencyHz
	ch <- metrics.VehicleCurrentAmps
	ch <- metrics.CurrentAPhaseAmps
	ch <- metrics.CurrentBPhaseAmps
	ch <- metrics.CurrentCPhaseAmps
	ch <- metrics.CurrentNeutralAmps
	ch <- metrics.VoltageAPhaseVolts
	ch <- metrics.VoltageBPhaseVolts
	ch <- metrics.VoltageCPhaseVolts
	ch <- metrics.RelayCoilVolts
	ch <- metrics.PilotHighVolts
	ch <- metrics.PilotLowVolts
	ch <- metrics.ProxVolts
	ch <- metrics.PcbaTempCelsius
	ch <- metrics.HandleTempCelsius
	ch <- metrics.McuTempCelsius
	ch <- metrics.InputThermopileMicrovolts
	ch <- metrics.UptimeSeconds
	ch <- metrics.LifetimeEnergyWh
	ch <- metrics.LifetimeChargeStartsTotal
	ch <- metrics.LifetimeChargingTimeSeconds
	ch <- metrics.LifetimeContactorCyclesTotal
	ch <- metrics.LifetimeContactorCyclesLoaded
	ch <- metrics.LifetimeConnectorCyclesTotal
	ch <- metrics.LifetimeAlertsTotal
	ch <- metrics.LifetimeThermalFoldbacksTotal
	ch <- metrics.LifetimeUptimeSecondsTotal
	ch <- metrics.AvgStartupTempCelsius
	ch <- metrics.WiFiSignalStrength
	ch <- metrics.WiFiRSSI
	ch <- metrics.WiFiSNR
	ch <- metrics.WiFiConnected
	ch <- metrics.InternetConnected
	ch <- metrics.WiFiInfo
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.collectVitals(ch)
	c.collectLifetime(ch)
	c.collectWiFi(ch)
	c.collectVersion(ch)
}

func (c *Collector) collectVitals(ch chan<- prometheus.Metric) {
	v, elapsed, err := c.client.FetchVitals()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(metrics.ScrapeSuccess, prometheus.GaugeValue, 0, client.EndpointVitals)
		ch <- prometheus.MustNewConstMetric(metrics.ScrapeDurationSeconds, prometheus.GaugeValue, elapsed, client.EndpointVitals)
		return
	}
	ch <- prometheus.MustNewConstMetric(metrics.ScrapeSuccess, prometheus.GaugeValue, 1, client.EndpointVitals)
	ch <- prometheus.MustNewConstMetric(metrics.ScrapeDurationSeconds, prometheus.GaugeValue, elapsed, client.EndpointVitals)

	ch <- prometheus.MustNewConstMetric(metrics.ContactorClosed, prometheus.GaugeValue, bool01(v.ContactorClosed))
	ch <- prometheus.MustNewConstMetric(metrics.VehicleConnected, prometheus.GaugeValue, bool01(v.VehicleConnected))
	ch <- prometheus.MustNewConstMetric(metrics.SessionSeconds, prometheus.GaugeValue, float64(v.SessionS))
	ch <- prometheus.MustNewConstMetric(metrics.SessionEnergyWh, prometheus.GaugeValue, v.SessionEnergyWh)
	ch <- prometheus.MustNewConstMetric(metrics.EvseState, prometheus.GaugeValue, float64(v.EvseState))
	ch <- prometheus.MustNewConstMetric(metrics.ConfigStatus, prometheus.GaugeValue, float64(v.ConfigStatus))
	ch <- prometheus.MustNewConstMetric(metrics.CurrentAlerts, prometheus.GaugeValue, float64(len(v.CurrentAlerts)))
	for _, a := range v.CurrentAlerts {
		ch <- prometheus.MustNewConstMetric(metrics.AlertInfo, prometheus.GaugeValue, 1, alertLabel(a))
	}
	ch <- prometheus.MustNewConstMetric(metrics.GridVoltageV, prometheus.GaugeValue, v.GridV)
	ch <- prometheus.MustNewConstMetric(metrics.GridFrequencyHz, prometheus.GaugeValue, v.GridHz)
	ch <- prometheus.MustNewConstMetric(metrics.VehicleCurrentAmps, prometheus.GaugeValue, v.VehicleCurrentA)
	ch <- prometheus.MustNewConstMetric(metrics.CurrentAPhaseAmps, prometheus.GaugeValue, v.CurrentAA)
	ch <- prometheus.MustNewConstMetric(metrics.CurrentBPhaseAmps, prometheus.GaugeValue, v.CurrentBA)
	ch <- prometheus.MustNewConstMetric(metrics.CurrentCPhaseAmps, prometheus.GaugeValue, v.CurrentCA)
	ch <- prometheus.MustNewConstMetric(metrics.CurrentNeutralAmps, prometheus.GaugeValue, v.CurrentNA)
	ch <- prometheus.MustNewConstMetric(metrics.VoltageAPhaseVolts, prometheus.GaugeValue, v.VoltageAV)
	ch <- prometheus.MustNewConstMetric(metrics.VoltageBPhaseVolts, prometheus.GaugeValue, v.VoltageBV)
	ch <- prometheus.MustNewConstMetric(metrics.VoltageCPhaseVolts, prometheus.GaugeValue, v.VoltageCV)
	ch <- prometheus.MustNewConstMetric(metrics.RelayCoilVolts, prometheus.GaugeValue, v.RelayCoilV)
	ch <- prometheus.MustNewConstMetric(metrics.PilotHighVolts, prometheus.GaugeValue, v.PilotHighV)
	ch <- prometheus.MustNewConstMetric(metrics.PilotLowVolts, prometheus.GaugeValue, v.PilotLowV)
	ch <- prometheus.MustNewConstMetric(metrics.ProxVolts, prometheus.GaugeValue, v.ProxV)
	ch <- prometheus.MustNewConstMetric(metrics.PcbaTempCelsius, prometheus.GaugeValue, v.PcbaTempC)
	ch <- prometheus.MustNewConstMetric(metrics.HandleTempCelsius, prometheus.GaugeValue, v.HandleTempC)
	ch <- prometheus.MustNewConstMetric(metrics.McuTempCelsius, prometheus.GaugeValue, v.McuTempC)
	ch <- prometheus.MustNewConstMetric(metrics.InputThermopileMicrovolts, prometheus.GaugeValue, v.InputThermopileUv)
	ch <- prometheus.MustNewConstMetric(metrics.UptimeSeconds, prometheus.GaugeValue, float64(v.UptimeS))
}

func (c *Collector) collectLifetime(ch chan<- prometheus.Metric) {
	lt, elapsed, err := c.client.FetchLifetime()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(metrics.ScrapeSuccess, prometheus.GaugeValue, 0, client.EndpointLifetime)
		ch <- prometheus.MustNewConstMetric(metrics.ScrapeDurationSeconds, prometheus.GaugeValue, elapsed, client.EndpointLifetime)
		return
	}
	ch <- prometheus.MustNewConstMetric(metrics.ScrapeSuccess, prometheus.GaugeValue, 1, client.EndpointLifetime)
	ch <- prometheus.MustNewConstMetric(metrics.ScrapeDurationSeconds, prometheus.GaugeValue, elapsed, client.EndpointLifetime)

	ch <- prometheus.MustNewConstMetric(metrics.LifetimeEnergyWh, prometheus.CounterValue, float64(lt.EnergyWh))
	ch <- prometheus.MustNewConstMetric(metrics.LifetimeChargeStartsTotal, prometheus.CounterValue, float64(lt.ChargeStarts))
	ch <- prometheus.MustNewConstMetric(metrics.LifetimeChargingTimeSeconds, prometheus.CounterValue, float64(lt.ChargingTimeS))
	ch <- prometheus.MustNewConstMetric(metrics.LifetimeContactorCyclesTotal, prometheus.CounterValue, float64(lt.ContactorCycles))
	ch <- prometheus.MustNewConstMetric(metrics.LifetimeContactorCyclesLoaded, prometheus.CounterValue, float64(lt.ContactorCyclesLoaded))
	ch <- prometheus.MustNewConstMetric(metrics.LifetimeConnectorCyclesTotal, prometheus.CounterValue, float64(lt.ConnectorCycles))
	ch <- prometheus.MustNewConstMetric(metrics.LifetimeAlertsTotal, prometheus.CounterValue, float64(lt.AlertCount))
	ch <- prometheus.MustNewConstMetric(metrics.LifetimeThermalFoldbacksTotal, prometheus.CounterValue, float64(lt.ThermalFoldbacks))
	ch <- prometheus.MustNewConstMetric(metrics.LifetimeUptimeSecondsTotal, prometheus.CounterValue, float64(lt.UptimeS))
	if lt.AvgStartupTemp != nil {
		ch <- prometheus.MustNewConstMetric(metrics.AvgStartupTempCelsius, prometheus.GaugeValue, *lt.AvgStartupTemp)
	}
}

func (c *Collector) collectWiFi(ch chan<- prometheus.Metric) {
	w, elapsed, err := c.client.FetchWiFiStatus()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(metrics.ScrapeSuccess, prometheus.GaugeValue, 0, client.EndpointWiFiStatus)
		ch <- prometheus.MustNewConstMetric(metrics.ScrapeDurationSeconds, prometheus.GaugeValue, elapsed, client.EndpointWiFiStatus)
		return
	}
	ch <- prometheus.MustNewConstMetric(metrics.ScrapeSuccess, prometheus.GaugeValue, 1, client.EndpointWiFiStatus)
	ch <- prometheus.MustNewConstMetric(metrics.ScrapeDurationSeconds, prometheus.GaugeValue, elapsed, client.EndpointWiFiStatus)

	ch <- prometheus.MustNewConstMetric(metrics.WiFiSignalStrength, prometheus.GaugeValue, float64(w.WiFiSignalStrength))
	ch <- prometheus.MustNewConstMetric(metrics.WiFiRSSI, prometheus.GaugeValue, float64(w.WiFiRSSI))
	ch <- prometheus.MustNewConstMetric(metrics.WiFiSNR, prometheus.GaugeValue, float64(w.WiFiSNR))
	ch <- prometheus.MustNewConstMetric(metrics.WiFiConnected, prometheus.GaugeValue, bool01(w.WiFiConnected))
	ch <- prometheus.MustNewConstMetric(metrics.InternetConnected, prometheus.GaugeValue, bool01(w.Internet))
	ch <- prometheus.MustNewConstMetric(metrics.WiFiInfo, prometheus.GaugeValue, 1, w.WiFiSSID, w.WiFiMAC, w.WiFiInfraIP)
}

func (c *Collector) collectVersion(ch chan<- prometheus.Metric) {
	v, elapsed, err := c.client.FetchVersion()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(metrics.ScrapeSuccess, prometheus.GaugeValue, 0, client.EndpointVersion)
		ch <- prometheus.MustNewConstMetric(metrics.ScrapeDurationSeconds, prometheus.GaugeValue, elapsed, client.EndpointVersion)
		return
	}
	ch <- prometheus.MustNewConstMetric(metrics.ScrapeSuccess, prometheus.GaugeValue, 1, client.EndpointVersion)
	ch <- prometheus.MustNewConstMetric(metrics.ScrapeDurationSeconds, prometheus.GaugeValue, elapsed, client.EndpointVersion)
	ch <- prometheus.MustNewConstMetric(metrics.Info, prometheus.GaugeValue, 1, v.FirmwareVersion, v.PartNumber, v.SerialNumber)
}

func bool01(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

func alertLabel(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	default:
		return fmt.Sprint(v)
	}
}
