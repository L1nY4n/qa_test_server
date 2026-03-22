package manager

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"qa_test_server/db"
	"qa_test_server/model"

	"gorm.io/gorm/clause"
)

var HistoryManagerGlobal = &HistoryManager{}

type HistoryQuery struct {
	DeviceSn string
	Start    time.Time
	End      time.Time
	Offset   int
	Limit    int
}

type HistoryMetricQuery struct {
	DeviceSn     string
	Start        time.Time
	End          time.Time
	Offset       int
	Limit        int
	TempIndex    int
	VoltageIndex int
	CurrentIndex int
}

type HistoryManager struct {
	deviceManager  *Manager
	retention      time.Duration
	sampleInterval time.Duration

	recoveredCount int
	recoveredAt    time.Time

	stopCh   chan struct{}
	wg       sync.WaitGroup
	initOnce sync.Once
	mu       sync.RWMutex
}

type historyTimelineRawDevice struct {
	Packet struct {
		FemtoInputReg struct {
			Mon struct {
				Temp    []uint16 `json:"Temp"`
				Vol     []uint16 `json:"Vol"`
				PumpMon []struct {
					ActualCur uint16 `json:"Actual_cur"`
					FpgaCur   uint16 `json:"Fpga_cur"`
				} `json:"Pump_mon"`
			} `json:"Mon"`
		} `json:"Femto_input_reg"`
	} `json:"Packet"`
}

func (m *HistoryManager) Init(deviceManager *Manager, sampleInterval, retention time.Duration) error {
	if db.DB == nil {
		return errors.New("database is unavailable")
	}
	if deviceManager == nil {
		return errors.New("device manager is required")
	}
	if sampleInterval <= 0 {
		sampleInterval = time.Minute
	}
	if retention <= 0 {
		retention = 10 * 24 * time.Hour
	}

	if err := db.DB.AutoMigrate(&model.DeviceHistory{}); err != nil {
		return err
	}

	recoveredCount := 0
	recoveredAt := time.Now()
	if c, err := m.restoreLatestDevices(context.Background(), deviceManager); err == nil {
		recoveredCount = c
	} else {
		log.Printf("restore device cache from history failed: %v", err)
	}

	m.initOnce.Do(func() {
		m.deviceManager = deviceManager
		m.sampleInterval = sampleInterval
		m.retention = retention
		m.mu.Lock()
		m.recoveredCount = recoveredCount
		m.recoveredAt = recoveredAt
		m.mu.Unlock()
		m.stopCh = make(chan struct{})
		m.wg.Add(1)
		go m.loop()
	})

	return nil
}

type RecoveryStats struct {
	RecoveredCount int       `json:"recoveredCount"`
	RecoveredAt    time.Time `json:"recoveredAt"`
}

func (m *HistoryManager) Stats() RecoveryStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return RecoveryStats{
		RecoveredCount: m.recoveredCount,
		RecoveredAt:    m.recoveredAt,
	}
}

func (m *HistoryManager) loop() {
	defer m.wg.Done()

	next := alignToInterval(time.Now(), m.sampleInterval)
	timer := time.NewTimer(time.Until(next))
	defer timer.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-timer.C:
			if err := m.FlushOnce(context.Background()); err != nil {
				log.Printf("history flush failed: %v", err)
			}
			if err := m.PurgeExpired(context.Background()); err != nil {
				log.Printf("history purge failed: %v", err)
			}

			next = next.Add(m.sampleInterval)
			wait := time.Until(next)
			if wait < 0 {
				next = alignToInterval(time.Now(), m.sampleInterval)
				wait = time.Until(next)
			}
			timer.Reset(wait)
		}
	}
}

func (m *HistoryManager) Stop() {
	if m.stopCh == nil {
		return
	}
	select {
	case <-m.stopCh:
		return
	default:
		close(m.stopCh)
	}
	m.wg.Wait()
}

func (m *HistoryManager) FlushOnce(ctx context.Context) error {
	if db.DB == nil || m.deviceManager == nil {
		return nil
	}

	devices := m.deviceManager.List()
	if len(devices) == 0 {
		return nil
	}

	sampledAt := alignToInterval(time.Now(), m.sampleInterval)
	records := make([]model.DeviceHistory, 0, len(devices))
	for _, dev := range devices {
		sn := strings.TrimSpace(dev.Sn)
		if sn == "" {
			continue
		}

		payload, err := json.Marshal(dev)
		if err != nil {
			continue
		}

		uptime := dev.Packet.Femto_input_reg.Time.Uptime
		seconds := uint32(uptime[0])*65536 + uint32(uptime[1])
		records = append(records, model.DeviceHistory{
			DeviceSn:      sn,
			DeviceName:    strings.TrimSpace(dev.Name),
			SampledAt:     sampledAt,
			Online:        !dev.Last_rx_time.IsZero() && time.Since(dev.Last_rx_time) <= 30*time.Second,
			HardwareBate:  dev.Packet.Femto_input_reg.Bate.Hardware_bate,
			UptimeSeconds: seconds,
			PumpCount:     len(dev.Packet.Femto_input_reg.Mon.Pump_mon),
			TempCount:     len(dev.Packet.Femto_input_reg.Mon.Temp),
			RawJSON:       string(payload),
		})
	}
	if len(records) == 0 {
		return nil
	}

	return db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "device_sn"},
			{Name: "sampled_at"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"device_name",
			"online",
			"hardware_bate",
			"uptime_seconds",
			"pump_count",
			"temp_count",
			"raw_json",
		}),
	}).Create(&records).Error
}

