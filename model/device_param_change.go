package model

import "time"

// DeviceParamChange stores one system parameter change event for a device.
type DeviceParamChange struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DeviceSn  string    `gorm:"size:64;index:idx_change_device_time,priority:1;not null" json:"deviceSn"`
	ParamPath string    `gorm:"size:512;index" json:"paramPath"`
	OldValue  string    `gorm:"type:text" json:"oldValue"`
	NewValue  string    `gorm:"type:text" json:"newValue"`
	ChangedAt time.Time `gorm:"index:idx_change_device_time,priority:2;index;not null" json:"changedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

func (DeviceParamChange) TableName() string {
	return "qa_device_param_changes"
}
