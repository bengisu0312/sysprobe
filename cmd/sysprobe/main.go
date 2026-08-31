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
		for _, c := range collector.DefaultCollectors() {
	metrics, err := c.Collect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "collector %s failed: %v\n", c.Name(), err)
		continue
	}
	for _, m := range metrics {
		fmt.Printf("%-20s %8.1f %s\n", m.Name, m.Value, m.Unit)
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
