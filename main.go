package main

import (
	"fmt"
	"io"
	"os"

	"charm.land/glamour/v2"
	adfmarkdown "github.com/ajbeck/adf-to-markdown"
)

func main() {
	in, err := readInput(os.Args[1:], os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "md:", err)
		os.Exit(1)
	}

	if err := run(in, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "md:", err)
		os.Exit(1)
	}
}

// run converts and renders the given Markdown according to the MD_ADF, MD_RAW
// and MD_WORD_WRAP environment variables, writing the result to out.
func run(in []byte, out io.Writer) error {
	if os.Getenv("MD_ADF") == "true" {
		var err error
		in, err = adfToMarkdown(in)
		if err != nil {
			return err
		}
	}

	if os.Getenv("MD_RAW") == "true" {
		fmt.Fprint(out, string(in))
		return nil
	}

	wordWrap := 80
	if w := os.Getenv("MD_WORD_WRAP"); w != "" {
		if _, err := fmt.Sscanf(w, "%d", &wordWrap); err != nil {
			return fmt.Errorf("invalid MD_WORD_WRAP value: %s", w)
		}
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithWordWrap(wordWrap),
		glamour.WithEnvironmentConfig(),
	)
	if err != nil {
		return err
	}

	rendered, err := r.Render(string(in))
	if err != nil {
		return err
	}

	fmt.Fprint(out, rendered)

	return nil
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

// adfToMarkdown converts the given ADF document to Markdown, using the
// github.com/ajbeck/adf-to-markdown package. It returns the converted Markdown
// or an error if the conversion fails.
func adfToMarkdown(adf []byte) ([]byte, error) {
	md, err := adfmarkdown.UnmarshalADF(adf)
	if err != nil {
		return nil, err
	}

	return md, nil
}
