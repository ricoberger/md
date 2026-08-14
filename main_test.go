package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	input := []byte("# Heading\n\nSome text.\n")

	t.Run("TestRunRaw", func(t *testing.T) {
		t.Setenv("MD_RAW", "true")

		var out bytes.Buffer
		if err := run(input, &out); err != nil {
			t.Fatalf("run returned error: %v", err)
		}

		// With MD_RAW set the input is written verbatim, without terminal
		// rendering.
		if out.String() != string(input) {
			t.Errorf("run = %q, want %q", out.String(), string(input))
		}
	})

	t.Run("TestRunRender", func(t *testing.T) {
		var out bytes.Buffer
		if err := run(input, &out); err != nil {
			t.Fatalf("run returned error: %v", err)
		}

		// Without MD_RAW the input is rendered for the terminal, so the output
		// differs from the raw input but still contains the heading text.
		if out.Len() == 0 {
			t.Fatal("run produced no output")
		}
		if out.String() == string(input) {
			t.Error("run did not render the input")
		}
		if !strings.Contains(out.String(), "Heading") {
			t.Errorf("run output %q does not contain heading text", out.String())
		}
	})

	t.Run("TestRunWordWrap", func(t *testing.T) {
		t.Setenv("MD_WORD_WRAP", "40")

		var out bytes.Buffer
		if err := run(input, &out); err != nil {
			t.Fatalf("run returned error: %v", err)
		}

		if out.Len() == 0 {
			t.Fatal("run produced no output")
		}
	})

	t.Run("TestRunInvalidWordWrap", func(t *testing.T) {
		t.Setenv("MD_WORD_WRAP", "notanumber")

		var out bytes.Buffer
		if err := run(input, &out); err == nil {
			t.Fatal("expected error for invalid MD_WORD_WRAP, got nil")
		}
	})

	t.Run("TestRunFrontMatter", func(t *testing.T) {
		fmInput := []byte("---\ntags: [daily-notes]\n---\n\n# Heading\n\nSome text.\n")

		var out bytes.Buffer
		if err := run(fmInput, &out); err != nil {
			t.Fatalf("run returned error: %v", err)
		}

		if strings.Contains(out.String(), "daily-notes") {
			t.Errorf("run output %q still contains front matter", out.String())
		}
		if !strings.Contains(out.String(), "Heading") {
			t.Errorf("run output %q does not contain heading text", out.String())
		}
	})

	t.Run("TestRunFrontMatterRaw", func(t *testing.T) {
		t.Setenv("MD_RAW", "true")
		fmInput := []byte("---\ntags: [daily-notes]\n---\n\n# Heading\n")

		var out bytes.Buffer
		if err := run(fmInput, &out); err != nil {
			t.Fatalf("run returned error: %v", err)
		}

		// MD_RAW prints the input verbatim, including any front matter.
		if out.String() != string(fmInput) {
			t.Errorf("run = %q, want %q", out.String(), string(fmInput))
		}
	})

	t.Run("TestRunFrontMatterEnabled", func(t *testing.T) {
		t.Setenv("MD_ENABLE_FRONT_MATTER", "true")
		fmInput := []byte("---\ntags: [daily-notes]\n---\n\n# Heading\n")

		var out bytes.Buffer
		if err := run(fmInput, &out); err != nil {
			t.Fatalf("run returned error: %v", err)
		}

		// With MD_ENABLE_FRONT_MATTER the front matter is rendered instead of
		// stripped, so its content still appears in the output.
		if !strings.Contains(out.String(), "daily-notes") {
			t.Errorf("run output %q does not contain front matter", out.String())
		}
	})

	t.Run("TestRunADF", func(t *testing.T) {
		t.Setenv("MD_ADF", "true")
		t.Setenv("MD_RAW", "true")

		var out bytes.Buffer
		err := run([]byte(adfDocument), &out)
		if err != nil {
			t.Fatalf("run returned error: %v", err)
		}

		// MD_ADF converts the ADF document to Markdown, and MD_RAW makes the
		// result deterministic by skipping terminal rendering.
		want := "## Overview\n\nHello from ADF"
		if out.String() != want {
			t.Errorf("run = %q, want %q", out.String(), want)
		}
	})

	t.Run("TestRunADFError", func(t *testing.T) {
		t.Setenv("MD_ADF", "true")

		var out bytes.Buffer
		err := run([]byte("{ invalid json"), &out)
		if err == nil {
			t.Fatal("expected error for invalid ADF, got nil")
		}
	})
}

