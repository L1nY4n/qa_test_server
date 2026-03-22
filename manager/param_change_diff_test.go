package manager

import (
	"testing"
	"time"
)

type demoHolding struct {
	Mode  int
	Items []int
	Nest  struct {
		Name string
	}
}

func TestBuildSystemParamChanges(t *testing.T) {
	oldVal := demoHolding{Mode: 1, Items: []int{1, 2}}
	oldVal.Nest.Name = "alpha"
	newVal := demoHolding{Mode: 2, Items: []int{1, 3, 4}}
	newVal.Nest.Name = "beta"

	changes := BuildSystemParamChanges("SN-001", oldVal, newVal, time.Unix(100, 0))
	if len(changes) < 3 {
		t.Fatalf("expected at least 3 change logs, got %d", len(changes))
	}

	var foundMode bool
	var foundArray bool
	var foundName bool
	for _, item := range changes {
		switch item.ParamPath {
		case "Femto_holding_reg.Mode":
			foundMode = true
		case "Femto_holding_reg.Items[1]", "Femto_holding_reg.Items[2]":
			foundArray = true
		case "Femto_holding_reg.Nest.Name":
			foundName = true
		}
	}
	if !foundMode {
		t.Fatalf("expected mode change path")
	}
	if !foundArray {
		t.Fatalf("expected array change path")
	}
	if !foundName {
		t.Fatalf("expected nested change path")
	}
}
