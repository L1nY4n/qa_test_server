package manager

import (
	"fmt"
	"testing"
	"time"
)

func TestStressPulseCleansTemporaryVirtualGroup(t *testing.T) {
	m := &VirtualDeviceManager{}
	prefix := fmt.Sprintf("UTST%d", time.Now().UnixNano()%1000000)

	result, err := m.StressPulse(VirtualDeviceConfig{
		Count:       6,
		IntervalMs:  100,
		Prefix:      prefix,
		NamePrefix:  "UT Stress Device",
		StartIndex:  1,
		MutateParam: true,
		WsBroadcast: false,
		PulseRepeat: 5,
	})
	if err != nil {
		t.Fatalf("stress pulse failed: %v", err)
	}

	if result.Group != "virtual-stress" {
		t.Fatalf("expected group virtual-stress, got %q", result.Group)
	}
	if !result.Cleaned {
		t.Fatalf("expected cleaned=true")
	}
	if result.Generated <= 0 {
		t.Fatalf("expected generated > 0, got %d", result.Generated)
	}

	items, total := ManagerGlabal.Query(prefix, "", false, 30*time.Second, 0, 50)
	if total != 0 || len(items) != 0 {
		t.Fatalf("expected stress virtual devices cleaned up, got total=%d len=%d", total, len(items))
	}
}
