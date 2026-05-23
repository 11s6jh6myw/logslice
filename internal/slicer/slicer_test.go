package slicer_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/slicer"
)

const sampleLog = `2024-01-15T08:00:00Z INFO  service started
2024-01-15T08:05:00Z DEBUG request received
2024-01-15T08:10:00Z ERROR connection timeout
2024-01-15T08:15:00Z INFO  retry attempt 1
2024-01-15T08:20:00Z INFO  service recovered
not a timestamped line
2024-01-15T08:25:00Z DEBUG cleanup done
`

func mustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSlicer_Slice_BasicRange(t *testing.T) {
	var out bytes.Buffer
	s := slicer.New(slicer.Options{
		From:   mustParse("2024-01-15T08:05:00Z"),
		To:     mustParse("2024-01-15T08:15:00Z"),
		Strict: true,
		Output: &out,
	})

	n, err := s.Slice(strings.NewReader(sampleLog))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 matched lines, got %d", n)
	}
	if !strings.Contains(out.String(), "DEBUG request received") {
		t.Error("expected output to contain 'DEBUG request received'")
	}
}

func TestSlicer_Slice_NonStrict_IncludesUntimestamped(t *testing.T) {
	var out bytes.Buffer
	s := slicer.New(slicer.Options{
		From:   mustParse("2024-01-15T08:00:00Z"),
		To:     mustParse("2024-01-15T08:25:00Z"),
		Strict: false,
		Output: &out,
	})

	n, err := s.Slice(strings.NewReader(sampleLog))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 7 {
		t.Errorf("expected 7 matched lines, got %d", n)
	}
}

func TestSlicer_Slice_EmptyInput(t *testing.T) {
	var out bytes.Buffer
	s := slicer.New(slicer.Options{
		From:   mustParse("2024-01-15T08:00:00Z"),
		To:     mustParse("2024-01-15T09:00:00Z"),
		Strict: true,
		Output: &out,
	})

	n, err := s.Slice(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 matched lines, got %d", n)
	}
}
