package config

// DefaultConfigFilenames lists the file names logslice probes for
// automatically when no --config flag is provided.
var DefaultConfigFilenames = []string{
	".logslice.yaml",
	".logslice.yml",
	".logslice.json",
	"logslice.yaml",
	"logslice.yml",
	"logslice.json",
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Strict:     false,
		Timestamps: false,
		OutputDir:  "",
	}
}

// FindDefault searches the current directory for any of the
// DefaultConfigFilenames and returns the path to the first one found.
// If none is found it returns ("", nil).
func FindDefault() (string, error) {
	for _, name := range DefaultConfigFilenames {
		if exists(name) {
			return name, nil
		}
	}
	return "", nil
}

// LoadOrDefault attempts to load path; if path is empty it calls FindDefault.
// If no config file is found it returns Default().
func LoadOrDefault(path string) (*Config, error) {
	if path == "" {
		var err error
		path, err = FindDefault()
		if err != nil {
			return nil, err
		}
	}
	if path == "" {
		return Default(), nil
	}
	return Load(path)
}
