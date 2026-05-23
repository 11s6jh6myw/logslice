// Package stats provides collection and reporting of log slicing statistics.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Stats holds counters accumulated during a slice operation.
type Stats struct {
	TotalLines      int
	MatchedLines    int
	SkippedLines    int
	FilteredByLevel int
	FilteredByPat   int
	Duration        time.Duration
	start           time.Time
}

// New creates a new Stats instance and records the start time.
func New() *Stats {
	return &Stats{start: time.Now()}
}

// RecordLine increments the total line counter.
func (s *Stats) RecordLine() {
	s.TotalLines++
}

// RecordMatch increments the matched line counter.
func (s *Stats) RecordMatch() {
	s.MatchedLines++
}

// RecordSkipped increments the skipped line counter.
func (s *Stats) RecordSkipped() {
	s.SkippedLines++
}

// RecordLevelFilter increments the level-filtered counter.
func (s *Stats) RecordLevelFilter() {
	s.FilteredByLevel++
	s.SkippedLines++
}

// RecordPatternFilter increments the pattern-filtered counter.
func (s *Stats) RecordPatternFilter() {
	s.FilteredByPat++
	s.SkippedLines++
}

// Finish records the elapsed duration.
func (s *Stats) Finish() {
	s.Duration = time.Since(s.start)
}

// Print writes a human-readable summary to w.
func (s *Stats) Print(w io.Writer) {
	fmt.Fprintf(w, "--- logslice stats ---\n")
	fmt.Fprintf(w, "  total lines   : %d\n", s.TotalLines)
	fmt.Fprintf(w, "  matched lines : %d\n", s.MatchedLines)
	fmt.Fprintf(w, "  skipped lines : %d\n", s.SkippedLines)
	if s.FilteredByLevel > 0 {
		fmt.Fprintf(w, "    by level    : %d\n", s.FilteredByLevel)
	}
	if s.FilteredByPat > 0 {
		fmt.Fprintf(w, "    by pattern  : %d\n", s.FilteredByPat)
	}
	fmt.Fprintf(w, "  duration      : %s\n", s.Duration.Round(time.Millisecond))
}
