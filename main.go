package main

import (
	"bufio"
	"bytes"
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

	if os.Getenv("MD_ENABLE_FRONT_MATTER") != "true" {
		in = stripFrontMatter(in)
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

// stripFrontMatter removes a leading YAML front matter block from the given
// input. A front matter block starts with a line containing only "---" at the
// very beginning of the input and ends at the next line containing only "---"
// (or "..."). If no valid front matter block is found the input is returned
// unchanged.
func stripFrontMatter(in []byte) []byte {
	// The input must start with the "---" delimiter, optionally preceded by a
	// UTF-8 byte order mark.
	rest := bytes.TrimPrefix(in, []byte{0xEF, 0xBB, 0xBF})
	if !bytes.HasPrefix(rest, []byte("---\n")) && !bytes.HasPrefix(rest, []byte("---\r\n")) {
		return in
	}

	scanner := bufio.NewScanner(bytes.NewReader(rest))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	// Skip the opening delimiter.
	if !scanner.Scan() {
		return in
	}

	var body bytes.Buffer
	closed := false
	for scanner.Scan() {
		line := scanner.Text()
		if !closed {
			if line == "---" || line == "..." {
				closed = true
			}
			continue
		}
		body.WriteString(line)
		body.WriteByte('\n')
	}
	if scanner.Err() != nil || !closed {
		return in
	}

	// Trim leading blank lines left after the front matter block.
	return bytes.TrimLeft(body.Bytes(), "\n")
}
