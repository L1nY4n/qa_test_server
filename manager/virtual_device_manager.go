package manager

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"qa_test_server/model"
)

var VirtualDeviceManagerGlobal = &VirtualDeviceManager{}

type VirtualDeviceConfig struct {
	Count       int    `json:"count"`
	IntervalMs  int    `json:"intervalMs"`
	Prefix      string `json:"prefix"`
	NamePrefix  string `json:"namePrefix"`
	StartIndex  int    `json:"startIndex"`
	Group       string `json:"group"`
	MutateParam bool   `json:"mutateParam"`
	WsBroadcast bool   `json:"wsBroadcast"`
	PulseRepeat int    `json:"pulseRepeat"`
}

type VirtualDeviceStatus struct {
	Running          bool      `json:"running"`
	Count            int       `json:"count"`
	IntervalMs       int       `json:"intervalMs"`
	Prefix           string    `json:"prefix"`
	NamePrefix       string    `json:"namePrefix"`
	StartIndex       int       `json:"startIndex"`
	Group            string    `json:"group"`
	MutateParam      bool      `json:"mutateParam"`
	WsBroadcast      bool      `json:"wsBroadcast"`
	StartedAt        time.Time `json:"startedAt"`
	LastTickAt       time.Time `json:"lastTickAt"`
	UpdatesGenerated int64     `json:"updatesGenerated"`
	UpdatesPerSecond float64   `json:"updatesPerSecond"`
	BroadcastDropped int64     `json:"broadcastDropped"`
	ActiveDevices    int       `json:"activeDevices"`
	SampleSN         string    `json:"sampleSn"`
}

type VirtualPulseResult struct {
	Generated        int           `json:"generated"`
	Elapsed          time.Duration `json:"elapsed"`
	ElapsedMs        int64         `json:"elapsedMs"`
	ElapsedNs        int64         `json:"elapsedNs"`
	UpdatesPerSecond float64       `json:"updatesPerSecond"`
	Group            string        `json:"group,omitempty"`
	Cleaned          bool          `json:"cleaned,omitempty"`
}

type virtualRuntimeConfig struct {
	Count       int
	Interval    time.Duration
	Prefix      string
	NamePrefix  string
	StartIndex  int
	Group       string
	MutateParam bool
	WsBroadcast bool
}

type VirtualDeviceManager struct {
	mu       sync.RWMutex
	running  bool
	cfg      virtualRuntimeConfig
	started  time.Time
	snList   []string
	cancel   context.CancelFunc
	doneChan chan struct{}

	updatesGenerated int64
	broadcastDropped int64
	lastTickUnixNano int64

	broadcastHook func(model.Device) bool
}

func (m *VirtualDeviceManager) SetBroadcastHook(hook func(model.Device) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.broadcastHook = hook
}

func (m *VirtualDeviceManager) Start(cfg VirtualDeviceConfig) (VirtualDeviceStatus, error) {
	runtimeCfg, err := normalizeVirtualConfig(cfg)
	if err != nil {
		return VirtualDeviceStatus{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return m.statusLocked(), errors.New("virtual device generator is already running")
	}

	ctx, cancel := context.WithCancel(context.Background())
	snList := buildSNList(runtimeCfg)
	m.running = true
	m.cfg = runtimeCfg
	m.started = time.Now()
	m.snList = snList
	m.cancel = cancel
	m.doneChan = make(chan struct{})
	atomic.StoreInt64(&m.updatesGenerated, 0)
	atomic.StoreInt64(&m.broadcastDropped, 0)
	atomic.StoreInt64(&m.lastTickUnixNano, m.started.UnixNano())

	go m.run(ctx, runtimeCfg, append([]string(nil), snList...), m.doneChan)

	return m.statusLocked(), nil
}

func (m *VirtualDeviceManager) Stop(remove bool) (VirtualDeviceStatus, error) {
	m.mu.RLock()
	if !m.running {
		m.mu.RUnlock()
		return m.Status(), nil
	}
	cancel := m.cancel
	done := m.doneChan
	snList := append([]string(nil), m.snList...)
	m.mu.RUnlock()

	if cancel != nil {
		cancel()
	}
	if done != nil {
		<-done
	}

	if remove {
		for _, sn := range snList {
			ManagerGlabal.Delete(sn)
		}
	}

	m.mu.Lock()
	m.running = false
	m.cancel = nil
	m.doneChan = nil
	m.snList = nil
	status := m.statusLocked()
	m.mu.Unlock()
	return status, nil
}

func (m *VirtualDeviceManager) Status() VirtualDeviceStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.statusLocked()
}