func (m *HistoryManager) PurgeExpired(ctx context.Context) error {
	if db.DB == nil {
		return nil
	}
	deadline := time.Now().Add(-m.retention)
	return db.DB.WithContext(ctx).Where("sampled_at < ?", deadline).Delete(&model.DeviceHistory{}).Error
}

func (m *HistoryManager) ClearAll(ctx context.Context) (int64, error) {
	if db.DB == nil {
		return 0, errors.New("database is unavailable")
	}
	res := db.DB.WithContext(ctx).Where("1 = 1").Delete(&model.DeviceHistory{})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (m *HistoryManager) Query(q HistoryQuery) ([]model.DeviceHistory, int64, error) {
	if db.DB == nil {
		return nil, 0, errors.New("database is unavailable")
	}

	sn := strings.TrimSpace(q.DeviceSn)
	if sn == "" {
		return nil, 0, errors.New("device sn is required")
	}

	if q.Limit <= 0 {
		q.Limit = 120
	}
	if q.Limit > 2000 {
		q.Limit = 2000
	}
	if q.Offset < 0 {
		q.Offset = 0
	}

	now := time.Now()
	end := q.End
	if end.IsZero() {
		end = now
	}
	start := q.Start
	if start.IsZero() {
		start = end.Add(-24 * time.Hour)
	}
	if start.After(end) {
		start, end = end, start
	}

	maxRange := m.retention
	if maxRange <= 0 {
		maxRange = 10 * 24 * time.Hour
	}
	if end.Sub(start) > maxRange {
		start = end.Add(-maxRange)
	}
	minStart := now.Add(-maxRange)
	if start.Before(minStart) {
		start = minStart
	}

	query := db.DB.Model(&model.DeviceHistory{}).
		Where("device_sn = ? AND sampled_at >= ? AND sampled_at <= ?", sn, start, end)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := make([]model.DeviceHistory, 0, q.Limit)
	if err := query.Order("sampled_at DESC").Offset(q.Offset).Limit(q.Limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (m *HistoryManager) QueryPoints(q HistoryQuery) ([]model.DeviceHistoryPoint, int64, error) {
	items, total, err := m.Query(q)
	if err != nil {
		return nil, 0, err
	}
	points := make([]model.DeviceHistoryPoint, 0, len(items))
	for _, item := range items {
		points = append(points, model.DeviceHistoryPoint{
			SampledAt:     item.SampledAt,
			Online:        item.Online,
			HardwareBate:  item.HardwareBate,
			UptimeSeconds: item.UptimeSeconds,
			PumpCount:     item.PumpCount,
			TempCount:     item.TempCount,
		})
	}
	return points, total, nil
}

func (m *HistoryManager) QueryTimeline(q HistoryQuery, maxPoints int) ([]model.DeviceHistoryTimelinePoint, int64, error) {
	if db.DB == nil {
		return nil, 0, errors.New("database is unavailable")
	}

	sn := strings.TrimSpace(q.DeviceSn)
	if sn == "" {
		return nil, 0, errors.New("device sn is required")
	}

	if q.Limit <= 0 {
		q.Limit = 16000
	}
	if q.Limit > 20000 {
		q.Limit = 20000
	}
	if q.Offset < 0 {
		q.Offset = 0
	}

	now := time.Now()
	end := q.End
	if end.IsZero() {
		end = now
	}
	start := q.Start
	if start.IsZero() {
		start = end.Add(-24 * time.Hour)
	}
	if start.After(end) {
		start, end = end, start
	}

	maxRange := m.retention
	if maxRange <= 0 {
		maxRange = 10 * 24 * time.Hour
	}
	if end.Sub(start) > maxRange {
		start = end.Add(-maxRange)
	}
	minStart := now.Add(-maxRange)
	if start.Before(minStart) {
		start = minStart
	}

	query := db.DB.Model(&model.DeviceHistory{}).
		Where("device_sn = ? AND sampled_at >= ? AND sampled_at <= ?", sn, start, end)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	rows := make([]model.DeviceHistory, 0, q.Limit)
	if err := query.Order("sampled_at ASC").Offset(q.Offset).Limit(q.Limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	points := make([]model.DeviceHistoryTimelinePoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, timelinePointFromHistory(row))
	}

	if maxPoints <= 0 {
		maxPoints = 4000
	}
	if maxPoints > 20000 {
		maxPoints = 20000
	}
	points = downsampleTimelinePoints(points, maxPoints)
	return points, total, nil
}

func (m *HistoryManager) QueryMetrics(q HistoryMetricQuery, maxPoints int) ([]model.DeviceHistoryMetricPoint, int64, error) {
	if db.DB == nil {
		return nil, 0, errors.New("database is unavailable")
	}

	sn := strings.TrimSpace(q.DeviceSn)
	if sn == "" {
		return nil, 0, errors.New("device sn is required")
	}

	if q.Limit <= 0 {
		q.Limit = 16000
	}
	if q.Limit > 20000 {
		q.Limit = 20000
	}
	if q.Offset < 0 {
		q.Offset = 0
	}

	now := time.Now()
	end := q.End
	if end.IsZero() {
		end = now
	}
	start := q.Start
	if start.IsZero() {
		start = end.Add(-24 * time.Hour)
	}
	if start.After(end) {
		start, end = end, start
	}

	maxRange := m.retention
	if maxRange <= 0 {
		maxRange = 10 * 24 * time.Hour
	}
	if end.Sub(start) > maxRange {
		start = end.Add(-maxRange)
	}
	minStart := now.Add(-maxRange)
	if start.Before(minStart) {
		start = minStart
	}

	tempIndex := normalizeMetricIndex(q.TempIndex)
	voltageIndex := normalizeMetricIndex(q.VoltageIndex)
	currentIndex := normalizeMetricIndex(q.CurrentIndex)

	query := db.DB.Model(&model.DeviceHistory{}).
		Where("device_sn = ? AND sampled_at >= ? AND sampled_at <= ?", sn, start, end)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	rows := make([]model.DeviceHistory, 0, q.Limit)
	if err := query.Order("sampled_at ASC").Offset(q.Offset).Limit(q.Limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	points := make([]model.DeviceHistoryMetricPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, metricPointFromHistory(row, tempIndex, voltageIndex, currentIndex))
	}

	if maxPoints <= 0 {
		maxPoints = 4000
	}
	if maxPoints > 20000 {
		maxPoints = 20000
	}
	points = downsampleMetricPoints(points, maxPoints)
	return points, total, nil
}

func timelinePointFromHistory(row model.DeviceHistory) model.DeviceHistoryTimelinePoint {
	point := model.DeviceHistoryTimelinePoint{
		SampledAt: row.SampledAt,
		Online:    row.Online,
	}

	raw := strings.TrimSpace(row.RawJSON)
	if raw == "" {
		return point
	}

	parsed := historyTimelineRawDevice{}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return point
	}

	mon := parsed.Packet.FemtoInputReg.Mon
	tempAvgRaw, tempMaxRaw := averageAndMaxPositiveU16(mon.Temp)
	voltageAvg, voltageMax := averageAndMaxPositiveU16(mon.Vol)

	currents := make([]uint16, 0, len(mon.PumpMon))
	for _, pump := range mon.PumpMon {
		if pump.ActualCur > 0 {
			currents = append(currents, pump.ActualCur)
		} else if pump.FpgaCur > 0 {
			currents = append(currents, pump.FpgaCur)
		}
	}
	currentAvg, currentMax := averageAndMaxPositiveU16(currents)

	point.TempAvg = roundFloat(tempAvgRaw/10.0, 2)
	point.TempMax = roundFloat(tempMaxRaw/10.0, 2)
	point.VoltageAvg = roundFloat(voltageAvg, 2)
	point.VoltageMax = roundFloat(voltageMax, 2)
	point.CurrentAvg = roundFloat(currentAvg, 2)
	point.CurrentMax = roundFloat(currentMax, 2)
	return point
}

func metricPointFromHistory(row model.DeviceHistory, tempIndex, voltageIndex, currentIndex int) model.DeviceHistoryMetricPoint {
	point := model.DeviceHistoryMetricPoint{
		SampledAt: row.SampledAt,
		Online:    row.Online,
	}

	raw := strings.TrimSpace(row.RawJSON)
	if raw == "" {
		return point
	}

	parsed := historyTimelineRawDevice{}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return point
	}

	mon := parsed.Packet.FemtoInputReg.Mon
	if v, ok := readMetricChannel(mon.Temp, tempIndex); ok {
		temp := roundFloat(float64(v)/10.0, 2)
		point.Temp = &temp
	}

	if v, ok := readMetricChannel(mon.Vol, voltageIndex); ok {
		voltage := roundFloat(float64(v), 2)
		point.Voltage = &voltage
	}

	if currentIndex >= 0 && currentIndex < len(mon.PumpMon) {
		current := roundFloat(float64(mon.PumpMon[currentIndex].ActualCur), 2)
		point.Current = &current
	}

	return point
}

func averageAndMaxPositiveU16(values []uint16) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}

	sum := 0.0
	maxVal := 0.0
	count := 0.0
	for _, raw := range values {
		if raw == 0 {
			continue
		}
		v := float64(raw)
		sum += v
		if v > maxVal {
			maxVal = v
		}
		count++
	}
	if count == 0 {
		return 0, 0
	}
	return sum / count, maxVal
}

