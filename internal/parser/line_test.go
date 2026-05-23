package parser

import (
	"testing"
	"time"
)

func TestParseLine_WithTimestamp(t *testing.T) {
	raw := "2024-01-15T10:30:00Z INFO starting server"
	l := ParseLine(raw)

	if !l.HasTime {
		t.Fatal("expected HasTime=true")
	}
	if l.Raw != raw {
		t.Errorf("expected Raw=%q, got %q", raw, l.Raw)
	}
	expected := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !l.Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, l.Timestamp)
	}
}

func TestParseLine_WithoutTimestamp(t *testing.T) {
	raw := "no timestamp here just plain text"
	l := ParseLine(raw)

	if l.HasTime {
		t.Error("expected HasTime=false for line without timestamp")
	}
	if l.Raw != raw {
		t.Errorf("expected Raw=%q, got %q", raw, l.Raw)
	}
}

func TestParseLine_EmptyLine(t *testing.T) {
	l := ParseLine("")
	if l.HasTime {
		t.Error("expected HasTime=false for empty line")
	}
}

func TestFilterByRange_Strict(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	lines := []LogLine{
		{Timestamp: time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC), HasTime: true, Raw: "before"},
		{Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), HasTime: true, Raw: "at start"},
		{Timestamp: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC), HasTime: true, Raw: "inside"},
		{Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), HasTime: true, Raw: "at end"},
		{Timestamp: time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC), HasTime: true, Raw: "after"},
		{HasTime: false, Raw: "no time"},
	}

	got := FilterByRange(lines, start, end, true)
	if len(got) != 3 {
		t.Errorf("strict: expected 3 lines, got %d", len(got))
	}
}

func TestFilterByRange_NonStrict(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	lines := []LogLine{
		{Timestamp: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC), HasTime: true, Raw: "inside"},
		{HasTime: false, Raw: "no time"},
	}

	got := FilterByRange(lines, start, end, false)
	if len(got) != 2 {
		t.Errorf("non-strict: expected 2 lines, got %d", len(got))
	}
}
