package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/logslice/internal/config"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestLoad_YAML(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "cfg.yaml", `
from: "2024-01-01T00:00:00Z"
to: "2024-01-02T00:00:00Z"
level: error
patterns:
  - "timeout"
output_dir: /tmp/out
strict: true
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.From != "2024-01-01T00:00:00Z" {
		t.Errorf("From = %q", cfg.From)
	}
	if cfg.Level != "error" {
		t.Errorf("Level = %q", cfg.Level)
	}
	if !cfg.Strict {
		t.Error("expected Strict=true")
	}
	if len(cfg.Patterns) != 1 || cfg.Patterns[0] != "timeout" {
		t.Errorf("Patterns = %v", cfg.Patterns)
	}
}

func TestLoad_JSON(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "cfg.json", `{"from":"2024-03-01T00:00:00Z","level":"warn"}`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Level != "warn" {
		t.Errorf("Level = %q", cfg.Level)
	}
}

func TestLoad_UnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "cfg.toml", `from = "2024-01-01"`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for unsupported extension")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/cfg.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestConfig_Merge(t *testing.T) {
	base := &config.Config{Level: "info", Strict: false}
	override := &config.Config{Level: "error", Strict: true, OutputDir: "/out"}
	merged := base.Merge(override)
	if merged.Level != "error" {
		t.Errorf("Level = %q, want error", merged.Level)
	}
	if !merged.Strict {
		t.Error("expected Strict=true after merge")
	}
	if merged.OutputDir != "/out" {
		t.Errorf("OutputDir = %q", merged.OutputDir)
	}
	// base should be unchanged
	if base.Level != "info" {
		t.Error("base config was mutated")
	}
}
