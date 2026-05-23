package slicer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Options configures the slicing behavior.
type Options struct {
	From     time.Time
	To       time.Time
	Strict   bool
	Output   io.Writer
}

// Slicer reads a log file and writes matching lines to the output.
type Slicer struct {
	opts Options
}

// New creates a new Slicer with the given options.
func New(opts Options) *Slicer {
	return &Slicer{opts: opts}
}

// SliceFile opens the given file path and slices it according to the options.
func (s *Slicer) SliceFile(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("slicer: open file: %w", err)
	}
	defer f.Close()
	return s.Slice(f)
}

// Slice reads from r, filters lines by the configured time range, and writes
// matching lines to the configured output writer.
func (s *Slicer) Slice(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	var lines []parser.Line

	for scanner.Scan() {
		line := parser.ParseLine(scanner.Text())
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("slicer: scan: %w", err)
	}

	matched := parser.FilterByRange(lines, s.opts.From, s.opts.To, s.opts.Strict)

	w := bufio.NewWriter(s.opts.Output)
	for _, l := range matched {
		if _, err := fmt.Fprintln(w, l.Raw); err != nil {
			return 0, fmt.Errorf("slicer: write: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return 0, fmt.Errorf("slicer: flush: %w", err)
	}

	return len(matched), nil
}
