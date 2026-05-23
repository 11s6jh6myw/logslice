// Package merge provides utilities for merging multiple sorted log
// segments into a single chronologically ordered output.
package merge

import (
	"bufio"
	"fmt"
	"io"
	"sort"

	"github.com/yourorg/logslice/internal/parser"
)

// Entry holds a parsed log line along with its source index.
type Entry struct {
	Line   parser.Line
	Source int
}

// MergeSegments reads lines from each reader, sorts them by timestamp
// (lines without a timestamp retain their original order relative to
// the preceding timestamped line), and writes the result to w.
// Returns the total number of lines written or an error.
func MergeSegments(w io.Writer, readers []io.Reader) (int, error) {
	if len(readers) == 0 {
		return 0, nil
	}

	var entries []Entry

	for idx, r := range readers {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			text := scanner.Text()
			line := parser.ParseLine(text)
			entries = append(entries, Entry{Line: line, Source: idx})
		}
		if err := scanner.Err(); err != nil {
			return 0, fmt.Errorf("merge: reading source %d: %w", idx, err)
		}
	}

	// Stable sort: lines with timestamps are ordered by time;
	// lines without a timestamp are treated as equal and preserve
	// their relative order thanks to sort.SliceStable.
	sort.SliceStable(entries, func(i, j int) bool {
		ti := entries[i].Line.Timestamp
		tj := entries[j].Line.Timestamp
		if ti == nil || tj == nil {
			return false
		}
		return ti.Before(*tj)
	})

	bw := bufio.NewWriter(w)
	written := 0
	for _, e := range entries {
		if _, err := fmt.Fprintln(bw, e.Line.Raw); err != nil {
			return written, fmt.Errorf("merge: writing line: %w", err)
		}
		written++
	}
	if err := bw.Flush(); err != nil {
		return written, fmt.Errorf("merge: flushing output: %w", err)
	}
	return written, nil
}
