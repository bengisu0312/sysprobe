package collector

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}

func TestCPUDelta(t *testing.T) {
	tests := []struct {
		name string
		prev CPUStat
		curr CPUStat
		want float64
	}{
		{"Normal Hesaplama", CPUStat{Total: 1000, Idle: 800}, CPUStat{Total: 1100, Idle: 870}, 30.0},
		{"Sıfıra Bölme Koruması", CPUStat{Total: 1000, Idle: 800}, CPUStat{Total: 1000, Idle: 800}, 0.0},
		{"Tamamen Boş (Idle)", CPUStat{Total: 1000, Idle: 1000}, CPUStat{Total: 1100, Idle: 1100}, 0.0},
		{"Tamamen Dolu (Max Load)", CPUStat{Total: 1000, Idle: 500}, CPUStat{Total: 1100, Idle: 500}, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cpuDelta(tt.prev, tt.curr)
			if !almostEqual(got, tt.want) {
				t.Errorf("cpuDelta() = %v, Beklenen = %v", got, tt.want)
			}
		})
	}
}

func TestParseCPUStat(t *testing.T) {
	validData := []byte("cpu  245789 1523 68432 8934521 12043 0 3421 0 0 0\ncpu0 123 45 67 89 0 0\n")

	stat, err := parseCPUStat(validData)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}

	expectedTotal := uint64(9265729)
	expectedIdle := uint64(8934521)

	if stat.Total != expectedTotal {
		t.Errorf("Total = %v, Beklenen = %v", stat.Total, expectedTotal)
	}
	if stat.Idle != expectedIdle {
		t.Errorf("Idle = %v, Beklenen = %v", stat.Idle, expectedIdle)
	}
}
