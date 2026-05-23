package stats_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/stats"
)

// TestStats_PrintOmitsZeroSubcounters verifies that level/pattern sub-lines
// are omitted when those counters are zero, keeping output concise.
func TestStats_PrintOmitsZeroSubcounters(t *testing.T) {
	s := stats.New()
	s.RecordLine()
	s.RecordMatch()
	s.Finish()

	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	if strings.Contains(out, "by level") {
		t.Error("expected 'by level' to be absent when FilteredByLevel==0")
	}
	if strings.Contains(out, "by pattern") {
		t.Error("expected 'by pattern' to be absent when FilteredByPat==0")
	}
}

// TestStats_ZeroValue exercises a Stats that has Finish called immediately.
func TestStats_ZeroValue(t *testing.T) {
	s := stats.New()
	s.Finish()

	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "total lines   : 0") {
		t.Errorf("expected zero total lines in output, got:\n%s", out)
	}
	if !strings.Contains(out, "matched lines : 0") {
		t.Errorf("expected zero matched lines in output, got:\n%s", out)
	}
}
