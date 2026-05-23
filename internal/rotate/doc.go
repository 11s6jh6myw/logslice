// Package rotate splits a flat slice of log lines into multiple named
// [Segment] values based on either a maximum line count or a maximum
// byte size per segment.
//
// # Usage
//
// Use [ByLines] when you want each output file to contain at most N lines:
//
//	segs, err := rotate.ByLines(lines, rotate.Options{
//		MaxLines: 500,
//		BaseName: "access",
//	})
//
// Use [ByBytes] when you need to cap the size of each output file:
//
//	segs, err := rotate.ByBytes(lines, rotate.Options{
//		MaxBytes: 1 << 20, // 1 MiB
//		BaseName: "access",
//	})
//
// Segments are named sequentially: "access-001", "access-002", etc.
// Each [Segment] can be passed directly to the writer package for output.
package rotate