func (m *VirtualDeviceManager) Pulse(cfg VirtualDeviceConfig) (VirtualPulseResult, error) {
	runtimeCfg, err := normalizeVirtualConfig(cfg)
	if err != nil {
		return VirtualPulseResult{}, err
	}
	repeat := cfg.PulseRepeat
	if repeat <= 0 {
		repeat = 20
	}
	if repeat > 500 {
		repeat = 500
	}

	start := time.Now()
	n := 0
	snList := buildSNList(runtimeCfg)
	for i := 0; i < repeat; i++ {
		n += m.generateBatch(runtimeCfg, snList, int64(i+1), time.Now())
	}
	elapsed := time.Since(start)
	ups := 0.0
	if elapsed > 0 {
		ups = float64(n) / elapsed.Seconds()
	}

	return VirtualPulseResult{
		Generated:        n,
		Elapsed:          elapsed,
		ElapsedMs:        elapsed.Milliseconds(),
		ElapsedNs:        elapsed.Nanoseconds(),
		UpdatesPerSecond: ups,
		Group:            runtimeCfg.Group,
	}, nil
}

func (m *VirtualDeviceManager) StressPulse(cfg VirtualDeviceConfig) (VirtualPulseResult, error) {
	// Stress testing uses an explicit temporary virtual group.
	cfg.Group = "virtual-stress"
	if strings.TrimSpace(cfg.Prefix) == "" {
		cfg.Prefix = "STRESS"
	}
	if strings.TrimSpace(cfg.NamePrefix) == "" {
		cfg.NamePrefix = "Stress Device"
	}

	runtimeCfg, err := normalizeVirtualConfig(cfg)
	if err != nil {
		return VirtualPulseResult{}, err
	}

	repeat := cfg.PulseRepeat
	if repeat <= 0 {
		repeat = 20
	}
	if repeat > 500 {
		repeat = 500
	}

	snList := buildSNList(runtimeCfg)
	snapshot := make(map[string]*model.Device, len(snList))
	for _, sn := range snList {
		if prev, ok := ManagerGlabal.Get(sn); ok && prev != nil {
			cp := *prev
			snapshot[sn] = &cp
		}
	}

	start := time.Now()
	n := 0
	for i := 0; i < repeat; i++ {
		n += m.generateBatch(runtimeCfg, snList, int64(i+1), time.Now())
	}
	elapsed := time.Since(start)
	ups := 0.0
	if elapsed > 0 {
		ups = float64(n) / elapsed.Seconds()
	}

	for _, sn := range snList {
		if prev, ok := snapshot[sn]; ok && prev != nil {
			ManagerGlabal.Set(sn, prev)
			continue
		}
		ManagerGlabal.Delete(sn)
	}

	return VirtualPulseResult{
		Generated:        n,
		Elapsed:          elapsed,
		ElapsedMs:        elapsed.Milliseconds(),
		ElapsedNs:        elapsed.Nanoseconds(),
		UpdatesPerSecond: ups,
		Group:            runtimeCfg.Group,
		Cleaned:          true,
	}, nil
}

func (m *VirtualDeviceManager) run(ctx context.Context, cfg virtualRuntimeConfig, snList []string, done chan struct{}) {
	defer close(done)

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	var tick int64
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			tick++
			count := m.generateBatch(cfg, snList, tick, now)
			atomic.AddInt64(&m.updatesGenerated, int64(count))
			atomic.StoreInt64(&m.lastTickUnixNano, now.UnixNano())
		}
	}
}

