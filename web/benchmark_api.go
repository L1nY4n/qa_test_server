package web

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"qa_test_server/manager"
)

func benchmarkInsights(c *gin.Context) {
	windowHours := queryInt(c, "windowHours", 24)
	insight, err := manager.BenchmarkManagerGlobal.Build(time.Duration(windowHours) * time.Hour)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, insight)
}

func exportBenchmarkReport(c *gin.Context) {
	windowHours := queryInt(c, "windowHours", 24)
	insight, err := manager.BenchmarkManagerGlobal.Build(time.Duration(windowHours) * time.Hour)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	_, _ = buf.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(buf)
	_ = writer.Write([]string{"section", "rank", "deviceSn", "deviceName", "metricA", "metricB", "metricC", "note"})
	writeFeatureSection(writer, insight.Features)
	writeTopSection(writer, "slaTop", insight.SlaTop)
	writeTopSection(writer, "completeTop", insight.CompleteTop)
	writeTopSection(writer, "restartTop", insight.RestartTop)
	writeTopSection(writer, "changeTop", insight.ChangeTop)
	writeFirmwareSection(writer, insight.Firmware)
	writeOfflineSection(writer, insight.OfflineHeat)
	writeRiskMatrixSection(writer, insight.RiskMatrix)
	writeAuditSection(writer, insight.Audit)
	writeCandidateSection(writer, insight.Candidates)
	writer.Flush()

	fileName := fmt.Sprintf("benchmark_report_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	c.Data(200, "text/csv; charset=utf-8", buf.Bytes())
}

func writeFeatureSection(writer *csv.Writer, items []manager.BenchmarkFeatureItem) {
	for idx, item := range items {
		_ = writer.Write([]string{
			"featureChecklist",
			fmt.Sprintf("%d", idx+1),
			"",
			item.Name,
			item.Code,
			item.Status,
			"",
			item.Description,
		})
	}
}

func writeTopSection(writer *csv.Writer, section string, items []manager.BenchmarkTopDevice) {
	for idx, item := range items {
		_ = writer.Write([]string{
			section,
			fmt.Sprintf("%d", idx+1),
			item.DeviceSn,
			item.DeviceName,
			fmt.Sprintf("online=%.2f%%", item.OnlineRatio*100),
			fmt.Sprintf("complete=%.2f%%", item.Completeness*100),
			fmt.Sprintf("restart=%d change=%d risk=%d", item.RestartCount, item.ParamChangeCount, item.RiskScore),
			item.RiskLevel,
		})
	}
}

func writeFirmwareSection(writer *csv.Writer, items []manager.BenchmarkFirmwareDistribution) {
	for idx, item := range items {
		_ = writer.Write([]string{
			"firmwareDistribution",
			fmt.Sprintf("%d", idx+1),
			"",
			"",
			fmt.Sprintf("hardwareBate=%d", item.HardwareBate),
			fmt.Sprintf("count=%d", item.Count),
			"",
			"",
		})
	}
}

func writeOfflineSection(writer *csv.Writer, items []manager.BenchmarkOfflineHour) {
	for _, item := range items {
		_ = writer.Write([]string{
			"offlineHeat",
			fmt.Sprintf("%d", item.Hour),
			"",
			"",
			fmt.Sprintf("offlineRate=%.2f%%", item.OfflineRate*100),
			fmt.Sprintf("offlineCount=%d", item.OfflineCount),
			fmt.Sprintf("total=%d", item.TotalCount),
			"",
		})
	}
}

func writeRiskMatrixSection(writer *csv.Writer, items []manager.BenchmarkRiskMatrixPoint) {
	for idx, item := range items {
		_ = writer.Write([]string{
			"riskMatrix",
			fmt.Sprintf("%d", idx+1),
			item.DeviceSn,
			item.DeviceName,
			fmt.Sprintf("online=%.2f%%", item.OnlineRatioPct),
			fmt.Sprintf("complete=%.2f%%", item.CompletenessPct),
			fmt.Sprintf("change=%d score=%d", item.ParamChangeCount, item.RiskScore),
			item.RiskLevel,
		})
	}
}

func writeAuditSection(writer *csv.Writer, items []manager.BenchmarkAuditSummary) {
	for idx, item := range items {
		_ = writer.Write([]string{
			"auditSummary",
			fmt.Sprintf("%d", idx+1),
			"",
			item.Operation,
			fmt.Sprintf("success=%d", item.Success),
			fmt.Sprintf("failed=%d", item.Failed),
			"",
			"",
		})
	}
}

func writeCandidateSection(writer *csv.Writer, items []manager.BenchmarkMaintenanceCandidate) {
	for idx, item := range items {
		_ = writer.Write([]string{
			"maintenanceCandidates",
			fmt.Sprintf("%d", idx+1),
			item.DeviceSn,
			item.DeviceName,
			fmt.Sprintf("risk=%d", item.RiskScore),
			fmt.Sprintf("online=%.2f%% complete=%.2f%%", item.OnlineRatio*100, item.Completeness*100),
			fmt.Sprintf("restart=%d change=%d", item.RestartCount, item.ParamChangeCnt),
			fmt.Sprintf("%s: %s", item.Priority, item.Reason),
		})
	}
}
