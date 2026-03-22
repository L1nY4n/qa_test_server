package manager

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	psload "github.com/shirou/gopsutil/v3/load"
	psmem "github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
	psproc "github.com/shirou/gopsutil/v3/process"
)

var SystemMonitorManagerGlobal = &SystemMonitorManager{}

type SystemHostInfo struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
	Architecture    string `json:"architecture"`
	BootTime        int64  `json:"bootTime"`
}

type SystemMetricSample struct {
	SampledAt         time.Time `json:"sampledAt"`
	CPUPercent        float64   `json:"cpuPercent"`
	Load1             *float64  `json:"load1"`
	Load5             *float64  `json:"load5"`
	Load15            *float64  `json:"load15"`
	MemoryTotalBytes  uint64    `json:"memoryTotalBytes"`
	MemoryUsedBytes   uint64    `json:"memoryUsedBytes"`
	MemoryUsedPercent float64   `json:"memoryUsedPercent"`
	SwapTotalBytes    uint64    `json:"swapTotalBytes"`
	SwapUsedBytes     uint64    `json:"swapUsedBytes"`
	SwapUsedPercent   float64   `json:"swapUsedPercent"`
	DiskTotalBytes    uint64    `json:"diskTotalBytes"`
	DiskUsedBytes     uint64    `json:"diskUsedBytes"`
	DiskUsedPercent   float64   `json:"diskUsedPercent"`
	NetSentBytes      uint64    `json:"netSentBytes"`
	NetRecvBytes      uint64    `json:"netRecvBytes"`
	NetSentBps        float64   `json:"netSentBps"`
	NetRecvBps        float64   `json:"netRecvBps"`
	ProcessCPUPercent float64   `json:"processCpuPercent"`
	ProcessRSSBytes   uint64    `json:"processRssBytes"`
	Goroutines        int       `json:"goroutines"`
}

type SystemMetricsPayload struct {
	SampleEverySeconds int64                `json:"sampleEverySeconds"`
	DiskPath           string               `json:"diskPath"`
	Host               SystemHostInfo       `json:"host"`
	Current            SystemMetricSample   `json:"current"`
	Points             []SystemMetricSample `json:"points"`
}

type systemNetSnapshot struct {
	SentBytes uint64
	RecvBytes uint64
	At        time.Time
	Valid     bool
}

type SystemMonitorManager struct {
	mu sync.RWMutex

	started     bool
	sampleEvery time.Duration
	maxPoints   int
	diskPath    string

	hostInfo SystemHostInfo
	process  *psproc.Process

	current SystemMetricSample
	points  []SystemMetricSample
	lastNet systemNetSnapshot
}

func (m *SystemMonitorManager) Init(sampleEvery time.Duration, maxPoints int) {
	if sampleEvery < time.Second {
		sampleEvery = 5 * time.Second
	}
	if maxPoints < 60 {
		maxPoints = 720
	}

	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return
	}
	m.sampleEvery = sampleEvery
	m.maxPoints = maxPoints
	m.diskPath = detectSystemDiskPath()
	m.hostInfo = collectSystemHostInfo()
	m.process, _ = psproc.NewProcess(int32(os.Getpid()))
	m.started = true
	m.mu.Unlock()

	m.collectOnce()
	go m.collectLoop()
}

func (m *SystemMonitorManager) Snapshot(points int) SystemMetricsPayload {
	m.mu.RLock()
	defer m.mu.RUnlock()

	limit := points
	if limit <= 0 || limit > len(m.points) {
		limit = len(m.points)
	}

	out := make([]SystemMetricSample, 0, limit)
	if limit > 0 {
		start := len(m.points) - limit
		out = append(out, m.points[start:]...)
	}

	current := m.current
	if current.SampledAt.IsZero() {
		current.SampledAt = time.Now()
	}

	return SystemMetricsPayload{
		SampleEverySeconds: int64(m.sampleEvery.Seconds()),
		DiskPath:           m.diskPath,
		Host:               m.hostInfo,
		Current:            current,
		Points:             out,
	}
}

