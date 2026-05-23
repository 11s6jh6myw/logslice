package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		input string
		want  filter.Level
	}{
		{"debug", filter.LevelDebug},
		{"INFO", filter.LevelInfo},
		{"Warn", filter.LevelWarn},
		{"WARNING", filter.LevelWarn},
		{"error", filter.LevelError},
		{"FATAL", filter.LevelFatal},
		{"TRACE", filter.LevelUnknown},
	}
	for _, tc := range cases {
		got := filter.ParseLevel(tc.input)
		if got != tc.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestDetectLevel(t *testing.T) {
	cases := []struct {
		line string
		want filter.Level
	}{
		{"2024-01-01 ERROR something failed", filter.LevelError},
		{"2024-01-01 INFO server started", filter.LevelInfo},
		{"2024-01-01 DEBUG low-level detail", filter.LevelDebug},
		{"2024-01-01 plain log line", filter.LevelUnknown},
	}
	for _, tc := range cases {
		got := filter.DetectLevel(tc.line)
		if got != tc.want {
			t.Errorf("DetectLevel(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestFilterByLevel_ErrorAndAbove(t *testing.T) {
	lines := []string{
		"2024-01-01 DEBUG verbose",
		"2024-01-01 INFO started",
		"2024-01-01 ERROR failed",
		"2024-01-01 FATAL crashed",
		"2024-01-01 plain line",
	}
	got := filter.FilterByLevel(lines, filter.LevelError)
	// expects ERROR, FATAL, and the unknown plain line
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(got), got)
	}
}

func TestFilterByLevel_IncludesUnknown(t *testing.T) {
	lines := []string{
		"plain line without level",
		"DEBUG noisy line",
	}
	got := filter.FilterByLevel(lines, filter.LevelInfo)
	if len(got) != 1 {
		t.Fatalf("expected 1 line (unknown), got %d: %v", len(got), got)
	}
	if got[0] != "plain line without level" {
		t.Errorf("unexpected line: %q", got[0])
	}
}
