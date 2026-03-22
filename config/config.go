package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	HTTPAddr           string
	TCPAddr            string
	ProxyFromPort      int
	ProxyToPort        int
	DBDSN              string
	DBAutoMigrate      bool
	DeviceActiveWithin time.Duration
	HistorySampleEvery time.Duration
	HistoryRetention   time.Duration
	ChangeLogRetention time.Duration

	AuthSecret           string
	AuthTokenTTL         time.Duration
	DefaultAdminUsername string
	DefaultAdminPassword string
}

func Load() AppConfig {
	return AppConfig{
		HTTPAddr:           normalizeAddr(envOrDefault("QA_HTTP_ADDR", ":8080")),
		TCPAddr:            normalizeAddr(envOrDefault("QA_TCP_ADDR", ":4001")),
		ProxyFromPort:      envInt("QA_PROXY_FROM_PORT", 4003),
		ProxyToPort:        envInt("QA_PROXY_TO_PORT", 7777),
		DBDSN:              envOrDefault("QA_DB_DSN", "root:L1nFen9.com@tcp(localhost:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"),
		DBAutoMigrate:      envBool("QA_DB_AUTOMIGRATE", false),
		DeviceActiveWithin: time.Duration(envInt("QA_DEVICE_ACTIVE_WITHIN_SECONDS", 30)) * time.Second,
		HistorySampleEvery: time.Duration(envInt("QA_HISTORY_SAMPLE_SECONDS", 60)) * time.Second,
		HistoryRetention:   time.Duration(envInt("QA_HISTORY_RETENTION_DAYS", 10)) * 24 * time.Hour,
		ChangeLogRetention: time.Duration(envInt("QA_CHANGE_LOG_RETENTION_DAYS", 30)) * 24 * time.Hour,

		AuthSecret:           envOrDefault("QA_AUTH_SECRET", "qa_test_server_default_secret_change_me"),
		AuthTokenTTL:         time.Duration(envInt("QA_AUTH_TOKEN_TTL_HOURS", 24)) * time.Hour,
		DefaultAdminUsername: envOrDefault("QA_DEFAULT_ADMIN_USERNAME", "admin"),
		DefaultAdminPassword: envOrDefault("QA_DEFAULT_ADMIN_PASSWORD", "Admin@123456"),
	}
}

func envOrDefault(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func envInt(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func envBool(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func normalizeAddr(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return ":8080"
	}
	if strings.HasPrefix(addr, ":") || strings.Contains(addr, ":") {
		return addr
	}
	return ":" + addr
}
