package manager

import (
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	"qa_test_server/db"
	"qa_test_server/model"
)

var BenchmarkManagerGlobal = &BenchmarkManager{}

type BenchmarkManager struct{}

type BenchmarkTopDevice struct {
	DeviceSn         string  `json:"deviceSn"`
	DeviceName       string  `json:"deviceName"`
	Samples          int     `json:"samples"`
	OnlineRatio      float64 `json:"onlineRatio"`
	Completeness     float64 `json:"completeness"`
	RestartCount     int     `json:"restartCount"`
	ParamChangeCount int64   `json:"paramChangeCount"`
	HardwareBate     uint16  `json:"hardwareBate"`
	RiskScore        int     `json:"riskScore"`
	RiskLevel        string  `json:"riskLevel"`
}

type BenchmarkFirmwareDistribution struct {
	HardwareBate uint16 `json:"hardwareBate"`
	Count        int    `json:"count"`
}

type BenchmarkOfflineHour struct {
	Hour         int     `json:"hour"`
	TotalCount   int     `json:"totalCount"`
	OfflineCount int     `json:"offlineCount"`
	OfflineRate  float64 `json:"offlineRate"`
}

type BenchmarkRiskMatrixPoint struct {
	DeviceSn         string  `json:"deviceSn"`
	DeviceName       string  `json:"deviceName"`
	OnlineRatioPct   float64 `json:"onlineRatioPct"`
	CompletenessPct  float64 `json:"completenessPct"`
	ParamChangeCount int64   `json:"paramChangeCount"`
	RiskScore        int     `json:"riskScore"`
	RiskLevel        string  `json:"riskLevel"`
}

type BenchmarkAuditSummary struct {
	Operation string `json:"operation"`
	Success   int64  `json:"success"`
	Failed    int64  `json:"failed"`
}

type BenchmarkMaintenanceCandidate struct {
	DeviceSn       string  `json:"deviceSn"`
	DeviceName     string  `json:"deviceName"`
	OnlineRatio    float64 `json:"onlineRatio"`
	Completeness   float64 `json:"completeness"`
	RestartCount   int     `json:"restartCount"`
	ParamChangeCnt int64   `json:"paramChangeCnt"`
	RiskScore      int     `json:"riskScore"`
	Priority       string  `json:"priority"`
	Reason         string  `json:"reason"`
}

type BenchmarkFeatureItem struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type BenchmarkInsight struct {
	WindowHours int64                           `json:"windowHours"`
	GeneratedAt time.Time                       `json:"generatedAt"`
	SlaTop      []BenchmarkTopDevice            `json:"slaTop"`
	CompleteTop []BenchmarkTopDevice            `json:"completeTop"`
	RestartTop  []BenchmarkTopDevice            `json:"restartTop"`
	ChangeTop   []BenchmarkTopDevice            `json:"changeTop"`
	Firmware    []BenchmarkFirmwareDistribution `json:"firmware"`
	OfflineHeat []BenchmarkOfflineHour          `json:"offlineHeat"`
	RiskMatrix  []BenchmarkRiskMatrixPoint      `json:"riskMatrix"`
	Audit       []BenchmarkAuditSummary         `json:"audit"`
	Candidates  []BenchmarkMaintenanceCandidate `json:"candidates"`
	Features    []BenchmarkFeatureItem          `json:"features"`
}

type benchmarkHistoryRow struct {
	DeviceSn     string
	DeviceName   string
	SampledAt    time.Time
	Online       bool
	HardwareBate uint16
	UptimeSecond uint32
}

type benchmarkAgg struct {
	DeviceSn       string
	DeviceName     string
	Samples        int
	Online         int
	RestartCount   int
	LastUptime     uint32
	LastHardware   uint16
	ParamChangeCnt int64
	OnlineRatio    float64
	Completeness   float64
	RiskScore      int
	RiskLevel      string
}

