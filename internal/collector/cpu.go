package collector

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CPUStat struct {
	Total uint64
	Idle  uint64
}

func readCPUStat() (CPUStat, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return CPUStat{}, err
	}
	return parseCPUStat(data)
}

func parseCPUStat(data []byte) (CPUStat, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)

			if len(fields) < 5 {
				return CPUStat{}, fmt.Errorf("hatalı cpu stat formatı")
			}

			var total uint64
			var idle uint64

			for i := 1; i < len(fields); i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					return CPUStat{}, err
				}

				total += val

				if i == 4 {
					idle = val
				}
			}
			return CPUStat{Total: total, Idle: idle}, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return CPUStat{}, err
	}

	return CPUStat{}, fmt.Errorf("cpu satırı bulunamadı")
}

func cpuDelta(prev, curr CPUStat) float64 {
	totalDiff := float64(curr.Total - prev.Total)
	idleDiff := float64(curr.Idle - prev.Idle)

	if totalDiff <= 0 {
		return 0.0
	}

	return ((totalDiff - idleDiff) / totalDiff) * 100.0
}

func CollectCPU() (float64, error) {
	prev, err := readCPUStat()
	if err != nil {
		return 0, err
	}

	time.Sleep(200 * time.Millisecond)

	curr, err := readCPUStat()
	if err != nil {
		return 0, err
	}

	return cpuDelta(prev, curr), nil
}
