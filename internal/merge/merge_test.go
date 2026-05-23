package merge_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/merge"
)

func reader(s string) *strings.Reader {
	return strings.NewReader(s)
}

func TestMergeSegments_EmptyInput(t *testing.T) {
	var buf bytes.Buffer
	n, err := merge.MergeSegments(&buf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Fatalf("expected 0 lines, got %d", n)
	}
}

func TestMergeSegments_SingleReader(t *testing.T) {
	input := "2024-01-01T10:00:00Z INFO hello\n2024-01-01T10:01:00Z INFO world\n"
	var buf bytes.Buffer
	n, err := merge.MergeSegments(&buf, []interface{ Read([]byte) (int, error) }{reader(input)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 lines, got %d", n)
	}
}

func TestMergeSegments_MergesChronologically(t *testing.T) {
	seg1 := "2024-01-01T10:00:00Z INFO first\n2024-01-01T10:02:00Z INFO third\n"
	seg2 := "2024-01-01T10:01:00Z INFO second\n2024-01-01T10:03:00Z INFO fourth\n"

	var buf bytes.Buffer
	_, err := merge.MergeSegments(&buf, []interface{ Read([]byte) (int, error) }{
		reader(seg1),
		reader(seg2),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}
	expected := []string{"first", "second", "third", "fourth"}
	for i, want := range expected {
		if !strings.Contains(lines[i], want) {
			t.Errorf("line %d: want %q in %q", i, want, lines[i])
		}
	}
}

func TestMergeSegments_UntimestampedLinesPreserveOrder(t *testing.T) {
	seg1 := "2024-01-01T10:00:00Z INFO start\nno timestamp line\n"
	seg2 := "2024-01-01T10:01:00Z INFO end\n"

	var buf bytes.Buffer
	n, err := merge.MergeSegments(&buf, []interface{ Read([]byte) (int, error) }{
		reader(seg1),
		reader(seg2),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3 lines, got %d", n)
	}
}