func downsampleTimelinePoints(points []model.DeviceHistoryTimelinePoint, maxPoints int) []model.DeviceHistoryTimelinePoint {
	if maxPoints <= 0 || len(points) <= maxPoints {
		return points
	}

	step := float64(len(points)) / float64(maxPoints)
	out := make([]model.DeviceHistoryTimelinePoint, 0, maxPoints)
	for i := 0; i < maxPoints; i++ {
		idx := int(float64(i) * step)
		if idx >= len(points) {
			idx = len(points) - 1
		}
		out = append(out, points[idx])
	}
	out[len(out)-1] = points[len(points)-1]
	return out
}

func downsampleMetricPoints(points []model.DeviceHistoryMetricPoint, maxPoints int) []model.DeviceHistoryMetricPoint {
	if maxPoints <= 0 || len(points) <= maxPoints {
		return points
	}

	step := float64(len(points)) / float64(maxPoints)
	out := make([]model.DeviceHistoryMetricPoint, 0, maxPoints)
	for i := 0; i < maxPoints; i++ {
		idx := int(float64(i) * step)
		if idx >= len(points) {
			idx = len(points) - 1
		}
		out = append(out, points[idx])
	}
	out[len(out)-1] = points[len(points)-1]
	return out
}

func readMetricChannel(values []uint16, index int) (uint16, bool) {
	if index < 0 || index >= len(values) {
		return 0, false
	}
	return values[index], true
}

