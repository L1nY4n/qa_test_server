package manager

import (
	"testing"
	"time"

	"qa_test_server/model"
)

func TestTimelinePointFromHistory(t *testing.T) {
	row := model.DeviceHistory{
		SampledAt: time.Now(),
		Online:    true,
		RawJSON: `{
			"Packet": {
				"Femto_input_reg": {
					"Mon": {
						"Temp": [300, 310, 0],
						"Vol": [1200, 1300, 0],
						"Pump_mon": [
							{"Actual_cur": 900, "Fpga_cur": 890},
							{"Actual_cur": 0, "Fpga_cur": 870}
						]
					}
				}
			}
		}`,
	}

	point := timelinePointFromHistory(row)
	if point.TempAvg != 30.5 {
		t.Fatalf("unexpected temp avg: %v", point.TempAvg)
	}
	if point.TempMax != 31.0 {
		t.Fatalf("unexpected temp max: %v", point.TempMax)
	}
	if point.VoltageAvg != 1250 {
		t.Fatalf("unexpected voltage avg: %v", point.VoltageAvg)
	}
	if point.VoltageMax != 1300 {
		t.Fatalf("unexpected voltage max: %v", point.VoltageMax)
	}
	if point.CurrentAvg != 885 {
		t.Fatalf("unexpected current avg: %v", point.CurrentAvg)
	}
	if point.CurrentMax != 900 {
		t.Fatalf("unexpected current max: %v", point.CurrentMax)
	}
}

func TestDownsampleTimelinePoints(t *testing.T) {
	points := make([]model.DeviceHistoryTimelinePoint, 0, 10)
	base := time.Now().Add(-10 * time.Minute)
	for i := 0; i < 10; i++ {
		points = append(points, model.DeviceHistoryTimelinePoint{
			SampledAt: base.Add(time.Duration(i) * time.Minute),
			TempAvg:   float64(i),
		})
	}

	out := downsampleTimelinePoints(points, 4)
	if len(out) != 4 {
		t.Fatalf("unexpected downsample len: %d", len(out))
	}
	if !out[0].SampledAt.Equal(points[0].SampledAt) {
		t.Fatalf("first point mismatch")
	}
	if !out[len(out)-1].SampledAt.Equal(points[len(points)-1].SampledAt) {
		t.Fatalf("last point mismatch")
	}
}

func TestMetricPointFromHistory(t *testing.T) {
	row := model.DeviceHistory{
		SampledAt: time.Now(),
		Online:    true,
		RawJSON: `{
			"Packet": {
				"Femto_input_reg": {
					"Mon": {
						"Temp": [301, 302, 303],
						"Vol": [1201, 1202, 1203],
						"Pump_mon": [
							{"Actual_cur": 901, "Fpga_cur": 0},
							{"Actual_cur": 902, "Fpga_cur": 0},
							{"Actual_cur": 903, "Fpga_cur": 0}
						]
					}
				}
			}
		}`,
	}

	point := metricPointFromHistory(row, 1, 2, 0)
	if point.Temp == nil || *point.Temp != 30.2 {
		t.Fatalf("unexpected temp: %+v", point.Temp)
	}
	if point.Voltage == nil || *point.Voltage != 1203 {
		t.Fatalf("unexpected voltage: %+v", point.Voltage)
	}
	if point.Current == nil || *point.Current != 901 {
		t.Fatalf("unexpected current: %+v", point.Current)
	}
}

func TestDownsampleMetricPoints(t *testing.T) {
	points := make([]model.DeviceHistoryMetricPoint, 0, 10)
	base := time.Now().Add(-10 * time.Minute)
	for i := 0; i < 10; i++ {
		v := float64(i)
		points = append(points, model.DeviceHistoryMetricPoint{
			SampledAt: base.Add(time.Duration(i) * time.Minute),
			Temp:      &v,
		})
	}

	out := downsampleMetricPoints(points, 4)
	if len(out) != 4 {
		t.Fatalf("unexpected downsample len: %d", len(out))
	}
	if !out[0].SampledAt.Equal(points[0].SampledAt) {
		t.Fatalf("first point mismatch")
	}
	if !out[len(out)-1].SampledAt.Equal(points[len(points)-1].SampledAt) {
		t.Fatalf("last point mismatch")
	}
}
