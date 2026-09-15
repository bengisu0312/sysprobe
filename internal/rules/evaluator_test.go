package rules

import (
	"testing"

	"github.com/bengisu0312/sysprobe/internal/collector"
	"github.com/bengisu0312/sysprobe/internal/config"
)

func TestEvaluate(t *testing.T) {
	cfg := config.DefaultConfig()

	tests := []struct {
		name     string
		metrics  []collector.Metric
		expected Status
	}{
		{
			name: "All OK",
			metrics: []collector.Metric{
				{Name: "cpu_percent", Value: 50.0},
				{Name: "disk_percent", Value: 60.0},
			},
			expected: OK,
		},
		{
			name: "CPU Warning",
			metrics: []collector.Metric{
				{Name: "cpu_percent", Value: 85.0}, // Default Warning is 80
				{Name: "disk_percent", Value: 50.0},
			},
			expected: WARNING,
		},
		{
			name: "Disk Critical overrides CPU Warning",
			metrics: []collector.Metric{
				{Name: "cpu_percent", Value: 85.0}, // WARNING
				{Name: "disk_percent", Value: 95.0}, // CRITICAL -> Worst-of
			},
			expected: CRITICAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := Evaluate(tt.metrics, cfg)
			if status != tt.expected {
				t.Errorf("Beklenen %v, alinan %v", tt.expected, status)
			}
		})
	}
}
