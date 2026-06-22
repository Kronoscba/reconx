package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Debug      bool
	OutputDir  string
	Target     string
}

func Load(target string) *Config {
	outputDir := "output"
	if target != "" {
		outputDir = filepath.Join("output", target)
	}

	// Ensure output directory exists
	os.MkdirAll(outputDir, 0755)

	return &Config{
		Debug:     os.Getenv("RECONX_DEBUG") == "true",
		OutputDir: outputDir,
		Target:    target,
	}
}

func (c *Config) GetPath(filename string) string {
	return filepath.Join(c.OutputDir, filename)
}
