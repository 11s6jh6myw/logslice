// Package tail provides utilities for reading the trailing portion
// of a log file by line count or byte count.
//
// It is designed to complement logslice's segmentation pipeline by
// allowing callers to quickly inspect the end of a large log before
// or after slicing, without loading the entire file into memory.
//
// # Usage
//
//	// Last 20 lines
//	lines, err := tail.ReadLines("/var/log/app.log", 20)
//
//	// Last 4 KiB
//	chunk, err := tail.ReadBytes("/var/log/app.log", 4096)
//
// ReadLines uses an in-memory ring buffer so memory usage is bounded
// by n regardless of file size. ReadBytes seeks directly to the
// appropriate file offset for O(1) positioning.
package tail
