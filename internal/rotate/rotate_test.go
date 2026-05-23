package rotate

import (
	"testing"
)

func lines(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = "line content here"
	}
	return out
}

func TestByLines_BasicSplit(t *testing.T) {
	segs, err := ByLines(lines(10), Options{MaxLines: 3, BaseName: "app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 4 {
		t.Fatalf("expected 4 segments, got %d", len(segs))
	}
	if len(segs[3].Lines) != 1 {
		t.Errorf("last segment: expected 1 line, got %d", len(segs[3].Lines))
	}
	if segs[0].Name != "app-001" {
		t.Errorf("expected name app-001, got %s", segs[0].Name)
	}
}

func TestByLines_ExactMultiple(t *testing.T) {
	segs, err := ByLines(lines(9), Options{MaxLines: 3, BaseName: "log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(segs))
	}
}

func TestByLines_InvalidMaxLines(t *testing.T) {
	_, err := ByLines(lines(5), Options{MaxLines: 0})
	if err == nil {
		t.Fatal("expected error for MaxLines=0")
	}
}

func TestByLines_EmptyInput(t *testing.T) {
	segs, err := ByLines([]string{}, Options{MaxLines: 10, BaseName: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 0 {
		t.Errorf("expected 0 segments for empty input, got %d", len(segs))
	}
}

func TestByBytes_BasicSplit(t *testing.T) {
	// each line is "line content here" = 17 bytes + 1 newline = 18 bytes
	segs, err := ByBytes(lines(6), Options{MaxBytes: 40, BaseName: "out"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) == 0 {
		t.Fatal("expected at least one segment")
	}
	for _, s := range segs {
		if s.Name == "" {
			t.Error("segment has empty name")
		}
	}
}

func TestByBytes_InvalidMaxBytes(t *testing.T) {
	_, err := ByBytes(lines(3), Options{MaxBytes: -1})
	if err == nil {
		t.Fatal("expected error for MaxBytes=-1")
	}
}

func TestByBytes_EmptyBaseName(t *testing.T) {
	segs, err := ByBytes(lines(2), Options{MaxBytes: 1024, BaseName: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segs))
	}
	if segs[0].Name != "segment-001" {
		t.Errorf("expected segment-001, got %s", segs[0].Name)
	}
}
