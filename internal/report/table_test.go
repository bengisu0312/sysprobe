package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bengisu0312/sysprobe/internal/collector"
	"github.com/bengisu0312/sysprobe/internal/config"
)

func TestPrintTable_Buffer(t *testing.T) {
	var buf bytes.Buffer
	metrics := []collector.Metric{
		{Name: "cpu_percent", Value: 50.0, Unit: "%"},
	}
	cfg := config.DefaultConfig()

	// os.Stdout yerine bytes.Buffer veriyoruz!
	PrintTable(&buf, metrics, cfg)
	
	output := buf.String()
	if !strings.Contains(output, "METRIC") || !strings.Contains(output, "cpu_percent") {
		t.Errorf("Tablo formati hatali veya eksik metrik var:\n%s", output)
	}
}

