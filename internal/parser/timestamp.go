package parser

import (
	"fmt"
	"regexp"
	"time"
)

// Common log timestamp formats to try in order
var timestampFormats = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.000Z07:00",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
	"Jan  2 15:04:05",
}

// timestampPattern matches common timestamp prefixes in log lines
var timestampPattern = regexp.MustCompile(
	`(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?|` +
		`\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2}|` +
		`\d{2}/\w{3}/\d{4}:\d{2}:\d{2}:\d{2}\s[+-]\d{4})`,
)

// ParseTimestamp extracts and parses the first timestamp found in a log line.
// Returns the parsed time and true on success, or zero time and false if no
// recognisable timestamp is present.
func ParseTimestamp(line string) (time.Time, bool) {
	match := timestampPattern.FindString(line)
	if match == "" {
		return time.Time{}, false
	}

	for _, format := range timestampFormats {
		t, err := time.Parse(format, match)
		if err == nil {
			return t, true
		}
	}

	return time.Time{}, false
}

// ErrNoTimestamp is returned when a line contains no parseable timestamp.
var ErrNoTimestamp = fmt.Errorf("no parseable timestamp found in line")
