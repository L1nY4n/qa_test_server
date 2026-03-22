package manager

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"qa_test_server/model"
)

const maxParamValueLength = 4096

func BuildSystemParamChanges(deviceSn string, oldHoldingReg, newHoldingReg interface{}, changedAt time.Time) []model.DeviceParamChange {
	deviceSn = strings.TrimSpace(deviceSn)
	if deviceSn == "" {
		return nil
	}
	if changedAt.IsZero() {
		changedAt = time.Now()
	}

	oldVal, okOld := normalizeToJSONValue(oldHoldingReg)
	newVal, okNew := normalizeToJSONValue(newHoldingReg)
	if !okOld || !okNew {
		return nil
	}
	if reflect.DeepEqual(oldVal, newVal) {
		return nil
	}

	out := make([]model.DeviceParamChange, 0, 32)
	diffJSON("Femto_holding_reg", oldVal, newVal, deviceSn, changedAt, &out)
	return out
}

func normalizeToJSONValue(v interface{}) (interface{}, bool) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, false
	}
	var out interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, false
	}
	return out, true
}

func diffJSON(path string, oldVal, newVal interface{}, deviceSn string, changedAt time.Time, out *[]model.DeviceParamChange) {
	switch oldTyped := oldVal.(type) {
	case map[string]interface{}:
		newTyped, ok := newVal.(map[string]interface{})
		if !ok {
			appendParamChange(path, oldVal, newVal, deviceSn, changedAt, out)
			return
		}
		keys := unionMapKeys(oldTyped, newTyped)
		for _, key := range keys {
			nextPath := path + "." + key
			o, okOld := oldTyped[key]
			n, okNew := newTyped[key]
			if !okOld {
				appendParamChange(nextPath, nil, n, deviceSn, changedAt, out)
				continue
			}
			if !okNew {
				appendParamChange(nextPath, o, nil, deviceSn, changedAt, out)
				continue
			}
			diffJSON(nextPath, o, n, deviceSn, changedAt, out)
		}
	case []interface{}:
		newTyped, ok := newVal.([]interface{})
		if !ok {
			appendParamChange(path, oldVal, newVal, deviceSn, changedAt, out)
			return
		}
		maxLen := len(oldTyped)
		if len(newTyped) > maxLen {
			maxLen = len(newTyped)
		}
		for i := 0; i < maxLen; i++ {
			nextPath := fmt.Sprintf("%s[%d]", path, i)
			hasOld := i < len(oldTyped)
			hasNew := i < len(newTyped)
			switch {
			case !hasOld && hasNew:
				appendParamChange(nextPath, nil, newTyped[i], deviceSn, changedAt, out)
			case hasOld && !hasNew:
				appendParamChange(nextPath, oldTyped[i], nil, deviceSn, changedAt, out)
			default:
				diffJSON(nextPath, oldTyped[i], newTyped[i], deviceSn, changedAt, out)
			}
		}
	default:
		if !reflect.DeepEqual(oldVal, newVal) {
			appendParamChange(path, oldVal, newVal, deviceSn, changedAt, out)
		}
	}
}

func appendParamChange(path string, oldVal, newVal interface{}, deviceSn string, changedAt time.Time, out *[]model.DeviceParamChange) {
	*out = append(*out, model.DeviceParamChange{
		DeviceSn:  deviceSn,
		ParamPath: path,
		OldValue:  encodeParamValue(oldVal),
		NewValue:  encodeParamValue(newVal),
		ChangedAt: changedAt,
	})
}

func encodeParamValue(v interface{}) string {
	if v == nil {
		return "null"
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	if len(b) <= maxParamValueLength {
		return string(b)
	}
	return string(b[:maxParamValueLength]) + "...(truncated)"
}

func unionMapKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
