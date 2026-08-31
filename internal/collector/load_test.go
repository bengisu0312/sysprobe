package collector

import (
	"math"
	"testing"
)

func TestPerCore(t *testing.T) {
	tests := []struct {
		name string
		load Load
		want float64
	}{
		{"Normal Load", Load{Load1: 2.0, CPUs: 4}, 0.5},
		{"Tam Kapasite", Load{Load1: 4.0, CPUs: 4}, 1.0},
		{"Sıfıra Bölme Koruması", Load{Load1: 1.0, CPUs: 0}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.load.PerCore()
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("PerCore() = %v, Beklenen = %v", got, tt.want)
			}
		})
	}
}

func TestParseLoadavg(t *testing.T) {
	validData := []byte("0.52 0.31 0.24 2/487 12043\n")
	
	load, err := parseLoadavg(validData)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	
	if math.Abs(load.Load1-0.52) > 0.01 {
		t.Errorf("Load1 = %v, Beklenen = 0.52", load.Load1)
	}
}
