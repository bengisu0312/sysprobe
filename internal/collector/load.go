package collector

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Load sistem yükü verilerini ve çekirdek sayısını tutar.
type Load struct {
	Load1  float64
	Load5  float64
	Load15 float64
	CPUs   int
}

// PerCore load'u çekirdek sayısına bölerek gerçek doygunluğu verir.
func (l Load) PerCore() float64 {
	// Sıfıra bölme koruması
	if l.CPUs == 0 {
		return 0.0
	}
	return l.Load1 / float64(l.CPUs)
}

// CollectLoad /proc/loadavg dosyasını okur ve parseLoadavg'a gönderir.
func CollectLoad() (Load, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return Load{}, err
	}
	return parseLoadavg(data)
}

// parseLoadavg ham metni alır ve ilk 3 float değeri çeker.
func parseLoadavg(data []byte) (Load, error) {
	line := strings.TrimSpace(string(data))
	fields := strings.Fields(line)

	if len(fields) < 3 {
		return Load{}, fmt.Errorf("hatalı loadavg formatı")
	}

	load1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return Load{}, err
	}

	load5, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return Load{}, err
	}

	load15, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return Load{}, err
	}

	// runtime.NumCPU() Go'nun standart kütüphanesinden çekirdek sayısını verir
	return Load{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
		CPUs:   runtime.NumCPU(),
	}, nil
}
