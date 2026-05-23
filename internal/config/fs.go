package config

import "os"

// exists reports whether a regular file exists at the given path.
func exists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
