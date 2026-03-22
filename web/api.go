package web

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"qa_test_server/db"
	"qa_test_server/manager"
	"qa_test_server/model"
)

var appStartedAt = time.Now()

func Api(r *gin.Engine) {
	authRoute := r.Group("/auth")
	{
		authRoute.POST("/register", register)
		authRoute.POST("/login", login)
	}

	apiRoute := r.Group("/")
	apiRoute.Use(authRequired())
	{
		secureAuthRoute := apiRoute.Group("/auth")
		{
			secureAuthRoute.GET("/profile", profile)
		}

		systemRoute := apiRoute.Group("/system")
		{
			systemRoute.GET("/info", sysInfo)
			systemRoute.GET("/health", health)
			systemRoute.POST("/reset", requireRoles(model.RoleAdmin), systemReset)
		}

		deviceRoute := apiRoute.Group("/device")
		{
			deviceRoute.GET("/list", deviceList)
			deviceRoute.GET("/stats", deviceStats)
			deviceRoute.GET("/info/:id", deviceInfo)
			deviceRoute.GET("/history/:id", deviceHistory)
			deviceRoute.GET("/history/:id/timeline", deviceHistoryTimeline)
			deviceRoute.GET("/history/:id/metrics", deviceHistoryMetrics)
			deviceRoute.GET("/changes/:id", deviceParamChanges)
			deviceRoute.GET("/changes/:id/export", exportDeviceParamChanges)
			deviceRoute.GET("/randomData/:current", randomData)
		}

		phmRoute := apiRoute.Group("/phm")
		phmRoute.Use(requireRoles(model.RoleAdmin, model.RoleOperator))
		{
			phmRoute.GET("/overview", phmOverview)
			phmRoute.GET("/device/:id", phmDeviceDetail)
			phmRoute.DELETE("/device/:id", phmDeleteSingle)
			phmRoute.POST("/devices/delete", phmDeleteBatch)
		}

		benchmarkRoute := apiRoute.Group("/benchmark")
		benchmarkRoute.Use(requireRoles(model.RoleAdmin, model.RoleOperator))
		{
			benchmarkRoute.GET("/insights", benchmarkInsights)
			benchmarkRoute.GET("/export", exportBenchmarkReport)
		}

		debugRoute := apiRoute.Group("/debug")
		debugRoute.Use(requireRoles(model.RoleAdmin, model.RoleOperator))
		{
			debugRoute.GET("/virtual/status", virtualStatus)
			debugRoute.POST("/virtual/start", virtualStart)
			debugRoute.POST("/virtual/stop", virtualStop)
			debugRoute.POST("/virtual/pulse", virtualPulse)
			debugRoute.POST("/virtual/stress/pulse", virtualStressPulse)
			debugRoute.POST("/decrypt/time-key", decryptTimeKey)
			debugRoute.POST("/decrypt/time-key/generate", generateTimeKey)
			debugRoute.GET("/decrypt/logs", decryptLogs)
		}

		userRoute := apiRoute.Group("/users")
		userRoute.Use(requireRoles(model.RoleAdmin))
		{
			userRoute.GET("/list", userList)
			userRoute.POST("", createUser)
			userRoute.PUT("/:id", updateUser)
			userRoute.PUT("/:id/password", resetUserPassword)
			userRoute.DELETE("/:id", deleteUser)
		}
	}
}

func sysInfo(c *gin.Context) {
	historyStats := manager.HistoryManagerGlobal.Stats()
	virtualStatus := manager.VirtualDeviceManagerGlobal.Status()

	c.JSON(200, gin.H{
		"ver":      15,
		"author":   "lf",
		"httpAddr": appConfig.HTTPAddr,
		"tcpAddr":  appConfig.TCPAddr,
		"history": gin.H{
			"sampleEverySeconds": int64(appConfig.HistorySampleEvery.Seconds()),
			"retentionDays":      int64(appConfig.HistoryRetention.Hours() / 24),
			"recoveredCount":     historyStats.RecoveredCount,
			"recoveredAt":        historyStats.RecoveredAt,
		},
		"changeLog": gin.H{
			"retentionDays": int64(appConfig.ChangeLogRetention.Hours() / 24),
		},
		"virtualDevice": virtualStatus,
	})
}