func TestStripFrontMatter(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "strips front matter",
			in:   "---\ntags: [daily-notes]\n---\n\n# Heading\n",
			want: "# Heading\n",
		},
		{
			name: "strips front matter with closing dots",
			in:   "---\ntags: [a]\n...\n# Heading\n",
			want: "# Heading\n",
		},
		{
			name: "no front matter is unchanged",
			in:   "# Heading\n\nSome text.\n",
			want: "# Heading\n\nSome text.\n",
		},
		{
			name: "horizontal rule is not treated as front matter",
			in:   "Some text.\n\n---\n\nMore text.\n",
			want: "Some text.\n\n---\n\nMore text.\n",
		},
		{
			name: "unterminated front matter is unchanged",
			in:   "---\ntags: [a]\n\n# Heading\n",
			want: "---\ntags: [a]\n\n# Heading\n",
		},
		{
			name: "preserves horizontal rules in body",
			in:   "---\ntags: [a]\n---\n\n# Heading\n\n---\n\nMore.\n",
			want: "# Heading\n\n---\n\nMore.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripFrontMatter([]byte(tt.in))
			if string(got) != tt.want {
				t.Errorf("stripFrontMatter = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadInput(t *testing.T) {
	t.Run("TestReadInputFromFile", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "input.md")
		want := "# Heading\n\nSome text.\n"

		if err := os.WriteFile(path, []byte(want), 0o600); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		got, err := readInput([]string{path}, nil)
		if err != nil {
			t.Fatalf("readInput returned error: %v", err)
		}

		if string(got) != want {
			t.Errorf("readInput = %q, want %q", got, want)
		}
	})

	t.Run("TestReadInputFromFileMissing", func(t *testing.T) {
		_, err := readInput([]string{filepath.Join(t.TempDir(), "missing.md")}, nil)
		if err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	})

	t.Run("TestReadInputFromStdin", func(t *testing.T) {
		want := "piped content\n"

		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}

		go func() {
			_, _ = w.WriteString(want)
			_ = w.Close()
		}()

		got, err := readInput(nil, r)
		if err != nil {
			t.Fatalf("readInput returned error: %v", err)
		}

		if string(got) != want {
			t.Errorf("readInput = %q, want %q", got, want)
		}
	})

	t.Run("TestReadInputNoInput", func(t *testing.T) {
		// A regular file reports a non-character device, but a terminal-like
		// device has ModeCharDevice set. os.DevNull behaves like a char device
		// on the supported platforms, so it triggers the "no input" path.
		f, err := os.Open(os.DevNull)
		if err != nil {
			t.Fatalf("failed to open %s: %v", os.DevNull, err)
		}
		defer f.Close()

		_, err = readInput(nil, f)
		if err == nil {
			t.Fatal("expected error when no input is provided, got nil")
		}
	})
}

// adfDocument is a minimal ADF document used to verify the conversion that runs
// when the MD_ADF environment variable is set to true.
const adfDocument = `{
    "version": 1,
    "type": "doc",
    "content": [
        {
            "type": "heading",
            "attrs": {"level": 2},
            "content": [{"type": "text", "text": "Overview"}]
        },
        {
            "type": "paragraph",
            "content": [{"type": "text", "text": "Hello from ADF"}]
        }
    ]
}`

func TestAdfToMarkdown(t *testing.T) {
	t.Run("TestAdfToMarkdownValid", func(t *testing.T) {
		want := "## Overview\n\nHello from ADF"

		got, err := adfToMarkdown([]byte(adfDocument))
		if err != nil {
			t.Fatalf("adfToMarkdown returned error: %v", err)
		}

		if string(got) != want {
			t.Errorf("adfToMarkdown = %q, want %q", got, want)
		}
	})

	t.Run("TestAdfToMarkdownInvalid", func(t *testing.T) {
		_, err := adfToMarkdown([]byte("{ invalid json"))
		if err == nil {
			t.Fatal("expected error for invalid ADF, got nil")
		}
	})
}
