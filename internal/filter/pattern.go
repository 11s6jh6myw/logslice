// Package filter provides pattern-based log line filtering.
package filter

import (
	"regexp"
)

// Pattern holds a compiled regular expression used to match log lines.
type Pattern struct {
	re      *regexp.Regexp
	negate  bool
}

// NewPattern compiles a regex pattern for log line matching.
// If negate is true, the filter matches lines that do NOT match the pattern.
func NewPattern(expr string, negate bool) (*Pattern, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &Pattern{re: re, negate: negate}, nil
}

// Match returns true if the line satisfies the pattern filter.
func (p *Pattern) Match(line string) bool {
	matched := p.re.MatchString(line)
	if p.negate {
		return !matched
	}
	return matched
}

// FilterLines returns only the lines that satisfy all provided patterns.
// If no patterns are provided, all lines are returned.
func FilterLines(lines []string, patterns []*Pattern) []string {
	if len(patterns) == 0 {
		return lines
	}
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if matchesAll(line, patterns) {
			result = append(result, line)
		}
	}
	return result
}

// matchesAll returns true if the line satisfies every pattern.
func matchesAll(line string, patterns []*Pattern) bool {
	for _, p := range patterns {
		if !p.Match(line) {
			return false
		}
	}
	return true
}