func (m *SystemMonitorManager) Latest() SystemMetricSample {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

func (m *SystemMonitorManager) collectLoop() {
	ticker := time.NewTicker(m.sampleEvery)
	defer ticker.Stop()

	for range ticker.C {
		m.collectOnce()
	}
}

func (m *SystemMonitorManager) collectOnce() {
	now := time.Now()

	sample := SystemMetricSample{
		SampledAt:  now,
		Goroutines: runtime.NumGoroutine(),
	}

	if percent, err := cpu.Percent(0, false); err == nil && len(percent) > 0 {
		sample.CPUPercent = roundTo2(percent[0])
	}

	if vm, err := psmem.VirtualMemory(); err == nil {
		sample.MemoryTotalBytes = vm.Total
		sample.MemoryUsedBytes = vm.Used
		sample.MemoryUsedPercent = roundTo2(vm.UsedPercent)
	}

	if sw, err := psmem.SwapMemory(); err == nil {
		sample.SwapTotalBytes = sw.Total
		sample.SwapUsedBytes = sw.Used
		sample.SwapUsedPercent = roundTo2(sw.UsedPercent)
	}

	usagePath := m.getDiskPath()
	if du, err := disk.Usage(usagePath); err == nil {
		sample.DiskTotalBytes = du.Total
		sample.DiskUsedBytes = du.Used
		sample.DiskUsedPercent = roundTo2(du.UsedPercent)
	} else if usagePath != "/" {
		if du2, err2 := disk.Usage("/"); err2 == nil {
			sample.DiskTotalBytes = du2.Total
			sample.DiskUsedBytes = du2.Used
			sample.DiskUsedPercent = roundTo2(du2.UsedPercent)
		}
	}

	if ioStats, err := psnet.IOCounters(false); err == nil && len(ioStats) > 0 {
		sample.NetSentBytes = ioStats[0].BytesSent
		sample.NetRecvBytes = ioStats[0].BytesRecv
	}

	if avg, err := psload.Avg(); err == nil {
		l1 := roundTo2(avg.Load1)
		l5 := roundTo2(avg.Load5)
		l15 := roundTo2(avg.Load15)
		sample.Load1 = &l1
		sample.Load5 = &l5
		sample.Load15 = &l15
	}

	m.mu.Lock()
	if m.lastNet.Valid {
		elapsed := now.Sub(m.lastNet.At).Seconds()
		if elapsed > 0 {
			sentDelta := uint64(0)
			recvDelta := uint64(0)
			if sample.NetSentBytes >= m.lastNet.SentBytes {
				sentDelta = sample.NetSentBytes - m.lastNet.SentBytes
			}
			if sample.NetRecvBytes >= m.lastNet.RecvBytes {
				recvDelta = sample.NetRecvBytes - m.lastNet.RecvBytes
			}
			sent := float64(sentDelta)
			recv := float64(recvDelta)
			sample.NetSentBps = roundTo2(sent / elapsed)
			sample.NetRecvBps = roundTo2(recv / elapsed)
		}
	}
	m.lastNet = systemNetSnapshot{
		SentBytes: sample.NetSentBytes,
		RecvBytes: sample.NetRecvBytes,
		At:        now,
		Valid:     true,
	}

	if m.process != nil {
		if cpuPercent, err := m.process.CPUPercent(); err == nil {
			sample.ProcessCPUPercent = roundTo2(cpuPercent)
		}
		if memInfo, err := m.process.MemoryInfo(); err == nil && memInfo != nil {
			sample.ProcessRSSBytes = memInfo.RSS
		}
	}

	m.current = sample
	m.points = append(m.points, sample)
	if len(m.points) > m.maxPoints {
		overflow := len(m.points) - m.maxPoints
		copied := make([]SystemMetricSample, len(m.points)-overflow)
		copy(copied, m.points[overflow:])
		m.points = copied
	}
	m.mu.Unlock()
}

func (m *SystemMonitorManager) getDiskPath() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.diskPath == "" {
		return "/"
	}
	return m.diskPath
}

func detectSystemDiskPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return "/"
	}

	volume := filepath.VolumeName(wd)
	if volume != "" {
		return volume + string(os.PathSeparator)
	}

	return "/"
}

func collectSystemHostInfo() SystemHostInfo {
	info := SystemHostInfo{
		Architecture: runtime.GOARCH,
		OS:           runtime.GOOS,
	}

	if hostInfo, err := host.Info(); err == nil && hostInfo != nil {
		info.Hostname = hostInfo.Hostname
		info.OS = hostInfo.OS
		info.Platform = hostInfo.Platform
		info.PlatformVersion = hostInfo.PlatformVersion
		info.KernelVersion = hostInfo.KernelVersion
		info.BootTime = int64(hostInfo.BootTime)
	}

	return info
}

func roundTo2(value float64) float64 {
	if value < 0 {
		return 0
	}
	return float64(int64(value*100+0.5)) / 100
}