func (m *VirtualDeviceManager) generateBatch(cfg virtualRuntimeConfig, snList []string, tick int64, now time.Time) int {
	if len(snList) == 0 {
		return 0
	}

	m.mu.RLock()
	hook := m.broadcastHook
	m.mu.RUnlock()

	for i, sn := range snList {
		index := cfg.StartIndex + i
		dev := buildVirtualDevice(sn, cfg.NamePrefix, cfg.Group, index, tick, cfg.Interval, now, cfg.MutateParam)
		ManagerGlabal.Update(dev)

		if cfg.WsBroadcast && hook != nil {
			if !hook(dev) {
				atomic.AddInt64(&m.broadcastDropped, 1)
			}
		}
	}
	return len(snList)
}

func (m *VirtualDeviceManager) statusLocked() VirtualDeviceStatus {
	lastTickUnix := atomic.LoadInt64(&m.lastTickUnixNano)
	lastTickAt := time.Time{}
	if lastTickUnix > 0 {
		lastTickAt = time.Unix(0, lastTickUnix)
	}

	updates := atomic.LoadInt64(&m.updatesGenerated)
	dropped := atomic.LoadInt64(&m.broadcastDropped)
	ups := 0.0
	if !m.started.IsZero() {
		elapsed := time.Since(m.started).Seconds()
		if elapsed > 0 {
			ups = float64(updates) / elapsed
		}
	}

	sampleSN := ""
	if len(m.snList) > 0 {
		sampleSN = m.snList[0]
	}

	return VirtualDeviceStatus{
		Running:          m.running,
		Count:            m.cfg.Count,
		IntervalMs:       int(m.cfg.Interval / time.Millisecond),
		Prefix:           m.cfg.Prefix,
		NamePrefix:       m.cfg.NamePrefix,
		StartIndex:       m.cfg.StartIndex,
		Group:            m.cfg.Group,
		MutateParam:      m.cfg.MutateParam,
		WsBroadcast:      m.cfg.WsBroadcast,
		StartedAt:        m.started,
		LastTickAt:       lastTickAt,
		UpdatesGenerated: updates,
		UpdatesPerSecond: ups,
		BroadcastDropped: dropped,
		ActiveDevices:    len(m.snList),
		SampleSN:         sampleSN,
	}
}

func normalizeVirtualConfig(cfg VirtualDeviceConfig) (virtualRuntimeConfig, error) {
	count := cfg.Count
	if count <= 0 {
		count = 10
	}
	if count > 5000 {
		return virtualRuntimeConfig{}, errors.New("count must be <= 5000")
	}

	intervalMs := cfg.IntervalMs
	if intervalMs <= 0 {
		intervalMs = 500
	}
	if intervalMs < 10 {
		return virtualRuntimeConfig{}, errors.New("intervalMs must be >= 10")
	}
	if intervalMs > 60_000 {
		return virtualRuntimeConfig{}, errors.New("intervalMs must be <= 60000")
	}

	prefix := strings.ToUpper(strings.TrimSpace(cfg.Prefix))
	if prefix == "" {
		prefix = "VDEV"
	}
	prefix = sanitizeIdentity(prefix)
	if prefix == "" {
		return virtualRuntimeConfig{}, errors.New("prefix contains no valid characters")
	}

	namePrefix := strings.TrimSpace(cfg.NamePrefix)
	if namePrefix == "" {
		namePrefix = "Virtual Device"
	}

	startIndex := cfg.StartIndex
	if startIndex <= 0 {
		startIndex = 1
	}

	return virtualRuntimeConfig{
		Count:       count,
		Interval:    time.Duration(intervalMs) * time.Millisecond,
		Prefix:      prefix,
		NamePrefix:  namePrefix,
		StartIndex:  startIndex,
		Group:       sanitizeGroup(cfg.Group),
		MutateParam: cfg.MutateParam,
		WsBroadcast: cfg.WsBroadcast,
	}, nil
}

