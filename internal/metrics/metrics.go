package metrics

import "github.com/prometheus/client_golang/prometheus"

const namespace = "tesla_wall_connector"

func fq(name string) string {
	return prometheus.BuildFQName(namespace, "", name)
}

var (
	ScrapeSuccess = prometheus.NewDesc(
		fq("scrape_success"),
		"1 if the last scrape of the endpoint succeeded.",
		[]string{"endpoint"},
		nil,
	)
	ScrapeDurationSeconds = prometheus.NewDesc(
		fq("scrape_duration_seconds"),
		"Duration of the last scrape per endpoint in seconds.",
		[]string{"endpoint"},
		nil,
	)

	// Version / info
	Info = prometheus.NewDesc(
		fq("info"),
		"Wall Connector firmware and identity (value is always 1).",
		[]string{"firmware_version", "part_number", "serial_number"},
		nil,
	)

	// Vitals — charging state
	ContactorClosed  = gauge("contactor_closed", "Whether the contactor is closed (1=yes).")
	VehicleConnected = gauge("vehicle_connected", "Whether a vehicle is connected (1=yes).")
	SessionSeconds   = gauge("session_seconds", "Current charging session duration in seconds.")
	SessionEnergyWh  = gauge("session_energy_watt_hours", "Energy delivered in the current session in watt-hours.")
	EvseState        = gauge("evse_state", "EVSE state code.")
	ConfigStatus     = gauge("config_status", "Configuration status code.")
	CurrentAlerts    = gauge("current_alerts", "Number of current alerts.")
	AlertInfo        = prometheus.NewDesc(
		fq("alert_info"),
		"Current alert present (value 1).",
		[]string{"alert"},
		nil,
	)

	// Vitals — electrical
	GridVoltageV       = gauge("grid_voltage_volts", "Grid voltage in volts.")
	GridFrequencyHz    = gauge("grid_frequency_hertz", "Grid frequency in hertz.")
	VehicleCurrentAmps = gauge("vehicle_current_amps", "Current delivered to the vehicle in amps.")
	CurrentAPhaseAmps  = gauge("current_phase_a_amps", "Phase A current in amps.")
	CurrentBPhaseAmps  = gauge("current_phase_b_amps", "Phase B current in amps.")
	CurrentCPhaseAmps  = gauge("current_phase_c_amps", "Phase C current in amps.")
	CurrentNeutralAmps = gauge("current_neutral_amps", "Neutral current in amps.")
	VoltageAPhaseVolts = gauge("voltage_phase_a_volts", "Phase A voltage in volts.")
	VoltageBPhaseVolts = gauge("voltage_phase_b_volts", "Phase B voltage in volts.")
	VoltageCPhaseVolts = gauge("voltage_phase_c_volts", "Phase C voltage in volts.")
	RelayCoilVolts     = gauge("relay_coil_volts", "Relay coil voltage in volts.")
	PilotHighVolts     = gauge("pilot_high_volts", "Pilot high voltage in volts.")
	PilotLowVolts      = gauge("pilot_low_volts", "Pilot low voltage in volts.")
	ProxVolts          = gauge("prox_volts", "Proximity voltage in volts.")

	// Vitals — temperature
	PcbaTempCelsius         = gauge("pcba_temperature_celsius", "PCBA temperature in Celsius.")
	HandleTempCelsius       = gauge("handle_temperature_celsius", "Handle temperature in Celsius.")
	McuTempCelsius          = gauge("mcu_temperature_celsius", "MCU temperature in Celsius.")
	InputThermopileMicrovolts = gauge("input_thermopile_microvolts", "Input thermopile reading in microvolts.")
	UptimeSeconds           = gauge("uptime_seconds", "Seconds since last restart.")

	// Lifetime counters
	LifetimeEnergyWh              = counter("lifetime_energy_watt_hours_total", "Total energy delivered in watt-hours.")
	LifetimeChargeStartsTotal     = counter("lifetime_charge_starts_total", "Total number of charge sessions started.")
	LifetimeChargingTimeSeconds   = counter("lifetime_charging_time_seconds_total", "Total time spent charging in seconds.")
	LifetimeContactorCyclesTotal  = counter("lifetime_contactor_cycles_total", "Total contactor cycles.")
	LifetimeContactorCyclesLoaded = counter("lifetime_contactor_cycles_loaded_total", "Contactor cycles while loaded.")
	LifetimeConnectorCyclesTotal  = counter("lifetime_connector_cycles_total", "Connector insertion cycles.")
	LifetimeAlertsTotal           = counter("lifetime_alerts_total", "Lifetime alert count.")
	LifetimeThermalFoldbacksTotal = counter("lifetime_thermal_foldbacks_total", "Thermal foldback events.")
	LifetimeUptimeSecondsTotal    = counter("lifetime_uptime_seconds_total", "Lifetime uptime in seconds.")
	AvgStartupTempCelsius         = gauge("avg_startup_temperature_celsius", "Average startup temperature in Celsius.")

	// WiFi
	WiFiSignalStrength = gauge("wifi_signal_strength_percent", "WiFi signal strength percentage.")
	WiFiRSSI           = gauge("wifi_rssi_dbm", "WiFi RSSI in dBm.")
	WiFiSNR            = gauge("wifi_snr_db", "WiFi signal-to-noise ratio in dB.")
	WiFiConnected      = gauge("wifi_connected", "WiFi connected (1=yes).")
	InternetConnected  = gauge("internet_connected", "Internet reachable from wall connector (1=yes).")
	WiFiInfo = prometheus.NewDesc(
		fq("wifi_info"),
		"WiFi connection details (value is always 1).",
		[]string{"ssid", "mac", "infra_ip"},
		nil,
	)
)

func gauge(name, help string) *prometheus.Desc {
	return prometheus.NewDesc(fq(name), help, nil, nil)
}

func counter(name, help string) *prometheus.Desc {
	return prometheus.NewDesc(fq(name), help, nil, nil)
}
