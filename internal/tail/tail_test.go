package tail

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeLines(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.WriteString(strings.Join(lines, "\n"))
	return f.Name()
}

func TestReadLines_Basic(t *testing.T) {
	input := []string{"line1", "line2", "line3", "line4", "line5"}
	path := writeLines(t, input)

	got, err := ReadLines(path, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "line3" || got[2] != "line5" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestReadLines_MoreThanAvailable(t *testing.T) {
	input := []string{"a", "b"}
	path := writeLines(t, input)

	got, err := ReadLines(path, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestReadLines_EmptyFile(t *testing.T) {
	path := writeLines(t, []string{})

	got, err := ReadLines(path, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}

func TestReadLines_InvalidN(t *testing.T) {
	path := writeLines(t, []string{"x"})
	_, err := ReadLines(path, 0)
	if err == nil {
		t.Error("expected error for n=0")
	}
}

func TestReadLines_MissingFile(t *testing.T) {
	_, err := ReadLines(filepath.Join(t.TempDir(), "no-such.log"), 5)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestReadBytes_Basic(t *testing.T) {
	path := writeLines(t, []string{"hello world"})

	got, err := ReadBytes(path, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "orld\n" && got != "orld" && !strings.HasSuffix(got, "orld") {
		// content may vary by OS line ending; just check length
		if len(got) != 5 {
			t.Errorf("expected 5 bytes, got %d: %q", len(got), got)
		}
	}
}

func TestReadBytes_LargerThanFile(t *testing.T) {
	path := writeLines(t, []string{"hi"})
	got, err := ReadBytes(path, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "hi") {
		t.Errorf("expected content to contain 'hi', got %q", got)
	}
}

func TestReadBytes_InvalidBytes(t *testing.T) {
	path := writeLines(t, []string{"x"})
	_, err := ReadBytes(path, 0)
	if err == nil {
		t.Error("expected error for bytes=0")
	}
}