func health(c *gin.Context) {
	dbOK, dbMessage := db.Health(2 * time.Second)
	wsStats := WsManager.Stats()

	payload := gin.H{
		"uptimeSeconds": int64(time.Since(appStartedAt).Seconds()),
		"db": gin.H{
			"ok":      dbOK,
			"message": dbMessage,
		},
		"websocket": wsStats,
		"server": gin.H{
			"httpAddr": appConfig.HTTPAddr,
			"tcpAddr":  appConfig.TCPAddr,
		},
	}

	if !dbOK {
		c.JSON(503, apiResponse{
			Success: false,
			Data:    payload,
			Error:   "database is unavailable",
		})
		return
	}
	ok(c, payload)
}

func systemReset(c *gin.Context) {
	if _, err := manager.VirtualDeviceManagerGlobal.Stop(true); err != nil {
		fail(c, http.StatusInternalServerError, "failed to stop virtual devices: "+err.Error())
		return
	}

	clearedDevices := manager.ManagerGlabal.ClearAll()

	historyDeleted, err := manager.HistoryManagerGlobal.ClearAll(c.Request.Context())
	if err != nil {
		fail(c, http.StatusInternalServerError, "failed to clear history data: "+err.Error())
		return
	}

	changeDeleted, err := manager.ParamChangeManagerGlobal.ClearAll(c.Request.Context())
	if err != nil {
		fail(c, http.StatusInternalServerError, "failed to clear parameter change data: "+err.Error())
		return
	}

	ok(c, gin.H{
		"resetAt":               time.Now(),
		"clearedDevices":        clearedDevices,
		"historyRowsDeleted":    historyDeleted,
		"paramRowsDeleted":      changeDeleted,
		"virtualDevicesStopped": true,
	})
}

func deviceList(c *gin.Context) {
	keyword := c.Query("keyword")
	group := c.Query("group")
	onlineOnly := queryBool(c, "onlineOnly", false)
	offset := queryInt(c, "offset", 0)
	limit := queryInt(c, "limit", 50)
	activeWithin := time.Duration(queryInt(c, "activeWithin", int(appConfig.DeviceActiveWithin.Seconds()))) * time.Second
	withMeta := queryBool(c, "meta", false)
	summaryOnly := queryBool(c, "summary", false)

	if summaryOnly {
		items, total := manager.ManagerGlabal.QuerySummary(keyword, group, onlineOnly, activeWithin, offset, limit)
		c.Header("X-Total-Count", strconv.Itoa(total))
		if withMeta {
			ok(c, gin.H{
				"items":               items,
				"total":               total,
				"offset":              offset,
				"limit":               limit,
				"group":               group,
				"onlineOnly":          onlineOnly,
				"activeWithinSeconds": int64(activeWithin.Seconds()),
			})
			return
		}
		ok(c, items)
		return
	}

	items, total := manager.ManagerGlabal.Query(keyword, group, onlineOnly, activeWithin, offset, limit)
	c.Header("X-Total-Count", strconv.Itoa(total))
	if withMeta {
		ok(c, gin.H{
			"items":               items,
			"total":               total,
			"offset":              offset,
			"limit":               limit,
			"group":               group,
			"onlineOnly":          onlineOnly,
			"activeWithinSeconds": int64(activeWithin.Seconds()),
		})
		return
	}
	ok(c, items)
}

func deviceStats(c *gin.Context) {
	activeWithin := time.Duration(queryInt(c, "activeWithin", int(appConfig.DeviceActiveWithin.Seconds()))) * time.Second
	ok(c, manager.ManagerGlabal.Stats(activeWithin))
}

func deviceInfo(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		fail(c, 400, "device id is required")
		return
	}

	device, exists := manager.ManagerGlabal.Get(id)
	if !exists {
		fail(c, 404, "device not found")
		return
	}

	activeWithin := time.Duration(queryInt(c, "activeWithin", int(appConfig.DeviceActiveWithin.Seconds()))) * time.Second
	ok(c, gin.H{
		"device": device,
		"online": !device.Last_rx_time.IsZero() && time.Since(device.Last_rx_time) <= activeWithin,
	})
}