func (m *BenchmarkManager) Build(window time.Duration) (BenchmarkInsight, error) {
	if db.DB == nil {
		return BenchmarkInsight{}, errors.New("database is unavailable")
	}
	start, end, window := normalizePHMWindow(window)
	windowHours := int64(window.Hours())
	expectedSamples := int(window / time.Minute)
	if expectedSamples <= 0 {
		expectedSamples = 1
	}

	rows := make([]benchmarkHistoryRow, 0, 4096)
	if err := db.DB.Model(&model.DeviceHistory{}).
		Select("device_sn, device_name, sampled_at, online, hardware_bate, uptime_seconds").
		Where("sampled_at >= ? AND sampled_at <= ?", start, end).
		Order("device_sn ASC, sampled_at ASC").
		Limit(300000).
		Find(&rows).Error; err != nil {
		return BenchmarkInsight{}, err
	}

	changeMap, err := m.queryParamChanges(start, end)
	if err != nil {
		return BenchmarkInsight{}, err
	}

	aggMap := make(map[string]*benchmarkAgg, 256)
	offlineTotal := make([]int, 24)
	offlineCount := make([]int, 24)

	for _, row := range rows {
		sn := strings.TrimSpace(row.DeviceSn)
		if sn == "" {
			continue
		}
		item, ok := aggMap[sn]
		if !ok {
			item = &benchmarkAgg{
				DeviceSn:   sn,
				DeviceName: strings.TrimSpace(row.DeviceName),
				LastUptime: row.UptimeSecond,
			}
			aggMap[sn] = item
		}
		if item.DeviceName == "" && strings.TrimSpace(row.DeviceName) != "" {
			item.DeviceName = strings.TrimSpace(row.DeviceName)
		}

		item.Samples++
		if row.Online {
			item.Online++
		}
		if row.UptimeSecond+30 < item.LastUptime {
			item.RestartCount++
		}
		item.LastUptime = row.UptimeSecond
		item.LastHardware = row.HardwareBate

		hour := row.SampledAt.Hour()
		if hour >= 0 && hour < 24 {
			offlineTotal[hour]++
			if !row.Online {
				offlineCount[hour]++
			}
		}
	}

	for sn, agg := range aggMap {
		agg.ParamChangeCnt = changeMap[sn]
		if agg.Samples > 0 {
			agg.OnlineRatio = round2(float64(agg.Online) / float64(agg.Samples))
			agg.Completeness = round2(math.Min(1, float64(agg.Samples)/float64(expectedSamples)))
		}
		agg.RiskScore = computeBenchmarkRisk(agg)
		agg.RiskLevel = scoreToRisk(100 - agg.RiskScore)
	}

	records := make([]BenchmarkTopDevice, 0, len(aggMap))
	for _, agg := range aggMap {
		name := agg.DeviceName
		if name == "" {
			name = agg.DeviceSn
		}
		records = append(records, BenchmarkTopDevice{
			DeviceSn:         agg.DeviceSn,
			DeviceName:       name,
			Samples:          agg.Samples,
			OnlineRatio:      agg.OnlineRatio,
			Completeness:     agg.Completeness,
			RestartCount:     agg.RestartCount,
			ParamChangeCount: agg.ParamChangeCnt,
			HardwareBate:     agg.LastHardware,
			RiskScore:        agg.RiskScore,
			RiskLevel:        agg.RiskLevel,
		})
	}

	slaTop := append([]BenchmarkTopDevice{}, records...)
	sort.Slice(slaTop, func(i, j int) bool {
		if slaTop[i].OnlineRatio == slaTop[j].OnlineRatio {
			return slaTop[i].Samples > slaTop[j].Samples
		}
		return slaTop[i].OnlineRatio > slaTop[j].OnlineRatio
	})
	if len(slaTop) > 15 {
		slaTop = slaTop[:15]
	}

	completeTop := append([]BenchmarkTopDevice{}, records...)
	sort.Slice(completeTop, func(i, j int) bool {
		if completeTop[i].Completeness == completeTop[j].Completeness {
			return completeTop[i].Samples > completeTop[j].Samples
		}
		return completeTop[i].Completeness > completeTop[j].Completeness
	})
	if len(completeTop) > 15 {
		completeTop = completeTop[:15]
	}

	restartTop := append([]BenchmarkTopDevice{}, records...)
	sort.Slice(restartTop, func(i, j int) bool {
		if restartTop[i].RestartCount == restartTop[j].RestartCount {
			return restartTop[i].RiskScore > restartTop[j].RiskScore
		}
		return restartTop[i].RestartCount > restartTop[j].RestartCount
	})
	if len(restartTop) > 15 {
		restartTop = restartTop[:15]
	}

	changeTop := append([]BenchmarkTopDevice{}, records...)
	sort.Slice(changeTop, func(i, j int) bool {
		if changeTop[i].ParamChangeCount == changeTop[j].ParamChangeCount {
			return changeTop[i].RiskScore > changeTop[j].RiskScore
		}
		return changeTop[i].ParamChangeCount > changeTop[j].ParamChangeCount
	})
	if len(changeTop) > 15 {
		changeTop = changeTop[:15]
	}

	firmware := buildFirmwareDistribution(records)
	offlineHeat := buildOfflineHeat(offlineTotal, offlineCount)
	riskMatrix := buildRiskMatrix(records)
	if len(riskMatrix) > 30 {
		riskMatrix = riskMatrix[:30]
	}
	audit, err := m.queryAuditSummary(start)
	if err != nil {
		return BenchmarkInsight{}, err
	}
	candidates := buildMaintenanceCandidates(records)
	if len(candidates) > 20 {
		candidates = candidates[:20]
	}

	return BenchmarkInsight{
		WindowHours: windowHours,
		GeneratedAt: time.Now(),
		SlaTop:      slaTop,
		CompleteTop: completeTop,
		RestartTop:  restartTop,
		ChangeTop:   changeTop,
		Firmware:    firmware,
		OfflineHeat: offlineHeat,
		RiskMatrix:  riskMatrix,
		Audit:       audit,
		Candidates:  candidates,
		Features:    benchmarkFeatures(),
	}, nil
}

