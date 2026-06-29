package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadInputFromFile(t *testing.T) {
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
}

func TestReadInputFromFileMissing(t *testing.T) {
	_, err := readInput([]string{filepath.Join(t.TempDir(), "missing.md")}, nil)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestReadInputFromStdin(t *testing.T) {
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
}

func TestReadInputNoInput(t *testing.T) {
	// A regular file reports a non-character device, but a terminal-like
	// device has ModeCharDevice set. os.DevNull behaves like a char device on
	// the supported platforms, so it triggers the "no input" path.
	f, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatalf("failed to open %s: %v", os.DevNull, err)
	}
	defer f.Close()

	_, err = readInput(nil, f)
	if err == nil {
		t.Fatal("expected error when no input is provided, got nil")
	}
}
