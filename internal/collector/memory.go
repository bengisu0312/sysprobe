package collector

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Memory struct {
	Total     uint64
	Available uint64
	SwapTotal uint64
	SwapFree  uint64
}

// UsedPercent, MemAvailable üzerinden kullanım yüzdesini hesaplar (sıfıra bölme korumalı)
func (m Memory) UsedPercent() float64 {
	if m.Total == 0 {
		return 0
	}
	used := m.Total - m.Available
	return float64(used) / float64(m.Total) * 100.0
}

// SwapUsedPercent, Swap kullanım yüzdesini hesaplar (Swap yoksa 0 döner)
func (m Memory) SwapUsedPercent() float64 {
	if m.SwapTotal == 0 {
		return 0
	}
	swapUsed := m.SwapTotal - m.SwapFree
	return float64(swapUsed) / float64(m.SwapTotal) * 100.0
}

func Collect() (Memory, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return Memory{}, fmt.Errorf("read /proc/meminfo: %w", err)
	}
	return parseMeminfo(data)
}

func parseMeminfo(data []byte) (Memory, error) {
	metrics := make(map[string]uint64)
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		// /proc/meminfo'daki kB değerini bayta çeviriyoruz (1 KiB = 1024 bytes)
		if len(fields) > 2 && fields[2] == "kB" {
			val *= 1024
		}

		metrics[key] = val
	}

	if err := scanner.Err(); err != nil {
		return Memory{}, fmt.Errorf("scan /proc/meminfo: %w", err)
	}

	total, ok := metrics["MemTotal"]
	if !ok {
		return Memory{}, fmt.Errorf("missing required field: MemTotal")
	}

	available, ok := metrics["MemAvailable"]
	if !ok {
		return Memory{}, fmt.Errorf("missing required field: MemAvailable")
	}

	return Memory{
		Total:     total,
		Available: available,
		SwapTotal: metrics["SwapTotal"],
		SwapFree:  metrics["SwapFree"],
	}, nil
}