func buildSNList(cfg virtualRuntimeConfig) []string {
	sns := make([]string, 0, cfg.Count)
	for i := 0; i < cfg.Count; i++ {
		index := cfg.StartIndex + i
		sns = append(sns, fmt.Sprintf("%s-%05d", cfg.Prefix, index))
	}
	return sns
}

func buildVirtualDevice(sn, namePrefix, group string, index int, tick int64, interval time.Duration, now time.Time, mutate bool) model.Device {
	packet := model.Femto_msg_packed{}
	packet.Femto_holding_reg.Laser_para.Head = 0x55aa
	packet.Femto_holding_reg.Laser_para.Laser_info.Laser_serial = uint16(index % 65535)
	packet.Femto_holding_reg.Laser_para.Laser_info.Laser_Power_level = uint16((tick % 100) + 1)
	copyFixedBytes(packet.Femto_holding_reg.Laser_para.Laser_info.SN[:], sn)
	copyFixedBytes(packet.Femto_holding_reg.Laser_para.Laser_info.Model[:], fmt.Sprintf("%s-%03d", namePrefix, index))

	packet.Femto_input_reg.Bate.Hardware_bate = uint16(100 + (index % 5))
	packet.Femto_input_reg.Status = 1
	packet.Femto_input_reg.Online[0] = 1

	uptimeSeconds := uint32((tick * int64(interval/time.Millisecond)) / 1000)
	packet.Femto_input_reg.Time.Uptime[0] = uint16((uptimeSeconds >> 16) & 0xffff)
	packet.Femto_input_reg.Time.Uptime[1] = uint16(uptimeSeconds & 0xffff)

	for i := range packet.Femto_input_reg.Mon.Temp {
		packet.Femto_input_reg.Mon.Temp[i] = uint16(250 + ((index + i + int(tick)) % 40))
	}
	for i := range packet.Femto_input_reg.Mon.Pump_mon {
		packet.Femto_input_reg.Mon.Pump_mon[i].Pump_sw = 1
		packet.Femto_input_reg.Mon.Pump_mon[i].Actual_cur = uint16(800 + ((index + i + int(tick)) % 120))
		packet.Femto_input_reg.Mon.Pump_mon[i].Fpga_cur = uint16(790 + ((index + i + int(tick)) % 100))
	}

	if mutate {
		packet.Femto_holding_reg.User_para.Freq = uint16(100 + ((index + int(tick)) % 40))
		packet.Femto_holding_reg.User_para.Puse_width = uint16(20 + ((index + int(tick/2)) % 20))
		packet.Femto_holding_reg.User_para.Laser_power = uint16(150 + ((index + int(tick)) % 250))
		packet.Femto_holding_reg.User_para.Div_factor_0 = uint16((tick + int64(index)) % 1024)
	}

	return model.Device{
		Sn:           sn,
		Name:         fmt.Sprintf("%s %03d", namePrefix, index),
		Group:        group,
		Packet:       packet,
		Last_rx_time: now,
	}
}

func copyFixedBytes(dst []uint8, input string) {
	for i := range dst {
		dst[i] = 0
	}
	if input == "" {
		return
	}
	plain := []byte(input)
	max := len(dst)
	if len(plain) < max {
		max = len(plain)
	}
	for i := 0; i < max; i++ {
		dst[i] = plain[i]
	}
}

func sanitizeIdentity(raw string) string {
	if raw == "" {
		return raw
	}
	var b strings.Builder
	b.Grow(len(raw))
	for _, ch := range raw {
		switch {
		case ch >= 'A' && ch <= 'Z':
			b.WriteRune(ch)
		case ch >= 'a' && ch <= 'z':
			b.WriteRune(ch - 32)
		case ch >= '0' && ch <= '9':
			b.WriteRune(ch)
		case ch == '-', ch == '_':
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func sanitizeGroup(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(raw))
	for _, ch := range raw {
		switch {
		case ch >= 'a' && ch <= 'z':
			b.WriteRune(ch)
		case ch >= '0' && ch <= '9':
			b.WriteRune(ch)
		case ch == '-', ch == '_':
			b.WriteRune(ch)
		}
	}
	return b.String()
}
