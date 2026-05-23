// Package slicer provides the core log-slicing engine for logslice.
//
// It combines the parser package's line parsing and range filtering into a
// high-level API that reads from any io.Reader (typically a log file) and
// writes matching lines to any io.Writer.
//
// Basic usage:
//
//	s := slicer.New(slicer.Options{
//		From:   from,
//		To:     to,
//		Strict: true,        // only include lines that have a parseable timestamp
//		Output: os.Stdout,
//	})
//
//	n, err := s.SliceFile("/var/log/app.log")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("wrote %d matching lines\n", n)
//
// When Strict is false, lines without a recognisable timestamp are included
// in the output alongside timestamped lines that fall within the range.
package slicer
