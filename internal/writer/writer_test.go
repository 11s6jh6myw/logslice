package writer_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/logslice/internal/writer"
)

func TestWriter_WriteSegment_ToStdout(t *testing.T) {
	var buf bytes.Buffer
	w := writer.New(writer.Options{Stdout: &buf})

	lines := []string{
		"2024-01-01T10:00:00Z INFO starting",
		"2024-01-01T10:00:01Z DEBUG ready",
	}

	if err := w.WriteSegment("segment", lines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	for _, line := range lines {
		if !strings.Contains(got, line) {
			t.Errorf("expected output to contain %q, got:\n%s", line, got)
		}
	}
}

func TestWriter_WriteSegment_ToFile(t *testing.T) {
	dir := t.TempDir()
	w := writer.New(writer.Options{
		OutputDir: dir,
		Prefix:    "slice_",
	})

	lines := []string{
		"2024-01-01T10:00:00Z INFO starting",
		"2024-01-01T10:00:01Z DEBUG ready",
	}

	if err := w.WriteSegment("jan01", lines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := filepath.Join(dir, "slice_jan01.log")
	data, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("expected file %q to exist: %v", expectedPath, err)
	}

	content := string(data)
	for _, line := range lines {
		if !strings.Contains(content, line) {
			t.Errorf("expected file to contain %q, got:\n%s", line, content)
		}
	}
}

func TestWriter_WriteSegment_EmptyLines(t *testing.T) {
	var buf bytes.Buffer
	w := writer.New(writer.Options{Stdout: &buf})

	if err := w.WriteSegment("empty", []string{}); err != nil {
		t.Fatalf("unexpected error on empty lines: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestWriter_WriteSegment_CreatesOutputDir(t *testing.T) {
	base := t.TempDir()
	nestedDir := filepath.Join(base, "nested", "output")

	w := writer.New(writer.Options{OutputDir: nestedDir})

	if err := w.WriteSegment("test", []string{"line one"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Errorf("expected output dir %q to be created", nestedDir)
	}
}