func normalizeMetricIndex(index int) int {
	if index < 0 {
		return 0
	}
	if index > 255 {
		return 255
	}
	return index
}

func roundFloat(value float64, decimals int) float64 {
	if decimals < 0 {
		decimals = 0
	}
	base := math.Pow(10, float64(decimals))
	return math.Round(value*base) / base
}

func alignToInterval(t time.Time, interval time.Duration) time.Time {
	if interval <= 0 {
		interval = time.Minute
	}
	return t.Truncate(interval)
}

func (m *HistoryManager) restoreLatestDevices(ctx context.Context, deviceManager *Manager) (int, error) {
	if db.DB == nil || deviceManager == nil {
		return 0, nil
	}

	table := model.DeviceHistory{}.TableName()
	subQuery := db.DB.WithContext(ctx).
		Model(&model.DeviceHistory{}).
		Select("device_sn, MAX(sampled_at) AS sampled_at").
		Group("device_sn")

	latestRows := make([]model.DeviceHistory, 0, 128)
	if err := db.DB.WithContext(ctx).
		Table(table+" AS h").
		Select("h.*").
		Joins("JOIN (?) latest ON h.device_sn = latest.device_sn AND h.sampled_at = latest.sampled_at", subQuery).
		Order("h.device_sn ASC").
		Find(&latestRows).Error; err != nil {
		return 0, err
	}

	recovered := 0
	for _, row := range latestRows {
		device, ok := historyRowToDevice(row)
		if !ok {
			continue
		}

		copyDevice := device
		deviceManager.Set(copyDevice.Sn, &copyDevice)
		recovered++
	}
	return recovered, nil
}

func historyRowToDevice(row model.DeviceHistory) (model.Device, bool) {
	sn := strings.TrimSpace(row.DeviceSn)
	if sn == "" {
		return model.Device{}, false
	}

	device := model.Device{
		Sn:           sn,
		Name:         strings.TrimSpace(row.DeviceName),
		Last_rx_time: row.SampledAt,
	}

	raw := strings.TrimSpace(row.RawJSON)
	if raw == "" {
		return device, true
	}

	if err := json.Unmarshal([]byte(raw), &device); err != nil {
		return device, true
	}

	device.Sn = strings.TrimSpace(device.Sn)
	if device.Sn == "" {
		device.Sn = sn
	}
	if strings.TrimSpace(device.Name) == "" {
		device.Name = strings.TrimSpace(row.DeviceName)
	}
	if device.Last_rx_time.IsZero() {
		device.Last_rx_time = row.SampledAt
	}
	return device, true
}
