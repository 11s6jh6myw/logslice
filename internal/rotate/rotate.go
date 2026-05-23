// Package rotate provides utilities for splitting a slice of log lines
// into multiple named segments based on a fixed line count or byte size.
package rotate

import (
	"fmt"
	"strings"
)

// Segment represents a named chunk of log lines produced by rotation.
type Segment struct {
	Name  string
	Lines []string
}

// Options controls how rotation is performed.
type Options struct {
	// MaxLines splits a new segment every MaxLines lines (0 = disabled).
	MaxLines int
	// MaxBytes splits a new segment when accumulated bytes exceed MaxBytes (0 = disabled).
	MaxBytes int
	// BaseName is the prefix used when naming segments (e.g. "app" → "app-001").
	BaseName string
}

// ByLines splits lines into segments of at most opts.MaxLines each.
// Returns an error if MaxLines is not positive.
func ByLines(lines []string, opts Options) ([]Segment, error) {
	if opts.MaxLines <= 0 {
		return nil, fmt.Errorf("rotate: MaxLines must be > 0, got %d", opts.MaxLines)
	}
	var segments []Segment
	for i := 0; i < len(lines); i += opts.MaxLines {
		end := i + opts.MaxLines
		if end > len(lines) {
			end = len(lines)
		}
		segments = append(segments, Segment{
			Name:  segmentName(opts.BaseName, len(segments)+1),
			Lines: lines[i:end],
		})
	}
	return segments, nil
}

// ByBytes splits lines into segments whose total byte size does not exceed
// opts.MaxBytes. A single line that exceeds MaxBytes is placed in its own segment.
// Returns an error if MaxBytes is not positive.
func ByBytes(lines []string, opts Options) ([]Segment, error) {
	if opts.MaxBytes <= 0 {
		return nil, fmt.Errorf("rotate: MaxBytes must be > 0, got %d", opts.MaxBytes)
	}
	var segments []Segment
	var current []string
	var currentBytes int

	flush := func() {
		if len(current) > 0 {
			segments = append(segments, Segment{
				Name:  segmentName(opts.BaseName, len(segments)+1),
				Lines: current,
			})
			current = nil
			currentBytes = 0
		}
	}

	for _, line := range lines {
		lineBytes := len(line) + 1 // +1 for newline
		if currentBytes+lineBytes > opts.MaxBytes && len(current) > 0 {
			flush()
		}
		current = append(current, line)
		currentBytes += lineBytes
	}
	flush()
	return segments, nil
}

func segmentName(base string, index int) string {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "segment"
	}
	return fmt.Sprintf("%s-%03d", base, index)
}
