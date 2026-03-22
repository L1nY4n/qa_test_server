package manager

import (
	"context"
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	"qa_test_server/db"
	"qa_test_server/model"
)

var PHMManagerGlobal = &PHMManager{}

const (
	PHMRiskLow      = "low"
	PHMRiskMedium   = "medium"
	PHMRiskHigh     = "high"
	PHMRiskCritical = "critical"
)

type PHMManager struct{}

type DevicePHMSummary struct {
	DeviceSn         string    `json:"deviceSn"`
	DeviceName       string    `json:"deviceName"`
	WindowHours      int64     `json:"windowHours"`
	Samples          int       `json:"samples"`
	OnlineRatio      float64   `json:"onlineRatio"`
	LastSampledAt    time.Time `json:"lastSampledAt"`
	FreshnessMinutes int64     `json:"freshnessMinutes"`
	RebootCount      int       `json:"rebootCount"`
	HardwareShift    int       `json:"hardwareShift"`
	DataGapCount     int       `json:"dataGapCount"`
	ParamChangeCount int64     `json:"paramChangeCount"`
	HealthScore      int       `json:"healthScore"`
	RiskLevel        string    `json:"riskLevel"`
	Reasons          []string  `json:"reasons"`
	Recommendation   string    `json:"recommendation"`
}

type DevicePHMDetail struct {
	Summary       DevicePHMSummary           `json:"summary"`
	HistoryPoints []model.DeviceHistoryPoint `json:"historyPoints"`
	RecentChanges []model.DeviceParamChange  `json:"recentChanges"`
	WindowStart   time.Time                  `json:"windowStart"`
	WindowEnd     time.Time                  `json:"windowEnd"`
}

type PHMDeleteResult struct {
	Requested          int       `json:"requested"`
	Cleared            []string  `json:"cleared"`
	Missing            []string  `json:"missing"`
	CacheRemoved       int       `json:"cacheRemoved"`
	HistoryRowsDeleted int64     `json:"historyRowsDeleted"`
	ChangeRowsDeleted  int64     `json:"changeRowsDeleted"`
	DeletedAt          time.Time `json:"deletedAt"`
}

func (m *PHMManager) Overview(window time.Duration, limit int) ([]DevicePHMSummary, error) {
	if db.DB == nil {
		return nil, errors.New("database is unavailable")
	}
	start, end, window := normalizePHMWindow(window)
	sns, err := m.listTrackedDeviceSNs()
	if err != nil {
		return nil, err
	}
	if len(sns) == 0 {
		return []DevicePHMSummary{}, nil
	}

	items := make([]DevicePHMSummary, 0, len(sns))
	for _, sn := range sns {
		summary, err := m.evaluateDevice(sn, start, end, window)
		if err != nil {
			continue
		}
		items = append(items, summary)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].HealthScore == items[j].HealthScore {
			return strings.Compare(items[i].DeviceSn, items[j].DeviceSn) < 0
		}
		return items[i].HealthScore < items[j].HealthScore
	})
	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (m *PHMManager) DeviceDetail(sn string, window time.Duration, maxChanges int, maxHistory int) (DevicePHMDetail, error) {
	if db.DB == nil {
		return DevicePHMDetail{}, errors.New("database is unavailable")
	}
	sn = strings.TrimSpace(sn)
	if sn == "" {
		return DevicePHMDetail{}, errors.New("device sn is required")
	}
	start, end, window := normalizePHMWindow(window)
	summary, err := m.evaluateDevice(sn, start, end, window)
	if err != nil {
		return DevicePHMDetail{}, err
	}

	if maxHistory <= 0 {
		maxHistory = 600
	}
	if maxHistory > 2000 {
		maxHistory = 2000
	}
	if maxChanges <= 0 {
		maxChanges = 200
	}
	if maxChanges > 2000 {
		maxChanges = 2000
	}

	historyRows := make([]model.DeviceHistory, 0, maxHistory)
	if err := db.DB.Model(&model.DeviceHistory{}).
		Where("device_sn = ? AND sampled_at >= ? AND sampled_at <= ?", sn, start, end).
		Order("sampled_at ASC").
		Limit(maxHistory).
		Find(&historyRows).Error; err != nil {
		return DevicePHMDetail{}, err
	}
	historyPoints := make([]model.DeviceHistoryPoint, 0, len(historyRows))
	for _, row := range historyRows {
		historyPoints = append(historyPoints, model.DeviceHistoryPoint{
			SampledAt:     row.SampledAt,
			Online:        row.Online,
			HardwareBate:  row.HardwareBate,
			UptimeSeconds: row.UptimeSeconds,
			PumpCount:     row.PumpCount,
			TempCount:     row.TempCount,
		})
	}

	changeRows := make([]model.DeviceParamChange, 0, maxChanges)
	if err := db.DB.Model(&model.DeviceParamChange{}).
		Where("device_sn = ? AND changed_at >= ? AND changed_at <= ?", sn, start, end).
		Order("changed_at DESC").
		Limit(maxChanges).
		Find(&changeRows).Error; err != nil {
		return DevicePHMDetail{}, err
	}

	return DevicePHMDetail{
		Summary:       summary,
		HistoryPoints: historyPoints,
		RecentChanges: changeRows,
		WindowStart:   start,
		WindowEnd:     end,
	}, nil
}

