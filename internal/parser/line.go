package parser

import (
	"strings"
	"time"
)

// LogLine represents a parsed log line with its timestamp and raw content.
type LogLine struct {
	Timestamp time.Time
	Raw       string
	HasTime   bool
}

// ParseLine attempts to extract a timestamp from a raw log line.
// It tries common log prefixes and returns a LogLine with HasTime=false
// if no timestamp could be found.
func ParseLine(raw string) LogLine {
	line := LogLine{Raw: raw}

	if raw == "" {
		return line
	}

	// Try to parse a timestamp from the beginning of the line.
	// We attempt progressively shorter prefixes to find a timestamp token.
	fields := strings.Fields(raw)
	for i := 1; i <= len(fields) && i <= 4; i++ {
		candidate := strings.Join(fields[:i], " ")
		ts, err := ParseTimestamp(candidate)
		if err == nil {
			line.Timestamp = ts
			line.HasTime = true
			return line
		}
	}

	return line
}

// FilterByRange returns only the log lines whose timestamps fall within
// [start, end] (inclusive). Lines without a parsed timestamp are excluded
// when strict is true, or included when strict is false.
func FilterByRange(lines []LogLine, start, end time.Time, strict bool) []LogLine {
	result := make([]LogLine, 0, len(lines))
	for _, l := range lines {
		if !l.HasTime {
			if !strict {
				result = append(result, l)
			}
			continue
		}
		if (l.Timestamp.Equal(start) || l.Timestamp.After(start)) &&
			(l.Timestamp.Equal(end) || l.Timestamp.Before(end)) {
			result = append(result, l)
		}
	}
	return result
}
