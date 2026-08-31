package main

import (
	"flag"
	"fmt"
	"os"
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
		statusCmd := flag.NewFlagSet("status", flag.ContinueOnError)
		_ = statusCmd.String("format", "text", "output format (text|json)")

		if err := statusCmd.Parse(args[2:]); err != nil {
			return 2
		}

		fmt.Println("sysprobe: status command not implemented yet")
		return 3

	case "version":
		fmt.Println("sysprobe dev")
		return 0

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown subcommand '%s'\n\n", subcommand)
		printUsage()
		return 2
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: sysprobe <subcommand> [flags]\n\n")
	fmt.Fprintf(os.Stderr, "Available subcommands:\n")
	fmt.Fprintf(os.Stderr, "  status     Check system metrics\n")
	fmt.Fprintf(os.Stderr, "  version    Print version information\n")
}