func (m *PHMManager) DeleteDevices(ctx context.Context, sns []string) (PHMDeleteResult, error) {
	if db.DB == nil {
		return PHMDeleteResult{}, errors.New("database is unavailable")
	}
	clean := normalizeSNList(sns)
	if len(clean) == 0 {
		return PHMDeleteResult{}, errors.New("device sn list is required")
	}

	historyRows := int64(0)
	if res := db.DB.WithContext(ctx).Where("device_sn IN ?", clean).Delete(&model.DeviceHistory{}); res.Error != nil {
		return PHMDeleteResult{}, res.Error
	} else {
		historyRows = res.RowsAffected
	}

	changeRows := int64(0)
	if res := db.DB.WithContext(ctx).Where("device_sn IN ?", clean).Delete(&model.DeviceParamChange{}); res.Error != nil {
		return PHMDeleteResult{}, res.Error
	} else {
		changeRows = res.RowsAffected
	}

	cleared := make([]string, 0, len(clean))
	missing := make([]string, 0, len(clean))
	cacheRemoved := 0
	for _, sn := range clean {
		_, exists := ManagerGlabal.Get(sn)
		ManagerGlabal.Delete(sn)
		if exists {
			cacheRemoved++
			cleared = append(cleared, sn)
		} else {
			missing = append(missing, sn)
		}
	}

	return PHMDeleteResult{
		Requested:          len(clean),
		Cleared:            cleared,
		Missing:            missing,
		CacheRemoved:       cacheRemoved,
		HistoryRowsDeleted: historyRows,
		ChangeRowsDeleted:  changeRows,
		DeletedAt:          time.Now(),
	}, nil
}

func (m *PHMManager) evaluateDevice(sn string, start, end time.Time, window time.Duration) (DevicePHMSummary, error) {
	rows := make([]model.DeviceHistory, 0, 720)
	if err := db.DB.Model(&model.DeviceHistory{}).
		Where("device_sn = ? AND sampled_at >= ? AND sampled_at <= ?", sn, start, end).
		Order("sampled_at ASC").
		Limit(5000).
		Find(&rows).Error; err != nil {
		return DevicePHMSummary{}, err
	}

	var changeCount int64
	if err := db.DB.Model(&model.DeviceParamChange{}).
		Where("device_sn = ? AND changed_at >= ? AND changed_at <= ?", sn, start, end).
		Count(&changeCount).Error; err != nil {
		return DevicePHMSummary{}, err
	}

	summary := DevicePHMSummary{
		DeviceSn:         sn,
		WindowHours:      int64(window.Hours()),
		Samples:          len(rows),
		ParamChangeCount: changeCount,
		Reasons:          []string{},
	}

	if len(rows) > 0 {
		summary.DeviceName = strings.TrimSpace(rows[len(rows)-1].DeviceName)
		summary.LastSampledAt = rows[len(rows)-1].SampledAt
		summary.FreshnessMinutes = int64(math.Max(0, time.Since(summary.LastSampledAt).Minutes()))
	} else {
		if dev, ok := ManagerGlabal.Get(sn); ok && dev != nil {
			summary.DeviceName = dev.Name
			summary.LastSampledAt = dev.Last_rx_time
			if !dev.Last_rx_time.IsZero() {
				summary.FreshnessMinutes = int64(math.Max(0, time.Since(dev.Last_rx_time).Minutes()))
			}
		}
	}
	if summary.DeviceName == "" {
		summary.DeviceName = sn
	}

	onlineCount := 0
	rebootCount := 0
	hardwareShift := 0
	dataGapCount := 0
	for i := range rows {
		if rows[i].Online {
			onlineCount++
		}
		if i > 0 {
			prev := rows[i-1]
			curr := rows[i]
			if curr.UptimeSeconds+30 < prev.UptimeSeconds {
				rebootCount++
			}
			if curr.HardwareBate != prev.HardwareBate {
				hardwareShift++
			}
			if curr.SampledAt.Sub(prev.SampledAt) > 2*time.Minute+30*time.Second {
				dataGapCount++
			}
		}
	}
	summary.RebootCount = rebootCount
	summary.HardwareShift = hardwareShift
	summary.DataGapCount = dataGapCount
	if summary.Samples > 0 {
		summary.OnlineRatio = float64(onlineCount) / float64(summary.Samples)
	} else {
		summary.OnlineRatio = 0
	}

	summary.HealthScore = computeHealthScore(&summary)
	summary.RiskLevel = scoreToRisk(summary.HealthScore)
	summary.Recommendation = buildRecommendation(summary)

	return summary, nil
}

