package manager

import (
	"testing"
	"time"

	"qa_test_server/model"
)

func TestQueryByKeywordAndPagination(t *testing.T) {
	m := &Manager{}
	m.Update(model.Device{Sn: "SN-003", Name: "gamma", Last_rx_time: time.Now()})
	m.Update(model.Device{Sn: "SN-001", Name: "alpha", Last_rx_time: time.Now()})
	m.Update(model.Device{Sn: "SN-002", Name: "beta", Last_rx_time: time.Now()})

	items, total := m.Query("SN-0", "", false, 30*time.Second, 1, 1)
	if total != 3 {
		t.Fatalf("expected total=3, got %d", total)
	}
	if len(items) != 1 {
		t.Fatalf("expected page size=1, got %d", len(items))
	}
	if items[0].Sn != "SN-002" {
		t.Fatalf("expected SN-002 on second page, got %s", items[0].Sn)
	}
}

func TestStatsOnlineOffline(t *testing.T) {
	m := &Manager{}
	m.Update(model.Device{Sn: "online", Last_rx_time: time.Now()})
	m.Update(model.Device{Sn: "offline", Last_rx_time: time.Now().Add(-2 * time.Minute)})

	stats := m.Stats(30 * time.Second)
	if stats.Total != 2 {
		t.Fatalf("expected total=2, got %d", stats.Total)
	}
	if stats.Online != 1 {
		t.Fatalf("expected online=1, got %d", stats.Online)
	}
	if stats.Offline != 1 {
		t.Fatalf("expected offline=1, got %d", stats.Offline)
	}
}

func TestQuerySummaryContainsLightweightFields(t *testing.T) {
	m := &Manager{}
	dev := model.Device{
		Sn:           "SN-100",
		Name:         "Alpha",
		Last_rx_time: time.Now(),
	}
	dev.Packet.Femto_input_reg.Bate.Hardware_bate = 321
	dev.Packet.Femto_input_reg.Time.Uptime = [2]uint16{12, 34}
	m.Update(dev)

	items, total := m.QuerySummary("SN-100", "", false, 30*time.Second, 0, 10)
	if total != 1 {
		t.Fatalf("expected total=1, got %d", total)
	}
	if len(items) != 1 {
		t.Fatalf("expected one summary item, got %d", len(items))
	}
	if items[0].Sn != "SN-100" {
		t.Fatalf("expected summary sn SN-100, got %s", items[0].Sn)
	}
	if items[0].Group != "" {
		t.Fatalf("expected empty group by default, got %q", items[0].Group)
	}
	if items[0].Hardware_bate != 321 {
		t.Fatalf("expected hardware version 321, got %d", items[0].Hardware_bate)
	}
	if items[0].Uptime[0] != 12 || items[0].Uptime[1] != 34 {
		t.Fatalf("expected uptime [12 34], got %+v", items[0].Uptime)
	}
	if items[0].Pump_count != 15 {
		t.Fatalf("expected pump count=15, got %d", items[0].Pump_count)
	}
	if items[0].Temp_count != 20 {
		t.Fatalf("expected temp count=20, got %d", items[0].Temp_count)
	}
}

func TestQueryByGroup(t *testing.T) {
	m := &Manager{}
	m.Update(model.Device{Sn: "S-A", Name: "A", Group: "virtual-stress", Last_rx_time: time.Now()})
	m.Update(model.Device{Sn: "S-B", Name: "B", Group: "", Last_rx_time: time.Now()})

	items, total := m.Query("", "virtual-stress", false, 30*time.Second, 0, 10)
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected one item in group filter, total=%d len=%d", total, len(items))
	}
	if items[0].Sn != "S-A" {
		t.Fatalf("expected S-A in virtual-stress group, got %s", items[0].Sn)
	}
}
