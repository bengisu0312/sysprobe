package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bengisu0312/sysprobe/internal/collector"
)

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		printUsage(stderr)
		return 2
	}

	subcommand := args[1]

	switch subcommand {
	case "status":
		cpuUsage, err := collector.CollectCPU()
		if err != nil {
			fmt.Fprintf(stderr, "Error collecting CPU: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "CPU:     %.1f%%\n", cpuUsage)

		mem, err := collector.Collect()
		if err != nil {
			fmt.Fprintf(stderr, "Error collecting memory: %v\n", err)
			return 1
		}

		totalGiB := float64(mem.Total) / (1024 * 1024 * 1024)
		availGiB := float64(mem.Available) / (1024 * 1024 * 1024)
		usedGiB := totalGiB - availGiB
		fmt.Fprintf(stdout, "Memory:  %.1f%% (%.1f GiB / %.1f GiB)\n", mem.UsedPercent(), usedGiB, totalGiB)

		if mem.SwapTotal == 0 {
			fmt.Fprintln(stdout, "Swap:     0.0% (swap disabled)")
		} else {
			swapTotalGiB := float64(mem.SwapTotal) / (1024 * 1024 * 1024)
			swapFreeGiB := float64(mem.SwapFree) / (1024 * 1024 * 1024)
			swapUsedGiB := swapTotalGiB - swapFreeGiB
			fmt.Fprintf(stdout, "Swap:     %.1f%% (%.1f GiB / %.1f GiB)\n", mem.SwapUsedPercent(), swapUsedGiB, swapTotalGiB)
		}

		disk, err := collector.CollectDisk("/")
		if err != nil {
			fmt.Fprintf(stderr, "Error collecting disk: %v\n", err)
			return 1
		}
		diskTotalGiB := float64(disk.Total) / (1024 * 1024 * 1024)
		diskUsedGiB := float64(disk.Used) / (1024 * 1024 * 1024)
		fmt.Fprintf(stdout, "Disk /:  %.1f%% (%.1f GiB / %.1f GiB)\n", disk.UsedPercent(), diskUsedGiB, diskTotalGiB)

		load, err := collector.CollectLoad()
		if err != nil {
			fmt.Fprintf(stderr, "Error collecting load: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "Load:    %.2f / %.2f / %.2f  (%.2f per core, %d cpus)\n", load.Load1, load.Load5, load.Load15, load.PerCore(), load.CPUs)

		return 0

	case "version":
		fmt.Fprintln(stdout, "sysprobe dev")
		return 0

	default:
		fmt.Fprintf(stderr, "Error: unknown subcommand '%s'\n\n", subcommand)
		printUsage(stderr)
		return 2
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: sysprobe <status|version>")