func computeHealthScore(s *DevicePHMSummary) int {
	score := 100
	if s.Samples == 0 {
		s.Reasons = append(s.Reasons, "no samples in selected window")
		score -= 85
		return clampScore(score)
	}

	if s.OnlineRatio < 0.98 {
		penalty := int((0.98 - s.OnlineRatio) * 120)
		if penalty > 40 {
			penalty = 40
		}
		score -= penalty
		s.Reasons = append(s.Reasons, "low online ratio")
	}

	if s.FreshnessMinutes > 5 {
		penalty := int(s.FreshnessMinutes - 5)
		if penalty > 25 {
			penalty = 25
		}
		score -= penalty
		s.Reasons = append(s.Reasons, "stale telemetry")
	}

	if s.RebootCount > 0 {
		penalty := s.RebootCount * 5
		if penalty > 20 {
			penalty = 20
		}
		score -= penalty
		s.Reasons = append(s.Reasons, "detected reboot or power instability")
	}

	if s.HardwareShift > 0 {
		penalty := s.HardwareShift * 4
		if penalty > 16 {
			penalty = 16
		}
		score -= penalty
		s.Reasons = append(s.Reasons, "hardware version changed")
	}

	if s.DataGapCount > 0 {
		penalty := s.DataGapCount * 3
		if penalty > 15 {
			penalty = 15
		}
		score -= penalty
		s.Reasons = append(s.Reasons, "sampling pipeline gaps")
	}

	if s.ParamChangeCount > 0 {
		penalty := int(s.ParamChangeCount / 4)
		if penalty < 2 {
			penalty = 2
		}
		if penalty > 20 {
			penalty = 20
		}
		score -= penalty
		s.Reasons = append(s.Reasons, "frequent system parameter changes")
	}

	return clampScore(score)
}

func clampScore(score int) int {
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func scoreToRisk(score int) string {
	switch {
	case score >= 85:
		return PHMRiskLow
	case score >= 70:
		return PHMRiskMedium
	case score >= 50:
		return PHMRiskHigh
	default:
		return PHMRiskCritical
	}
}

func buildRecommendation(s DevicePHMSummary) string {
	if s.Samples == 0 {
		return "Check collector service, device link, and network connectivity."
	}
	if s.HealthScore < 50 {
		return "Schedule immediate maintenance and inspect power, communication, and key modules."
	}
	if s.RebootCount >= 2 {
		return "Inspect power stability and thermal conditions due to repeated reboot behavior."
	}
	if s.ParamChangeCount >= 20 {
		return "Freeze configuration and review parameter strategy to reduce frequent write operations."
	}
	if s.DataGapCount > 0 || s.OnlineRatio < 0.95 {
		return "Inspect network quality and telemetry pipeline to reduce data gaps and offline duration."
	}
	return "Device health is stable. Keep routine inspections and trend tracking."
}

func (m *PHMManager) listTrackedDeviceSNs() ([]string, error) {
	seen := make(map[string]struct{})
	sns := make([]string, 0, 64)

	for _, dev := range ManagerGlabal.List() {
		sn := strings.TrimSpace(dev.Sn)
		if sn == "" {
			continue
		}
		if _, ok := seen[sn]; ok {
			continue
		}
		seen[sn] = struct{}{}
		sns = append(sns, sn)
	}

	dbSNS := make([]string, 0, 64)
	retentionStart := time.Now().Add(-30 * 24 * time.Hour)
	if err := db.DB.Model(&model.DeviceHistory{}).
		Where("sampled_at >= ?", retentionStart).
		Distinct("device_sn").
		Pluck("device_sn", &dbSNS).Error; err != nil {
		return nil, err
	}
	for _, raw := range dbSNS {
		sn := strings.TrimSpace(raw)
		if sn == "" {
			continue
		}
		if _, ok := seen[sn]; ok {
			continue
		}
		seen[sn] = struct{}{}
		sns = append(sns, sn)
	}

	sort.Strings(sns)
	return sns, nil
}

func normalizePHMWindow(window time.Duration) (time.Time, time.Time, time.Duration) {
	if window <= 0 {
		window = 24 * time.Hour
	}
	if window > 10*24*time.Hour {
		window = 10 * 24 * time.Hour
	}
	end := time.Now()
	start := end.Add(-window)
	return start, end, window
}

func normalizeSNList(input []string) []string {
	if len(input) == 0 {
		return []string{}
	}
	seen := make(map[string]struct{}, len(input))
	out := make([]string, 0, len(input))
	for _, raw := range input {
		sn := strings.TrimSpace(raw)
		if sn == "" {
			continue
		}
		if _, ok := seen[sn]; ok {
			continue
		}
		seen[sn] = struct{}{}
		out = append(out, sn)
	}
	return out
}
