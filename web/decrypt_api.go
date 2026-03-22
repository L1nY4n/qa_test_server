package web

import (
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"qa_test_server/manager"

	"github.com/gin-gonic/gin"
)

type decryptTimeKeyRequest struct {
	DeviceSN string `json:"deviceSn"`
	Key      string `json:"key"`
	KeyHex   string `json:"keyHex"`
}

type generateTimeKeyRequest struct {
	DeviceSN string `json:"deviceSn"`
	TimeText string `json:"time"`
}

func decryptTimeKey(c *gin.Context) {
	req := decryptTimeKeyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid decrypt payload")
		return
	}

	claims, _ := currentClaims(c)
	deviceSN := strings.TrimSpace(req.DeviceSN)

	keyValue, mode, err := parseDecryptKey(req.Key, req.KeyHex)
	if err != nil {
		rawInput := strings.TrimSpace(req.Key)
		if strings.TrimSpace(req.KeyHex) != "" {
			rawInput = strings.TrimSpace(req.KeyHex)
		}
		_ = manager.DecryptManagerGlobal.WriteDecodeLog(c.Request.Context(), manager.TimeKeyDecodeLogInput{
			OperatorID:   claims.UserID,
			OperatorName: claims.Username,
			OperatorRole: claims.Role,
			Operation:    "decode",
			DeviceSn:     deviceSN,
			InputMode:    mode,
			KeyRaw:       rawInput,
			SourceIP:     clientIP(c),
			UserAgent:    c.GetHeader("User-Agent"),
			Err:          err,
		})
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	keyForLog := keyValue
	if mode == "hex" {
		keyForLog = strings.TrimSpace(req.KeyHex)
	}

	result, decodeErr := manager.DecryptManagerGlobal.DecodeTimeKey(deviceSN, keyValue)
	logErr := manager.DecryptManagerGlobal.WriteDecodeLog(c.Request.Context(), manager.TimeKeyDecodeLogInput{
		OperatorID:    claims.UserID,
		OperatorName:  claims.Username,
		OperatorRole:  claims.Role,
		Operation:     "decode",
		DeviceSn:      deviceSN,
		InputMode:     mode,
		KeyRaw:        keyForLog,
		SourceIP:      clientIP(c),
		UserAgent:     c.GetHeader("User-Agent"),
		Result:        &result,
		DecodedYear:   result.DecodedYear,
		DecodedMonth:  result.DecodedMonth,
		DecodedDay:    result.DecodedDay,
		DecodedHour:   result.DecodedHour,
		DecodedMinute: result.DecodedMinute,
		DecodedSecond: result.DecodedSecond,
		Err:           decodeErr,
	})
	if logErr != nil {
		// Do not block decode result by logging failures.
	}

	if decodeErr != nil {
		fail(c, http.StatusBadRequest, decodeErr.Error())
		return
	}
	ok(c, result)
}

func generateTimeKey(c *gin.Context) {
	req := generateTimeKeyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid generate payload")
		return
	}

	claims, _ := currentClaims(c)
	deviceSN := strings.TrimSpace(req.DeviceSN)
	timeText := strings.TrimSpace(req.TimeText)

	target, err := parseGenerateTime(timeText)
	if err != nil {
		_ = manager.DecryptManagerGlobal.WriteDecodeLog(c.Request.Context(), manager.TimeKeyDecodeLogInput{
			OperatorID:   claims.UserID,
			OperatorName: claims.Username,
			OperatorRole: claims.Role,
			Operation:    "generate",
			DeviceSn:     deviceSN,
			InputMode:    "time",
			KeyRaw:       timeText,
			SourceIP:     clientIP(c),
			UserAgent:    c.GetHeader("User-Agent"),
			Err:          err,
		})
		fail(c, http.StatusBadRequest, err.Error())
		return
	}

	result, genErr := manager.DecryptManagerGlobal.GenerateTimeKey(deviceSN, target)
	keyForLog := result.Key
	if keyForLog == "" {
		keyForLog = timeText
	}
	_ = manager.DecryptManagerGlobal.WriteDecodeLog(c.Request.Context(), manager.TimeKeyDecodeLogInput{
		OperatorID:    claims.UserID,
		OperatorName:  claims.Username,
		OperatorRole:  claims.Role,
		Operation:     "generate",
		DeviceSn:      deviceSN,
		InputMode:     "time",
		KeyRaw:        keyForLog,
		SourceIP:      clientIP(c),
		UserAgent:     c.GetHeader("User-Agent"),
		DecodedYear:   result.DecodedYear,
		DecodedMonth:  result.DecodedMonth,
		DecodedDay:    result.DecodedDay,
		DecodedHour:   result.DecodedHour,
		DecodedMinute: result.DecodedMinute,
		DecodedSecond: result.DecodedSecond,
		Err:           genErr,
	})

	if genErr != nil {
		fail(c, http.StatusBadRequest, genErr.Error())
		return
	}
	ok(c, result)
}

func decryptLogs(c *gin.Context) {
	offset := queryInt(c, "offset", 0)
	limit := queryInt(c, "limit", 20)
	deviceSN := strings.TrimSpace(c.Query("deviceSn"))
	keyword := strings.TrimSpace(c.Query("keyword"))
	operation := strings.TrimSpace(strings.ToLower(c.Query("operation")))

	items, total, err := manager.DecryptManagerGlobal.QueryLogs(manager.DecryptLogQuery{
		DeviceSn:  deviceSN,
		Keyword:   keyword,
		Operation: operation,
		Offset:    offset,
		Limit:     limit,
	})
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	ok(c, gin.H{
		"items":     items,
		"total":     total,
		"offset":    offset,
		"limit":     limit,
		"deviceSn":  deviceSN,
		"keyword":   keyword,
		"operation": operation,
	})
}

func parseDecryptKey(key, keyHex string) (string, string, error) {
	keyHex = strings.TrimSpace(keyHex)
	if keyHex != "" {
		plainHex := strings.ReplaceAll(keyHex, " ", "")
		plainHex = strings.ReplaceAll(plainHex, "-", "")
		raw, err := hex.DecodeString(plainHex)
		if err != nil {
			return "", "hex", err
		}
		return string(raw), "hex", nil
	}

	key = strings.TrimSpace(key)
	if key == "" {
		return "", "plain", manager.ErrDecryptKeyRequired
	}
	return key, "plain", nil
}

func parseGenerateTime(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, manager.ErrDecryptTimeRequired
	}

	if unix, err := strconv.ParseInt(raw, 10, 64); err == nil {
		if unix > 1e12 {
			return time.UnixMilli(unix), nil
		}
		return time.Unix(unix, 0), nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if ts, err := time.ParseInLocation(layout, raw, time.Local); err == nil {
			return ts, nil
		}
	}
	return time.Time{}, manager.ErrDecryptTimeRequired
}

func clientIP(c *gin.Context) string {
	if ip := strings.TrimSpace(c.ClientIP()); ip != "" {
		return ip
	}
	return strings.TrimSpace(c.Request.RemoteAddr)
}
