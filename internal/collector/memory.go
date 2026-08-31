package collector

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Memory /proc/meminfo verilerini tutar.
type Memory struct {
	Total     uint64
	Free      uint64
	Available uint64
	SwapTotal uint64
	SwapFree  uint64
}

// Collect /proc/meminfo dosyasını okur ve parseMeminfo'ya gönderir.
func Collect() (Memory, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return Memory{}, err
	}
	return parseMeminfo(data)
}

// parseMeminfo ham bayt verisini alır ve Memory yapısına parse eder.
func parseMeminfo(data []byte) (Memory, error) {
	var mem Memory
	scanner := bufio.NewScanner(bytes.NewReader(data))
	hasMemTotal := false

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			// Eğer kritik bir alanda sayı formatı bozuksa hata dönmeliyiz
			if key == "MemTotal" {
				return Memory{}, err
			}
			continue
		}
		// /proc/meminfo içindeki değerler kB cinsindendir, bayta çeviriyoruz
		val *= 1024

		switch key {
		case "MemTotal":
			mem.Total = val
			hasMemTotal = true
		case "MemFree":
			mem.Free = val
		case "MemAvailable":
			mem.Available = val
		case "SwapTotal":
			mem.SwapTotal = val
		case "SwapFree":
			mem.SwapFree = val
		}
	}

	if err := scanner.Err(); err != nil {
		return Memory{}, err
	}

	// MemTotal yoksa veya 0 ise geçersiz veri olarak kabul edip hata dönüyoruz
	if !hasMemTotal || mem.Total == 0 {
		return Memory{}, fmt.Errorf("geçersiz veya eksik MemTotal")
	}

	return mem, nil
}

// UsedPercent RAM kullanım yüzdesini hesaplar.
func (m Memory) UsedPercent() float64 {
	if m.Total == 0 {
		return 0.0
	}
	used := m.Total - m.Available
	return (float64(used) / float64(m.Total)) * 100.0
}

// SwapUsedPercent Swap kullanım yüzdesini hesaplar.
func (m Memory) SwapUsedPercent() float64 {
	if m.SwapTotal == 0 {
		return 0.0
	}
	used := m.SwapTotal - m.SwapFree
	return (float64(used) / float64(m.SwapTotal)) * 100.0
}

// ==========================================
// Interface Sarmalayıcısı (Wrapper) - Gün 6
// ==========================================

type MemoryCollector struct{}

func (c MemoryCollector) Name() string { return "memory" }

func (c MemoryCollector) Collect() ([]Metric, error) {
	mem, err := Collect()
	if err != nil {
		return nil, err
	}

	metrics := []Metric{
		{Name: "memory_percent", Value: mem.UsedPercent(), Unit: "%"},
	}

	if mem.SwapTotal > 0 {
		metrics = append(metrics, Metric{Name: "swap_percent", Value: mem.SwapUsedPercent(), Unit: "%"})
	}

	return metrics, nil
}
