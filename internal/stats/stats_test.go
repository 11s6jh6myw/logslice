package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/stats"
)

func TestStats_Counters(t *testing.T) {
	s := stats.New()

	s.RecordLine()
	s.RecordLine()
	s.RecordLine()
	s.RecordMatch()
	s.RecordMatch()
	s.RecordSkipped()

	if s.TotalLines != 3 {
		t.Errorf("expected TotalLines=3, got %d", s.TotalLines)
	}
	if s.MatchedLines != 2 {
		t.Errorf("expected MatchedLines=2, got %d", s.MatchedLines)
	}
	if s.SkippedLines != 1 {
		t.Errorf("expected SkippedLines=1, got %d", s.SkippedLines)
	}
}

func TestStats_LevelAndPatternFilter(t *testing.T) {
	s := stats.New()
	s.RecordLevelFilter()
	s.RecordPatternFilter()

	if s.FilteredByLevel != 1 {
		t.Errorf("expected FilteredByLevel=1, got %d", s.FilteredByLevel)
	}
	if s.FilteredByPat != 1 {
		t.Errorf("expected FilteredByPat=1, got %d", s.FilteredByPat)
	}
	// Both should also increment SkippedLines
	if s.SkippedLines != 2 {
		t.Errorf("expected SkippedLines=2, got %d", s.SkippedLines)
	}
}

func TestStats_Finish_SetsDuration(t *testing.T) {
	s := stats.New()
	time.Sleep(5 * time.Millisecond)
	s.Finish()

	if s.Duration < 5*time.Millisecond {
		t.Errorf("expected Duration >= 5ms, got %s", s.Duration)
	}
}

func TestStats_Print_ContainsExpectedFields(t *testing.T) {
	s := stats.New()
	s.RecordLine()
	s.RecordMatch()
	s.RecordLevelFilter()
	s.RecordPatternFilter()
	s.Finish()

	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	for _, want := range []string{
		"total lines",
		"matched lines",
		"skipped lines",
		"by level",
		"by pattern",
		"duration",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q", want)
		}
	}
}
