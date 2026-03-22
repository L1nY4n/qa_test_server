package model

import "time"

// DeviceHistory stores sampled device data for historical replay.
// One row represents one device snapshot at one minute granularity.
type DeviceHistory struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DeviceSn      string    `gorm:"size:64;index:idx_device_time,priority:1;uniqueIndex:uk_device_minute,priority:1;not null" json:"deviceSn"`
	DeviceName    string    `gorm:"size:128" json:"deviceName"`
	SampledAt     time.Time `gorm:"index:idx_device_time,priority:2;index;uniqueIndex:uk_device_minute,priority:2;not null" json:"sampledAt"`
	Online        bool      `gorm:"index" json:"online"`
	HardwareBate  uint16    `json:"hardwareBate"`
	UptimeSeconds uint32    `json:"uptimeSeconds"`
	PumpCount     int       `json:"pumpCount"`
	TempCount     int       `json:"tempCount"`
	RawJSON       string    `gorm:"type:longtext" json:"rawJson"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (DeviceHistory) TableName() string {
	return "qa_device_history"
}

type DeviceHistoryPoint struct {
	SampledAt     time.Time `json:"sampledAt"`
	Online        bool      `json:"online"`
	HardwareBate  uint16    `json:"hardwareBate"`
	UptimeSeconds uint32    `json:"uptimeSeconds"`
	PumpCount     int       `json:"pumpCount"`
	TempCount     int       `json:"tempCount"`
}

// DeviceHistoryTimelinePoint is optimized for long-range trend plotting.
// Temperature is converted to Celsius (raw / 10), voltage/current keep raw values.
type DeviceHistoryTimelinePoint struct {
	SampledAt  time.Time `json:"sampledAt"`
	Online     bool      `json:"online"`
	TempAvg    float64   `json:"tempAvg"`
	TempMax    float64   `json:"tempMax"`
	VoltageAvg float64   `json:"voltageAvg"`
	VoltageMax float64   `json:"voltageMax"`
	CurrentAvg float64   `json:"currentAvg"`
	CurrentMax float64   `json:"currentMax"`
}

// DeviceHistoryMetricPoint is aligned by sampledAt for selected channel metrics.
// Nil metric fields mean channel index out of range or source data missing.
type DeviceHistoryMetricPoint struct {
	SampledAt time.Time `json:"sampledAt"`
	Online    bool      `json:"online"`
	Temp      *float64  `json:"temp"`
	Voltage   *float64  `json:"voltage"`
	Current   *float64  `json:"current"`
}
