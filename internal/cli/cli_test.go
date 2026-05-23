package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/cli"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

const sampleLog = `2024-01-15T10:00:00 INFO  service started
2024-01-15T10:01:00 DEBUG checking config
2024-01-15T10:02:00 ERROR connection refused
2024-01-15T10:03:00 WARN  retrying request
`

func TestRun_MissingInput(t *testing.T) {
	err := cli.Run([]string{})
	if err == nil || !strings.Contains(err.Error(), "--input is required") {
		t.Errorf("expected --input required error, got: %v", err)
	}
}

func TestRun_InvalidFromTime(t *testing.T) {
	f := writeTempLog(t, sampleLog)
	err := cli.Run([]string{"--input", f, "--from", "not-a-time"})
	if err == nil || !strings.Contains(err.Error(), "invalid --from time") {
		t.Errorf("expected invalid --from error, got: %v", err)
	}
}

func TestRun_InvalidLevel(t *testing.T) {
	f := writeTempLog(t, sampleLog)
	err := cli.Run([]string{"--input", f, "--level", "verbose"})
	if err == nil || !strings.Contains(err.Error(), "invalid --level") {
		t.Errorf("expected invalid --level error, got: %v", err)
	}
}

func TestRun_BasicSlice_ToFile(t *testing.T) {
	input := writeTempLog(t, sampleLog)
	outDir := t.TempDir()
	output := filepath.Join(outDir, "out.log")

	err := cli.Run([]string{
		"--input", input,
		"--output", output,
		"--from", "2024-01-15T10:01:00",
		"--to", "2024-01-15T10:02:30",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	got := string(data)
	if !strings.Contains(got, "DEBUG checking config") {
		t.Errorf("expected DEBUG line in output, got:\n%s", got)
	}
	if !strings.Contains(got, "ERROR connection refused") {
		t.Errorf("expected ERROR line in output, got:\n%s", got)
	}
	if strings.Contains(got, "service started") {
		t.Errorf("expected INFO line excluded, got:\n%s", got)
	}
}

func TestRun_PatternFilter(t *testing.T) {
	input := writeTempLog(t, sampleLog)
	output := filepath.Join(t.TempDir(), "out.log")

	err := cli.Run([]string{
		"--input", input,
		"--output", output,
		"--pattern", "ERROR|WARN",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(output)
	got := string(data)
	if strings.Contains(got, "INFO") || strings.Contains(got, "DEBUG") {
		t.Errorf("pattern filter should exclude INFO/DEBUG lines, got:\n%s", got)
	}
}
