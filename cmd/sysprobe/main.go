package main

import (
	"fmt"
	"os"

	"github.com/bengisu0312/sysprobe/internal/collector"
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	if len(args) < 2 {
		printUsage()
		return 2
	}

	subcommand := args[1]

	switch subcommand {
	case "status":
		mem, err := collector.Collect()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error collecting memory: %v\n", err)
			return 1
		}

		// Bellek hesaplamaları ve GiB dönüşümleri
		totalGiB := float64(mem.Total) / (1024 * 1024 * 1024)
		availGiB := float64(mem.Available) / (1024 * 1024 * 1024)
		usedGiB := totalGiB - availGiB
		fmt.Printf("Memory:  %.1f%% (%.1f GiB / %.1f GiB)\n", mem.UsedPercent(), usedGiB, totalGiB)

		// Swap hesaplamaları
		if mem.SwapTotal == 0 {
			fmt.Println("Swap:     0.0% (swap disabled)")
		} else {
			swapTotalGiB := float64(mem.SwapTotal) / (1024 * 1024 * 1024)
			swapFreeGiB := float64(mem.SwapFree) / (1024 * 1024 * 1024)
			swapUsedGiB := swapTotalGiB - swapFreeGiB
			fmt.Printf("Swap:     %.1f%% (%.1f GiB / %.1f GiB)\n", mem.SwapUsedPercent(), swapUsedGiB, swapTotalGiB)
		}

		return 0

	case "version":
		fmt.Println("sysprobe dev")
		return 0

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown subcommand '%s'\n\n", subcommand)
		printUsage()
		return 2
	}
}

printUsage() {
	fmt.Println("Usage: sysprobe <status|version>")
}