func (m *BenchmarkManager) queryParamChanges(start, end time.Time) (map[string]int64, error) {
	type row struct {
		DeviceSn string
		Cnt      int64
	}
	changes := make([]row, 0, 256)
	if err := db.DB.Model(&model.DeviceParamChange{}).
		Select("device_sn, COUNT(*) AS cnt").
		Where("changed_at >= ? AND changed_at <= ?", start, end).
		Group("device_sn").
		Find(&changes).Error; err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(changes))
	for _, item := range changes {
		sn := strings.TrimSpace(item.DeviceSn)
		if sn == "" {
			continue
		}
		out[sn] = item.Cnt
	}
	return out, nil
}

func (m *BenchmarkManager) queryAuditSummary(start time.Time) ([]BenchmarkAuditSummary, error) {
	type row struct {
		Operation string
		Success   bool
		Cnt       int64
	}
	rows := make([]row, 0, 16)
	if err := db.DB.Model(&model.DecryptLog{}).
		Select("operation, success, COUNT(*) AS cnt").
		Where("created_at >= ?", start).
		Group("operation, success").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	byOp := map[string]*BenchmarkAuditSummary{}
	for _, item := range rows {
		op := strings.TrimSpace(strings.ToLower(item.Operation))
		if op == "" {
			op = "decode"
		}
		if _, ok := byOp[op]; !ok {
			byOp[op] = &BenchmarkAuditSummary{
				Operation: op,
			}
		}
		if item.Success {
			byOp[op].Success += item.Cnt
		} else {
			byOp[op].Failed += item.Cnt
		}
	}

	out := make([]BenchmarkAuditSummary, 0, len(byOp))
	for _, item := range byOp {
		out = append(out, *item)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Operation < out[j].Operation
	})
	return out, nil
}

