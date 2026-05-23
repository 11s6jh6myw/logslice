// Package tail provides utilities for reading the last N lines
// or bytes from a log file, useful for previewing log segments.
package tail

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Options configures tail behaviour.
type Options struct {
	// Lines is the number of trailing lines to return.
	// Ignored when Bytes > 0.
	Lines int
	// Bytes is the number of trailing bytes to return.
	// When non-zero, Lines is ignored.
	Bytes int64
}

// ReadLines returns the last n lines from the named file.
func ReadLines(path string, n int) ([]string, error) {
	if n <= 0 {
		return nil, fmt.Errorf("tail: n must be greater than zero")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tail: open %s: %w", path, err)
	}
	defer f.Close()
	return readLastLines(f, n)
}

// ReadBytes returns the last b bytes from the named file as a string.
func ReadBytes(path string, b int64) (string, error) {
	if b <= 0 {
		return "", fmt.Errorf("tail: bytes must be greater than zero")
	}
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("tail: open %s: %w", path, err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	offset := fi.Size() - b
	if offset < 0 {
		offset = 0
	}
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return "", err
	}
	buf, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// readLastLines reads the last n lines from r using a ring buffer.
func readLastLines(r io.Reader, n int) ([]string, error) {
	ring := make([]string, n)
	idx := 0
	count := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ring[idx%n] = scanner.Text()
		idx++
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if count == 0 {
		return []string{}, nil
	}
	start := idx % n
	result := make([]string, 0, min(count, n))
	for i := 0; i < min(count, n); i++ {
		result = append(result, ring[(start+i)%n])
	}
	return result, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
