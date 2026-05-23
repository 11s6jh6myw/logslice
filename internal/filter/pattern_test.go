package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestNewPattern_InvalidRegex(t *testing.T) {
	_, err := filter.NewPattern("[invalid", false)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestPattern_Match_Basic(t *testing.T) {
	p, err := filter.NewPattern(`ERROR`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.Match("2024-01-01 ERROR something went wrong") {
		t.Error("expected match for ERROR line")
	}
	if p.Match("2024-01-01 INFO all good") {
		t.Error("expected no match for INFO line")
	}
}

func TestPattern_Match_Negate(t *testing.T) {
	p, err := filter.NewPattern(`DEBUG`, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.Match("2024-01-01 INFO message") {
		t.Error("expected negated match for non-DEBUG line")
	}
	if p.Match("2024-01-01 DEBUG verbose output") {
		t.Error("expected no match for DEBUG line when negated")
	}
}

func TestFilterLines_NoPatterns(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	got := filter.FilterLines(lines, nil)
	if len(got) != len(lines) {
		t.Errorf("expected %d lines, got %d", len(lines), len(got))
	}
}

func TestFilterLines_MultiplePatterns(t *testing.T) {
	lines := []string{
		"ERROR timeout in service A",
		"ERROR disk full",
		"INFO timeout reached",
		"WARN timeout warning",
	}
	p1, _ := filter.NewPattern(`ERROR`, false)
	p2, _ := filter.NewPattern(`timeout`, false)
	got := filter.FilterLines(lines, []*filter.Pattern{p1, p2})
	if len(got) != 1 {
		t.Fatalf("expected 1 line, got %d: %v", len(got), got)
	}
	if got[0] != "ERROR timeout in service A" {
		t.Errorf("unexpected line: %q", got[0])
	}
}

func TestFilterLines_EmptyInput(t *testing.T) {
	p, _ := filter.NewPattern(`ERROR`, false)
	got := filter.FilterLines([]string{}, []*filter.Pattern{p})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
