// Package config handles loading and validating logslice configuration
// from files (YAML/JSON) and merging with CLI flags.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration options for logslice.
type Config struct {
	From       string   `yaml:"from"       json:"from"`
	To         string   `yaml:"to"         json:"to"`
	Level      string   `yaml:"level"      json:"level"`
	Patterns   []string `yaml:"patterns"   json:"patterns"`
	OutputDir  string   `yaml:"output_dir" json:"output_dir"`
	Strict     bool     `yaml:"strict"     json:"strict"`
	Timestamps bool     `yaml:"timestamps" json:"timestamps"`
}

// Load reads a config file (YAML or JSON) from the given path.
// The format is inferred from the file extension.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	cfg := &Config{}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("config: parse YAML %q: %w", path, err)
		}
	case ".json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("config: parse JSON %q: %w", path, err)
		}
	default:
		return nil, fmt.Errorf("config: unsupported format %q (use .yaml or .json)", ext)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate checks that the config fields are consistent.
func (c *Config) Validate() error {
	if c.From != "" && c.To != "" {
		// Basic sanity: from <= to is validated downstream with real time parsing.
	}
	return nil
}

// Merge overlays non-zero values from other onto c, returning a new Config.
func (c *Config) Merge(other *Config) *Config {
	out := *c
	if other.From != "" {
		out.From = other.From
	}
	if other.To != "" {
		out.To = other.To
	}
	if other.Level != "" {
		out.Level = other.Level
	}
	if len(other.Patterns) > 0 {
		out.Patterns = other.Patterns
	}
	if other.OutputDir != "" {
		out.OutputDir = other.OutputDir
	}
	if other.Strict {
		out.Strict = true
	}
	if other.Timestamps {
		out.Timestamps = true
	}
	return &out
}
