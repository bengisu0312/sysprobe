package collector

import (
	"math"
	"os"
	"testing"
)

func TestParseMeminfo(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		wantErr       bool
		wantTotal     uint64
		wantAvailable uint64
		wantPercent   float64
		wantSwapPct   float64
	}{
		{
			name: "valid minimal input",
			input: []byte("MemTotal:       1000 kB\n" +
				"MemAvailable:    250 kB\n" +
				"SwapTotal:       100 kB\n" +
				"SwapFree:         40 kB\n"),
			wantErr:       false,
			wantTotal:     1024000,
			wantAvailable: 256000,
			wantPercent:   75.0,
			wantSwapPct:   60.0,
		},
		{
			name:    "empty input",
			input:   []byte(""),
			wantErr: true,
		},
		{
			name:    "missing MemTotal",
			input:   []byte("MemAvailable: 512 kB\n"),
			wantErr: true,
		},
		{
			name:    "invalid number format",
			input:   []byte("MemTotal: abc kB\nMemAvailable: 512 kB\n"),
			wantErr: true,
		},
		{
			name: "line without unit (e.g. HugePages)",
			input: []byte("MemTotal:       1000 kB\n" +
				"MemAvailable:    500 kB\n" +
				"HugePages_Total:     0\n"),
			wantErr:       false,
			wantTotal:     1024000,
			wantAvailable: 512000,
			wantPercent:   50.0,
			wantSwapPct:   0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMeminfo(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// t.Errorf kullanarak tüm alanları aynı anda kontrol ediyoruz
			if got.Total != tt.wantTotal {
				t.Errorf("Total = %d, want %d", got.Total, tt.wantTotal)
			}
			if got.Available != tt.wantAvailable {
				t.Errorf("Available = %d, want %d", got.Available, tt.wantAvailable)
			}

			const epsilon = 0.01
			if math.Abs(got.UsedPercent()-tt.wantPercent) > epsilon {
				t.Errorf("UsedPercent() = %v, want %v", got.UsedPercent(), tt.wantPercent)
			}
			if math.Abs(got.SwapUsedPercent()-tt.wantSwapPct) > epsilon {
				t.Errorf("SwapUsedPercent() = %v, want %v", got.SwapUsedPercent(), tt.wantSwapPct)
			}
		})
	}
}

func TestParseMeminfoRealFixture(t *testing.T) {
	data, err := os.ReadFile("testdata/meminfo.txt")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	m, err := parseMeminfo(data)
	if err != nil {
		t.Fatalf("unexpected error parsing real fixture: %v", err)
	}

	if m.Total == 0 {
		t.Errorf("expected Total > 0, got %d", m.Total)
	}
	if m.Available > m.Total {
		t.Errorf("Available (%d) cannot be greater than Total (%d)", m.Available, m.Total)
	}
	
	p := m.UsedPercent()
	if p < 0 || p > 100 {
		t.Errorf("UsedPercent() = %f, expected between 0 and 100", p)
	}
}

func TestMemoryUsedPercent(t *testing.T) {
	// Sıfır Total durumu (sıfıra bölme koruması)
	mZero := Memory{Total: 0, Available: 0, SwapTotal: 0, SwapFree: 0}
	if got := mZero.UsedPercent(); got != 0.0 {
		t.Errorf("UsedPercent() with zero Total = %f, want 0.0", got)
	}
	if got := mZero.SwapUsedPercent(); got != 0.0 {
		t.Errorf("SwapUsedPercent() with zero SwapTotal = %f, want 0.0", got)
	}
}