func buildFirmwareDistribution(records []BenchmarkTopDevice) []BenchmarkFirmwareDistribution {
	counter := map[uint16]int{}
	for _, item := range records {
		counter[item.HardwareBate]++
	}
	out := make([]BenchmarkFirmwareDistribution, 0, len(counter))
	for bate, cnt := range counter {
		out = append(out, BenchmarkFirmwareDistribution{
			HardwareBate: bate,
			Count:        cnt,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count == out[j].Count {
			return out[i].HardwareBate < out[j].HardwareBate
		}
		return out[i].Count > out[j].Count
	})
	return out
}

func buildOfflineHeat(total, offline []int) []BenchmarkOfflineHour {
	out := make([]BenchmarkOfflineHour, 0, 24)
	for i := 0; i < 24; i++ {
		rate := 0.0
		if total[i] > 0 {
			rate = round2(float64(offline[i]) / float64(total[i]))
		}
		out = append(out, BenchmarkOfflineHour{
			Hour:         i,
			TotalCount:   total[i],
			OfflineCount: offline[i],
			OfflineRate:  rate,
		})
	}
	return out
}

func buildRiskMatrix(records []BenchmarkTopDevice) []BenchmarkRiskMatrixPoint {
	out := make([]BenchmarkRiskMatrixPoint, 0, len(records))
	for _, item := range records {
		out = append(out, BenchmarkRiskMatrixPoint{
			DeviceSn:         item.DeviceSn,
			DeviceName:       item.DeviceName,
			OnlineRatioPct:   round2(item.OnlineRatio * 100),
			CompletenessPct:  round2(item.Completeness * 100),
			ParamChangeCount: item.ParamChangeCount,
			RiskScore:        item.RiskScore,
			RiskLevel:        item.RiskLevel,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].RiskScore > out[j].RiskScore
	})
	return out
}

func buildMaintenanceCandidates(records []BenchmarkTopDevice) []BenchmarkMaintenanceCandidate {
	out := make([]BenchmarkMaintenanceCandidate, 0, len(records))
	for _, item := range records {
		if item.RiskScore < 20 {
			continue
		}
		priority := "low"
		if item.RiskScore >= 60 {
			priority = "high"
		} else if item.RiskScore >= 35 {
			priority = "medium"
		}
		reason := buildCandidateReason(item)
		out = append(out, BenchmarkMaintenanceCandidate{
			DeviceSn:       item.DeviceSn,
			DeviceName:     item.DeviceName,
			OnlineRatio:    item.OnlineRatio,
			Completeness:   item.Completeness,
			RestartCount:   item.RestartCount,
			ParamChangeCnt: item.ParamChangeCount,
			RiskScore:      item.RiskScore,
			Priority:       priority,
			Reason:         reason,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].RiskScore > out[j].RiskScore
	})
	return out
}

func buildCandidateReason(item BenchmarkTopDevice) string {
	reasons := make([]string, 0, 4)
	if item.OnlineRatio < 0.95 {
		reasons = append(reasons, "在线率偏低")
	}
	if item.Completeness < 0.85 {
		reasons = append(reasons, "数据完整率不足")
	}
	if item.RestartCount > 0 {
		reasons = append(reasons, "存在重启波动")
	}
	if item.ParamChangeCount > 20 {
		reasons = append(reasons, "参数变更频繁")
	}
	if len(reasons) == 0 {
		return "常规巡检"
	}
	return strings.Join(reasons, "、")
}

func computeBenchmarkRisk(item *benchmarkAgg) int {
	if item == nil {
		return 0
	}
	risk := 0
	if item.OnlineRatio < 0.95 {
		risk += int((0.95 - item.OnlineRatio) * 180)
	}
	if item.Completeness < 0.85 {
		risk += int((0.85 - item.Completeness) * 150)
	}
	risk += item.RestartCount * 8
	if item.ParamChangeCnt > 0 {
		changePenalty := int(item.ParamChangeCnt / 3)
		if changePenalty > 25 {
			changePenalty = 25
		}
		risk += changePenalty
	}
	if risk < 0 {
		return 0
	}
	if risk > 100 {
		return 100
	}
	return risk
}

func benchmarkFeatures() []BenchmarkFeatureItem {
	return []BenchmarkFeatureItem{
		{Code: "F01", Name: "SLA在线率排行", Description: "按设备在线率进行服务水平排名", Status: "已适配"},
		{Code: "F02", Name: "数据完整率监控", Description: "按分钟采样窗口评估历史数据完整度", Status: "已适配"},
		{Code: "F03", Name: "重启波动识别", Description: "基于运行时长回退识别异常重启", Status: "已适配"},
		{Code: "F04", Name: "参数变更热榜", Description: "统计窗口内参数变更频次Top设备", Status: "已适配"},
		{Code: "F05", Name: "固件版本分布", Description: "统计硬件版本分布与占比", Status: "已适配"},
		{Code: "F06", Name: "离线时段热力", Description: "按小时聚合离线发生频率", Status: "已适配"},
		{Code: "F07", Name: "风险矩阵", Description: "在线率与完整率双维风险评分", Status: "已适配"},
		{Code: "F08", Name: "运维审计汇总", Description: "统计密钥工具操作成功/失败", Status: "已适配"},
		{Code: "F09", Name: "维护候选清单", Description: "自动筛选高风险设备形成候选列表", Status: "已适配"},
		{Code: "F10", Name: "运营报告导出", Description: "导出窗口期设备健康报告CSV", Status: "已适配"},
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
