package config

import (
	"os"
)

type Config struct {
	DatabasePath string
	OutputFormat string
}

func Load() *Config {
	dbPath := os.Getenv("NEWAPI_DB_PATH")
	if dbPath == "" {
		dbPath = "newapi.db"
	}

	outputFormat := os.Getenv("NEWAPI_OUTPUT_FORMAT")
	if outputFormat == "" {
		outputFormat = "json"
	}

	return &Config{
		DatabasePath: dbPath,
		OutputFormat: outputFormat,
	}
}
