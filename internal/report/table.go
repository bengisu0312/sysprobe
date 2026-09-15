package report

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/bengisu0312/sysprobe/internal/collector"
	"github.com/bengisu0312/sysprobe/internal/config"
	"github.com/bengisu0312/sysprobe/internal/rules"
)

// isTerminal, çıktının bir terminale mi yoksa pipe/dosyaya mı gittiğini kontrol eder (test -t 1 mantığı)
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		stat, err := f.Stat()
		if err == nil {
			return (stat.Mode() & os.ModeCharDevice) != 0
		}
	}
	return false
}

func colorize(text string, status rules.Status, tty bool) string {
	if !tty {
		return text // Yönlendirme varsa ANSI renk kodu ekleme
	}
	switch status {
	case rules.OK:
		return fmt.Sprintf("\033[32m%s\033[0m", text) // Yeşil
	case rules.WARNING:
		return fmt.Sprintf("\033[33m%s\033[0m", text) // Sarı
	case rules.CRITICAL:
		return fmt.Sprintf("\033[31m%s\033[0m", text) // Kırmızı
	default:
		return text
	}
}

// PrintTable, metrikleri formatlı bir tablo olarak yazar
func PrintTable(w io.Writer, metrics []collector.Metric, cfg *config.Config) {
	tty := isTerminal(w)
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	fmt.Fprintln(tw, "METRIC\tVALUE\tSTATUS\t")
	fmt.Fprintln(tw, "------\t-----\t------\t")

	for _, m := range metrics {
		// Her bir metriğin durumunu tekil olarak değerlendir
		status := rules.Evaluate([]collector.Metric{m}, cfg)
		coloredStatus := colorize(status.String(), status, tty)

		fmt.Fprintf(tw, "%s\t%.1f %s\t%s\t\n", m.Name, m.Value, m.Unit, coloredStatus)
	}
	tw.Flush()
}
