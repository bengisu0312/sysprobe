package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bengisu0312/sysprobe/internal/collector"
	"github.com/bengisu0312/sysprobe/internal/config"
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
		// Flag okumak için Go'nun standart 'flag' kütüphanesini kullanıyoruz
		statusCmd := flag.NewFlagSet("status", flag.ContinueOnError)
		statusCmd.SetOutput(stderr)
		configFlag := statusCmd.String("config", "", "Config dosyasi yolu")

		// args[2:] çünkü args[0]=sysprobe, args[1]=status
		if err := statusCmd.Parse(args[2:]); err != nil {
			return 2
		}

		// Config'i yüklüyoruz
		cfg, err := config.Load(*configFlag)
		if err != nil {
			fmt.Fprintf(stderr, "Config yukleme hatasi: %v\n", err)
			return 1
		}

		// Config'in yüklendiğini kanıtlamak için ekrana basalım
		fmt.Fprintf(stdout, "[Config Loaded] CPU Warn: %.1f | Disk Warn: %.1f\n\n", *cfg.Thresholds.CPUWarning, *cfg.Thresholds.DiskWarning)

		for _, c := range collector.DefaultCollectors() {
			metrics, err := c.Collect()
			if err != nil {
				fmt.Fprintf(stderr, "collector %s failed: %v\n", c.Name(), err)
				continue
			}
			for _, m := range metrics {
				fmt.Fprintf(stdout, "%-20s %8.1f %s\n", m.Name, m.Value, m.Unit)
			}
		}
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
}
