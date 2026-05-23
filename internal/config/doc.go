// Package config provides loading, parsing, and merging of logslice
// configuration files.
//
// Supported formats are YAML (.yaml / .yml) and JSON (.json).
// A config file can supply defaults for all CLI flags; CLI flags always
// take precedence when merged with [Config.Merge].
//
// Typical usage:
//
//	cfg, err := config.Load(".logslice.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
//	// merge with flags parsed from os.Args
//	effective := cfg.Merge(flagConfig)
package config
