package parser

import (
	"testing"
	"time"
)

type tsCase struct {
	line    string
	wantOK  bool
	wantUTC string // expected time in RFC3339, empty means skip check
}

var timestampCases = []tsCase{
	{
		line:    "2024-03-15T08:32:01Z INFO server started",
		wantOK:  true,
		wantUTC: "2024-03-15T08:32:01Z",
	},
	{
		line:    "2024-03-15T08:32:01.123+02:00 ERROR disk full",
		wantOK:  true,
		wantUTC: "2024-03-15T06:32:01Z",
	},
	{
		line:    "2024-03-15 08:32:01 WARN retry attempt 3",
		wantOK:  true,
		wantUTC: "2024-03-15T08:32:01Z",
	},
	{
		line:    "Mar 15 08:32:01 myhost sshd[1234]: Accepted",
		wantOK:  true,
		wantUTC: "", // year is unknown in syslog format — skip exact check
	},
	{
		line:    "127.0.0.1 - - [15/Mar/2024:08:32:01 +0000] \"GET / HTTP/1.1\" 200",
		wantOK:  true,
		wantUTC: "2024-03-15T08:32:01Z",
	},
	{
		line:    "no timestamp here at all",
		wantOK:  false,
	},
	{
		line:    "",
		wantOK:  false,
	},
}

func TestParseTimestamp(t *testing.T) {
	for _, tc := range timestampCases {
		t.Run(tc.line, func(t *testing.T) {
			got, ok := ParseTimestamp(tc.line)
			if ok != tc.wantOK {
				t.Fatalf("ParseTimestamp(%q) ok=%v, want %v", tc.line, ok, tc.wantOK)
			}
			if !tc.wantOK || tc.wantUTC == "" {
				return
			}
			want, err := time.Parse(time.RFC3339, tc.wantUTC)
			if err != nil {
				t.Fatalf("bad test case wantUTC %q: %v", tc.wantUTC, err)
			}
			if !got.UTC().Equal(want.UTC()) {
				t.Errorf("ParseTimestamp(%q) = %v, want %v", tc.line, got.UTC(), want.UTC())
			}
		})
	}
}
