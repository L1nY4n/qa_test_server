package manager

import (
	"testing"
	"time"

	"qa_test_server/model"
)

func TestHistoryRowToDeviceFromRaw(t *testing.T) {
	last := time.Now().Add(-time.Minute).Round(time.Second)
	row := model.DeviceHistory{
		DeviceSn:   "SN-001",
		DeviceName: "Fallback",
		SampledAt:  time.Now(),
		RawJSON:    `{"Sn":"SN-001","Name":"Laser A","Last_rx_time":"` + last.Format(time.RFC3339) + `"}`,
	}

	device, ok := historyRowToDevice(row)
	if !ok {
		t.Fatalf("expected row to be converted")
	}
	if device.Sn != "SN-001" {
		t.Fatalf("unexpected sn: %s", device.Sn)
	}
	if device.Name != "Laser A" {
		t.Fatalf("unexpected name: %s", device.Name)
	}
	if !device.Last_rx_time.Equal(last) {
		t.Fatalf("unexpected last_rx_time: %v", device.Last_rx_time)
	}
}

func TestHistoryRowToDeviceFallbackWhenRawInvalid(t *testing.T) {
	sampledAt := time.Now().Add(-2 * time.Hour).Round(time.Second)
	row := model.DeviceHistory{
		DeviceSn:   "SN-002",
		DeviceName: "Recovered Device",
		SampledAt:  sampledAt,
		RawJSON:    "{broken json}",
	}

	device, ok := historyRowToDevice(row)
	if !ok {
		t.Fatalf("expected row to fallback")
	}
	if device.Sn != row.DeviceSn {
		t.Fatalf("expected fallback sn %s, got %s", row.DeviceSn, device.Sn)
	}
	if device.Name != row.DeviceName {
		t.Fatalf("expected fallback name %s, got %s", row.DeviceName, device.Name)
	}
	if !device.Last_rx_time.Equal(sampledAt) {
		t.Fatalf("expected fallback time %v, got %v", sampledAt, device.Last_rx_time)
	}
}

func TestHistoryRowToDeviceRejectEmptySN(t *testing.T) {
	row := model.DeviceHistory{DeviceSn: "   "}
	_, ok := historyRowToDevice(row)
	if ok {
		t.Fatalf("expected empty sn row to be rejected")
	}
}
