package client

// Vitals from GET /api/1/vitals
type Vitals struct {
	ContactorClosed  bool          `json:"contactor_closed"`
	VehicleConnected bool          `json:"vehicle_connected"`
	SessionS         int           `json:"session_s"`
	GridV            float64       `json:"grid_v"`
	GridHz           float64       `json:"grid_hz"`
	VehicleCurrentA  float64       `json:"vehicle_current_a"`
	CurrentAA        float64       `json:"currentA_a"`
	CurrentBA        float64       `json:"currentB_a"`
	CurrentCA        float64       `json:"currentC_a"`
	CurrentNA        float64       `json:"currentN_a"`
	VoltageAV        float64       `json:"voltageA_v"`
	VoltageBV        float64       `json:"voltageB_v"`
	VoltageCV        float64       `json:"voltageC_v"`
	RelayCoilV       float64       `json:"relay_coil_v"`
	PcbaTempC        float64       `json:"pcba_temp_c"`
	HandleTempC      float64       `json:"handle_temp_c"`
	McuTempC         float64       `json:"mcu_temp_c"`
	UptimeS          int           `json:"uptime_s"`
	InputThermopileUv float64      `json:"input_thermopile_uv"`
	ProxV            float64       `json:"prox_v"`
	PilotHighV       float64       `json:"pilot_high_v"`
	PilotLowV        float64       `json:"pilot_low_v"`
	SessionEnergyWh  float64       `json:"session_energy_wh"`
	ConfigStatus     int           `json:"config_status"`
	EvseState        int           `json:"evse_state"`
	CurrentAlerts    []interface{} `json:"current_alerts"`
}

// Lifetime from GET /api/1/lifetime
type Lifetime struct {
	ContactorCycles       int      `json:"contactor_cycles"`
	ContactorCyclesLoaded int      `json:"contactor_cycles_loaded"`
	AlertCount            int      `json:"alert_count"`
	ThermalFoldbacks      int      `json:"thermal_foldbacks"`
	AvgStartupTemp        *float64 `json:"avg_startup_temp"`
	ChargeStarts          int      `json:"charge_starts"`
	EnergyWh              int      `json:"energy_wh"`
	ConnectorCycles       int      `json:"connector_cycles"`
	UptimeS               int      `json:"uptime_s"`
	ChargingTimeS         int      `json:"charging_time_s"`
}

// WiFiStatus from GET /api/1/wifi_status
type WiFiStatus struct {
	WiFiSSID            string `json:"wifi_ssid"`
	WiFiSignalStrength  int    `json:"wifi_signal_strength"`
	WiFiRSSI            int    `json:"wifi_rssi"`
	WiFiSNR             int    `json:"wifi_snr"`
	WiFiConnected       bool   `json:"wifi_connected"`
	WiFiInfraIP         string `json:"wifi_infra_ip"`
	Internet            bool   `json:"internet"`
	WiFiMAC             string `json:"wifi_mac"`
}

// Version from GET /api/1/version
type Version struct {
	FirmwareVersion string `json:"firmware_version"`
	PartNumber      string `json:"part_number"`
	SerialNumber    string `json:"serial_number"`
}

const (
	EndpointVitals     = "vitals"
	EndpointLifetime   = "lifetime"
	EndpointWiFiStatus = "wifi_status"
	EndpointVersion    = "version"
)
