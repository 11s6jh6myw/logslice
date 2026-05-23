// Package filter provides utilities for filtering log lines by pattern and severity level.
//
// Pattern-based filtering supports regular expressions with optional negation,
// allowing callers to include or exclude lines matching a given expression.
// Multiple patterns can be combined; a line must satisfy all patterns to pass.
//
// Level-based filtering detects common severity keywords (DEBUG, INFO, WARN,
// ERROR, FATAL) within log lines and retains only those at or above a minimum
// severity threshold. Lines with no detectable level are always preserved to
// avoid accidentally discarding unstructured log output.
//
// Typical usage:
//
//	// Keep only ERROR lines mentioning "timeout"
//	p1, _ := filter.NewPattern(`ERROR`, false)
//	p2, _ := filter.NewPattern(`timeout`, false)
//	filtered := filter.FilterLines(lines, []*filter.Pattern{p1, p2})
//
//	// Keep ERROR and above
//	filtered = filter.FilterByLevel(filtered, filter.LevelError)
package filter