func deviceHistory(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		fail(c, 400, "device id is required")
		return
	}

	start := queryTime(c, "start")
	end := queryTime(c, "end")
	offset := queryInt(c, "offset", 0)
	limit := queryInt(c, "limit", 360)
	withRaw := queryBool(c, "raw", false)

	q := manager.HistoryQuery{
		DeviceSn: id,
		Start:    start,
		End:      end,
		Offset:   offset,
		Limit:    limit,
	}

	if withRaw {
		items, total, err := manager.HistoryManagerGlobal.Query(q)
		if err != nil {
			fail(c, 400, err.Error())
			return
		}
		ok(c, gin.H{
			"items":  items,
			"total":  total,
			"offset": offset,
			"limit":  limit,
			"raw":    true,
		})
		return
	}

	points, total, err := manager.HistoryManagerGlobal.QueryPoints(q)
	if err != nil {
		fail(c, 400, err.Error())
		return
	}
	ok(c, gin.H{
		"items":  points,
		"total":  total,
		"offset": offset,
		"limit":  limit,
		"raw":    false,
	})
}

func deviceHistoryTimeline(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		fail(c, 400, "device id is required")
		return
	}

	start := queryTime(c, "start")
	end := queryTime(c, "end")
	offset := queryInt(c, "offset", 0)
	limit := queryInt(c, "limit", 16000)
	maxPoints := queryInt(c, "maxPoints", 4000)

	items, total, err := manager.HistoryManagerGlobal.QueryTimeline(manager.HistoryQuery{
		DeviceSn: id,
		Start:    start,
		End:      end,
		Offset:   offset,
		Limit:    limit,
	}, maxPoints)
	if err != nil {
		fail(c, 400, err.Error())
		return
	}

	ok(c, gin.H{
		"items":     items,
		"total":     total,
		"offset":    offset,
		"limit":     limit,
		"maxPoints": maxPoints,
	})
}

func deviceHistoryMetrics(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		fail(c, 400, "device id is required")
		return
	}

	start := queryTime(c, "start")
	end := queryTime(c, "end")
	offset := queryInt(c, "offset", 0)
	limit := queryInt(c, "limit", 16000)
	maxPoints := queryInt(c, "maxPoints", 4000)

	tempIndex := clampMetricIndex(queryInt(c, "tempIndex", 0))
	voltageIndex := clampMetricIndex(queryInt(c, "voltageIndex", 0))
	currentIndex := clampMetricIndex(queryInt(c, "currentIndex", 0))

	items, total, err := manager.HistoryManagerGlobal.QueryMetrics(manager.HistoryMetricQuery{
		DeviceSn:     id,
		Start:        start,
		End:          end,
		Offset:       offset,
		Limit:        limit,
		TempIndex:    tempIndex,
		VoltageIndex: voltageIndex,
		CurrentIndex: currentIndex,
	}, maxPoints)
	if err != nil {
		fail(c, 400, err.Error())
		return
	}

	ok(c, gin.H{
		"items":        items,
		"total":        total,
		"offset":       offset,
		"limit":        limit,
		"maxPoints":    maxPoints,
		"tempIndex":    tempIndex,
		"voltageIndex": voltageIndex,
		"currentIndex": currentIndex,
	})
}

func deviceParamChanges(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		fail(c, 400, "device id is required")
		return
	}

	start := queryTime(c, "start")
	end := queryTime(c, "end")
	pathKeyword := strings.TrimSpace(c.Query("path"))
	offset := queryInt(c, "offset", 0)
	limit := queryInt(c, "limit", 200)

	items, total, err := manager.ParamChangeManagerGlobal.Query(manager.ParamChangeQuery{
		DeviceSn:    id,
		Start:       start,
		End:         end,
		PathKeyword: pathKeyword,
		Offset:      offset,
		Limit:       limit,
	})
	if err != nil {
		fail(c, 400, err.Error())
		return
	}

	ok(c, gin.H{
		"items":  items,
		"total":  total,
		"offset": offset,
		"limit":  limit,
		"path":   pathKeyword,
	})
}

