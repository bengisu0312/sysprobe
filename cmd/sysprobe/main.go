package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bengisu0312/sysprobe/internal/collector"
	"github.com/bengisu0312/sysprobe/internal/config"
	"github.com/bengisu0312/sysprobe/internal/report"
	"github.com/bengisu0312/sysprobe/internal/rules"
)

func main() {
	// os.Exit sadece burada çağrılır. Bu sayede içerideki defer'lar ezilmez.
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
		statusCmd := flag.NewFlagSet("status", flag.ContinueOnError)
		statusCmd.SetOutput(stderr)
		configFlag := statusCmd.String("config", "", "Config dosyasi yolu")

		if err := statusCmd.Parse(args[2:]); err != nil {
			return 2
		}

		cfg, err := config.Load(*configFlag)
		if err != nil {
			fmt.Fprintf(stderr, "Config yukleme hatasi: %v\n", err)
			return 1 // Hata durumu
		}

		fmt.Fprintf(stdout, "[Config Loaded] CPU Warn: %.1f | Disk Warn: %.1f\n\n", *cfg.Thresholds.CPUWarning, *cfg.Thresholds.DiskWarning)

		var allMetrics []collector.Metric

		for _, c := range collector.DefaultCollectors() {
			metrics, err := c.Collect()
			if err != nil {
				fmt.Fprintf(stderr, "collector %s failed: %v\n", c.Name(), err)
				continue
			}
			
			// Toplananları listeye ekle
			allMetrics = append(allMetrics, metrics...)
		} // <-- Döngü burada bitiyor, metrikleri tek tek ekrana basma kısmını sildik!

		// YENİ: İnsan dostu hizalı ve renkli tabloyu stdout'a bas
		fmt.Fprintln(stdout)
		report.PrintTable(stdout, allMetrics, cfg)

		// Kural motorunu çalıştır ve son durumu bul
		overallStatus := rules.Evaluate(allMetrics, cfg)
		fmt.Fprintf(stdout, "\nOVERALL STATUS: %s\n", overallStatus)

		// Nagios standardına uygun Exit Code döndür (OK=0, WARN=1, CRIT=2)
		return int(overallStatus)

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
