package main

import (
	"fmt"
	"io"
	"os"

	"charm.land/glamour/v2"
)

func main() {
	in, err := readInput(os.Args[1:], os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "md:", err)
		os.Exit(1)
	}

	out, err := glamour.RenderWithEnvironmentConfig(string(in))
	if err != nil {
		fmt.Fprintln(os.Stderr, "md:", err)
		os.Exit(1)
	}

	fmt.Print(out)
}

// readInput returns the Markdown to render. If a file is passed as the first
// argument it is read from disk, otherwise the input is read from stdin.
func readInput(args []string, stdin *os.File) ([]byte, error) {
	if len(args) > 0 {
		//nolint:gosec
		return os.ReadFile(args[0])
	}

	stat, err := stdin.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Mode()&os.ModeCharDevice != 0 {
		return nil, fmt.Errorf("no input: pass a file as argument or pipe text via stdin")
	}

	return io.ReadAll(stdin)
}
