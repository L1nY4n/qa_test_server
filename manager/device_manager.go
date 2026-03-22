package manager

import (
	"sort"
	"strings"
	"sync"
	"time"

	"qa_test_server/model"
)

var ManagerGlabal Manager

type DeviceStats struct {
	Total               int       `json:"total"`
	Online              int       `json:"online"`
	Offline             int       `json:"offline"`
	ActiveWithinSeconds int64     `json:"activeWithinSeconds"`
	GeneratedAt         time.Time `json:"generatedAt"`
}

type DeviceSummary struct {
	Sn               string    `json:"Sn"`
	Name             string    `json:"Name"`
	Group            string    `json:"Group"`
	Model            string    `json:"Model"`
	PN               string    `json:"PN"`
	Last_rx_time     time.Time `json:"Last_rx_time"`
	Hardware_bate    uint16    `json:"Hardware_bate"`
	Uptime           []uint16  `json:"Uptime"`
	Uptime_seconds   uint32    `json:"Uptime_seconds"`
	Pump_count       int       `json:"Pump_count"`
	Temp_count       int       `json:"Temp_count"`
	Laser_status     uint16    `json:"Laser_status"`
	Laser_ready      uint16    `json:"Laser_ready"`
	Laser_wavelength uint16    `json:"Laser_wavelength"`
	Online           bool      `json:"Online"`
}

type Manager struct {
	devices sync.Map
}

func (m *Manager) ClearAll() int {
	cleared := 0
	m.devices.Range(func(key, _ interface{}) bool {
		m.devices.Delete(key)
		cleared++
		return true
	})
	return cleared
}

func (m *Manager) Update(device model.Device) {
	sn := strings.TrimSpace(device.Sn)
	if sn == "" {
		return
	}

	if prev, ok := m.Get(sn); ok && prev != nil {
		changes := BuildSystemParamChanges(sn, prev.Packet.Femto_holding_reg, device.Packet.Femto_holding_reg, device.Last_rx_time)
		ParamChangeManagerGlobal.Enqueue(changes)
	}
	m.Set(sn, &device)
}

func (m *Manager) List() []model.Device {
	var devicelist []model.Device
	m.devices.Range(func(_, dev interface{}) bool {
		device := dev.(*model.Device)
		devicelist = append(devicelist, *device)
		return true
	})
	sort.Slice(devicelist, func(i, j int) bool {
		return strings.Compare(devicelist[i].Sn, devicelist[j].Sn) < 0
	})
	return devicelist
}

func (m *Manager) Query(keyword, group string, onlineOnly bool, activeWithin time.Duration, offset, limit int) ([]model.Device, int) {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	group = strings.ToLower(strings.TrimSpace(group))
	offset, limit = normalizeRange(offset, limit)

	all := m.List()
	filtered := make([]model.Device, 0, len(all))
	for _, item := range all {
		if group != "" && strings.ToLower(strings.TrimSpace(item.Group)) != group {
			continue
		}
		if keyword != "" {
			modelText := strings.ToLower(decodeASCIIBytes(item.Packet.Femto_holding_reg.Laser_para.Laser_info.Model[:]))
			pnText := strings.ToLower(decodeASCIIBytes(item.Packet.Femto_holding_reg.Laser_para.Laser_info.PN[:]))
			if !strings.Contains(strings.ToLower(item.Sn), keyword) &&
				!strings.Contains(strings.ToLower(item.Name), keyword) &&
				!strings.Contains(modelText, keyword) &&
				!strings.Contains(pnText, keyword) &&
				!strings.Contains(strings.ToLower(item.Group), keyword) {
				continue
			}
		}
		if onlineOnly && !isOnline(item, activeWithin) {
			continue
		}
		filtered = append(filtered, item)
	}

	total := len(filtered)
	if offset >= total {
		return []model.Device{}, total
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return filtered[offset:end], total
}

func (m *Manager) Stats(activeWithin time.Duration) DeviceStats {
	stats := DeviceStats{
		ActiveWithinSeconds: int64(activeWithin.Seconds()),
		GeneratedAt:         time.Now(),
	}
	m.devices.Range(func(_, dev interface{}) bool {
		stats.Total++
		device := dev.(*model.Device)
		if isOnline(*device, activeWithin) {
			stats.Online++
		}
		return true
	})
	stats.Offline = stats.Total - stats.Online
	return stats
}

func (m *Manager) QuerySummary(keyword, group string, onlineOnly bool, activeWithin time.Duration, offset, limit int) ([]DeviceSummary, int) {
	items, total := m.Query(keyword, group, onlineOnly, activeWithin, offset, limit)
	summaries := make([]DeviceSummary, 0, len(items))
	for _, item := range items {
		summaries = append(summaries, deviceToSummary(item, activeWithin))
	}
	return summaries, total
}

func (m *Manager) Get(sn string) (*model.Device, bool) {
	val, ok := m.devices.Load(strings.TrimSpace(sn))
	if !ok {
		return nil, false
	}
	device, ok := val.(*model.Device)
	return device, ok
}

func (m *Manager) Set(sn string, dev *model.Device) {
	m.devices.Store(sn, dev)
}

func (m *Manager) Delete(sn string) {
	m.devices.Delete(sn)
}

func normalizeRange(offset, limit int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	return offset, limit
}

func isOnline(device model.Device, activeWithin time.Duration) bool {
	if activeWithin <= 0 {
		activeWithin = 30 * time.Second
	}
	if device.Last_rx_time.IsZero() {
		return false
	}
	return time.Since(device.Last_rx_time) <= activeWithin
}

func deviceToSummary(device model.Device, activeWithin time.Duration) DeviceSummary {
	modelText := decodeASCIIBytes(device.Packet.Femto_holding_reg.Laser_para.Laser_info.Model[:])
	pnText := decodeASCIIBytes(device.Packet.Femto_holding_reg.Laser_para.Laser_info.PN[:])
	uptime := device.Packet.Femto_input_reg.Time.Uptime
	uptimeSeconds := (uint32(uptime[0]) << 16) | uint32(uptime[1])
	return DeviceSummary{
		Sn:               device.Sn,
		Name:             device.Name,
		Group:            device.Group,
		Model:            modelText,
		PN:               pnText,
		Last_rx_time:     device.Last_rx_time,
		Hardware_bate:    device.Packet.Femto_input_reg.Bate.Hardware_bate,
		Uptime:           []uint16{uptime[0], uptime[1]},
		Uptime_seconds:   uptimeSeconds,
		Pump_count:       len(device.Packet.Femto_input_reg.Mon.Pump_mon),
		Temp_count:       len(device.Packet.Femto_input_reg.Mon.Temp),
		Laser_status:     device.Packet.Femto_input_reg.Status,
		Laser_ready:      device.Packet.Femto_holding_reg.User_para.Laser_ready,
		Laser_wavelength: device.Packet.Femto_holding_reg.User_para.Laser_wavelength,
		Online:           isOnline(device, activeWithin),
	}
}

func decodeASCIIBytes(src []uint8) string {
	if len(src) == 0 {
		return ""
	}
	buf := make([]byte, 0, len(src))
	for _, b := range src {
		if b == 0 {
			break
		}
		if b < 32 || b > 126 {
			continue
		}
		buf = append(buf, b)
	}
	return string(buf)
}
