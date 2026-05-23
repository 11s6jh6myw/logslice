package filter

import (
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelUnknown Level = -1
)

var levelNames = map[string]Level{
	"DEBUG": LevelDebug,
	"INFO":  LevelInfo,
	"WARN":  LevelWarn,
	"WARNING": LevelWarn,
	"ERROR": LevelError,
	"FATAL": LevelFatal,
}

// ParseLevel parses a level string (case-insensitive) into a Level constant.
func ParseLevel(s string) Level {
	if l, ok := levelNames[strings.ToUpper(s)]; ok {
		return l
	}
	return LevelUnknown
}

// DetectLevel scans a log line for a known severity keyword and returns its Level.
// Returns LevelUnknown if no level keyword is found.
func DetectLevel(line string) Level {
	upper := strings.ToUpper(line)
	for name, level := range levelNames {
		if strings.Contains(upper, name) {
			return level
		}
	}
	return LevelUnknown
}

// FilterByLevel returns lines whose detected level is >= minLevel.
// Lines with LevelUnknown are always included.
func FilterByLevel(lines []string, minLevel Level) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		lvl := DetectLevel(line)
		if lvl == LevelUnknown || lvl >= minLevel {
			result = append(result, line)
		}
	}
	return result
}