func exportDeviceParamChanges(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		fail(c, 400, "device id is required")
		return
	}

	start := queryTime(c, "start")
	end := queryTime(c, "end")
	pathKeyword := strings.TrimSpace(c.Query("path"))
	limit := queryInt(c, "limit", 5000)

	items, _, err := manager.ParamChangeManagerGlobal.Query(manager.ParamChangeQuery{
		DeviceSn:    id,
		Start:       start,
		End:         end,
		PathKeyword: pathKeyword,
		Offset:      0,
		Limit:       limit,
	})
	if err != nil {
		fail(c, 400, err.Error())
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	// Add UTF-8 BOM for better compatibility with spreadsheet software.
	_, _ = buf.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(buf)
	_ = writer.Write([]string{"deviceSn", "changedAt", "paramPath", "oldValue", "newValue"})
	for _, item := range items {
		_ = writer.Write([]string{
			item.DeviceSn,
			item.ChangedAt.Format(time.RFC3339),
			item.ParamPath,
			item.OldValue,
			item.NewValue,
		})
	}
	writer.Flush()

	fileSafeSN := strings.NewReplacer("/", "_", "\\", "_", ":", "_", " ", "_").Replace(id)
	fileName := fmt.Sprintf("device_changes_%s_%s.csv", fileSafeSN, time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	c.Data(200, "text/csv; charset=utf-8", buf.Bytes())
}

func randomData(c *gin.Context) {
	curr := c.Param("current")
	num, err := strconv.Atoi(curr)
	if err != nil {
		fail(c, 400, "current must be an integer")
		return
	}

	delta := time.Now().UnixNano()%1000 - 500
	now := time.Now().Format("2006/01/02 15:04:05")
	ok(c, gin.H{
		"name": now,
		"data": gin.H{
			"time":  now,
			"value": num + int(delta),
		},
	})
}

func phmOverview(c *gin.Context) {
	window := time.Duration(queryInt(c, "windowHours", 24)) * time.Hour
	limit := queryInt(c, "limit", 100)
	items, err := manager.PHMManagerGlobal.Overview(window, limit)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}

	riskStats := gin.H{
		"low":      0,
		"medium":   0,
		"high":     0,
		"critical": 0,
	}
	for _, item := range items {
		switch item.RiskLevel {
		case manager.PHMRiskLow:
			riskStats["low"] = riskStats["low"].(int) + 1
		case manager.PHMRiskMedium:
			riskStats["medium"] = riskStats["medium"].(int) + 1
		case manager.PHMRiskHigh:
			riskStats["high"] = riskStats["high"].(int) + 1
		default:
			riskStats["critical"] = riskStats["critical"].(int) + 1
		}
	}

	ok(c, gin.H{
		"windowHours": int64(window.Hours()),
		"total":       len(items),
		"riskStats":   riskStats,
		"items":       items,
		"generatedAt": time.Now(),
	})
}

func phmDeviceDetail(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		fail(c, 400, "device id is required")
		return
	}

	window := time.Duration(queryInt(c, "windowHours", 24)) * time.Hour
	maxHistory := queryInt(c, "historyLimit", 600)
	maxChanges := queryInt(c, "changeLimit", 200)

	detail, err := manager.PHMManagerGlobal.DeviceDetail(id, window, maxChanges, maxHistory)
	if err != nil {
		fail(c, 400, err.Error())
		return
	}
	ok(c, detail)
}

type phmBatchDeleteRequest struct {
	DeviceSNs []string `json:"deviceSns"`
}

func phmDeleteSingle(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		fail(c, 400, "device id is required")
		return
	}
	result, err := manager.PHMManagerGlobal.DeleteDevices(c.Request.Context(), []string{id})
	if err != nil {
		fail(c, 400, err.Error())
		return
	}
	ok(c, result)
}

func phmDeleteBatch(c *gin.Context) {
	req := phmBatchDeleteRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "invalid payload")
		return
	}

	result, err := manager.PHMManagerGlobal.DeleteDevices(c.Request.Context(), req.DeviceSNs)
	if err != nil {
		fail(c, 400, err.Error())
		return
	}
	ok(c, result)
}

func queryInt(c *gin.Context, key string, fallback int) int {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return n
}

func queryBool(c *gin.Context, key string, fallback bool) bool {
	raw := strings.TrimSpace(strings.ToLower(c.Query(key)))
	if raw == "" {
		return fallback
	}
	switch raw {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return fallback
	}
}

func queryTime(c *gin.Context, key string) time.Time {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return time.Time{}
	}

	if unix, err := strconv.ParseInt(raw, 10, 64); err == nil {
		if unix > 1e12 {
			return time.UnixMilli(unix)
		}
		return time.Unix(unix, 0)
	}

	if ts, err := time.Parse(time.RFC3339, raw); err == nil {
		return ts
	}
	if ts, err := time.Parse("2006-01-02 15:04:05", raw); err == nil {
		return ts
	}
	if ts, err := time.Parse("2006-01-02", raw); err == nil {
		return ts
	}
	return time.Time{}
}

func clampMetricIndex(value int) int {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}
