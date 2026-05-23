// Package stats collects and reports operational statistics for a logslice run.
//
// A Stats value tracks how many lines were read, matched, skipped, and filtered
// by level or pattern, along with the total elapsed wall-clock duration.
//
// Typical usage:
//
//	s := stats.New()
//	for _, line := range lines {
//		s.RecordLine()
//		if passesFilter(line) {
//			s.RecordMatch()
//		} else {
//			s.RecordSkipped()
//		}
//	}
//	s.Finish()
//	s.Print(os.Stderr)
package stats
