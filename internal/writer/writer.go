// Package writer handles writing sliced log segments to output destinations.
package writer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Options configures the writer behavior.
type Options struct {
	// OutputDir is the directory where output files are written.
	// If empty, output is written to Stdout.
	OutputDir string

	// Prefix is prepended to output filenames.
	Prefix string

	// Stdout is used when OutputDir is empty. Defaults to os.Stdout.
	Stdout io.Writer
}

// Writer writes log lines to an output destination.
type Writer struct {
	opts Options
}

// New creates a new Writer with the given options.
func New(opts Options) *Writer {
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	return &Writer{opts: opts}
}

// WriteSegment writes a slice of log lines to a destination.
// If OutputDir is set, lines are written to a file named "<prefix><name>.log".
// Otherwise, lines are written to Stdout.
func (w *Writer) WriteSegment(name string, lines []string) error {
	if w.opts.OutputDir != "" {
		return w.writeToFile(name, lines)
	}
	return w.writeToStream(w.opts.Stdout, lines)
}

func (w *Writer) writeToFile(name string, lines []string) error {
	if err := os.MkdirAll(w.opts.OutputDir, 0o755); err != nil {
		return fmt.Errorf("writer: create output dir: %w", err)
	}

	filename := fmt.Sprintf("%s%s.log", w.opts.Prefix, name)
	path := filepath.Join(w.opts.OutputDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("writer: create file %q: %w", path, err)
	}
	defer f.Close()

	return w.writeToStream(f, lines)
}

func (w *Writer) writeToStream(dst io.Writer, lines []string) error {
	bw := bufio.NewWriter(dst)
	for _, line := range lines {
		if _, err := fmt.Fprintln(bw, line); err != nil {
			return fmt.Errorf("writer: write line: %w", err)
		}
	}
	return bw.Flush()
}
