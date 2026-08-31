package collector

import (
	"math"
	"testing"
)

func TestDiskUsedPercent(t *testing.T) {
	tests := []struct {
		name string
		disk Disk
		want float64
	}{
		{"Normal Kullanım", Disk{Total: 1000, Used: 750, Free: 250}, 75.0},
		{"Boş Disk", Disk{Total: 1000, Used: 0, Free: 1000}, 0.0},
		{"Sıfıra Bölme Koruması", Disk{Total: 0, Used: 0, Free: 0}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.disk.UsedPercent()
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("UsedPercent() = %v, Beklenen = %v", got, tt.want)
			}
		})
	}
}
