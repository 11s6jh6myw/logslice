// Package cli provides the command-line interface for logslice.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/slicer"
	"github.com/yourorg/logslice/internal/writer"
)

const timeLayout = "2006-01-02T15:04:05"

// Config holds all parsed CLI options.
type Config struct {
	Input     string
	Output    string
	From      string
	To        string
	Level     string
	Patterns  []string
	Negate    []string
	Strict    bool
}

// Run parses arguments and executes the slicing pipeline.
func Run(args []string) error {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	input := fs.String("input", "", "path to input log file (required)")
	output := fs.String("output", "", "output file path (defaults to stdout)")
	from := fs.String("from", "", "start time filter, format: 2006-01-02T15:04:05")
	to := fs.String("to", "", "end time filter, format: 2006-01-02T15:04:05")
	level := fs.String("level", "", "minimum log level (debug|info|warn|error|fatal)")
	pattern := fs.String("pattern", "", "include lines matching this regex")
	negate := fs.String("negate", "", "exclude lines matching this regex")
	strict := fs.Bool("strict", false, "exclude lines without a timestamp")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	if *input == "" {
		fs.Usage()
		return fmt.Errorf("--input is required")
	}

	var fromTime, toTime time.Time
	var err error
	if *from != "" {
		if fromTime, err = time.Parse(timeLayout, *from); err != nil {
			return fmt.Errorf("invalid --from time: %w", err)
		}
	}
	if *to != "" {
		if toTime, err = time.Parse(timeLayout, *to); err != nil {
			return fmt.Errorf("invalid --to time: %w", err)
		}
	}

	f, err := os.Open(*input)
	if err != nil {
		return fmt.Errorf("opening input: %w", err)
	}
	defer f.Close()

	var patterns []string
	if *pattern != "" {
		patterns = append(patterns, *pattern)
	}
	var negates []string
	if *negate != "" {
		negates = append(negates, *negate)
	}

	pat, err := filter.NewPattern(patterns, negates)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	var lvl filter.Level
	if *level != "" {
		lvl, err = filter.ParseLevel(*level)
		if err != nil {
			return fmt.Errorf("invalid --level: %w", err)
		}
	}

	s := slicer.New(f, fromTime, toTime, *strict)
	lines, err := s.Slice()
	if err != nil {
		return fmt.Errorf("slicing: %w", err)
	}

	lines = filter.FilterLines(lines, pat)
	if *level != "" {
		lines = filter.FilterByLevel(lines, lvl, true)
	}

	w := writer.New(*output)
	return w.WriteSegment(lines)
}
