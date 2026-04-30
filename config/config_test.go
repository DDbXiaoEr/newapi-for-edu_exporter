package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Unsetenv("NEWAPI_DB_PATH")
	os.Unsetenv("NEWAPI_OUTPUT_FORMAT")

	cfg := Load()

	if cfg.DatabasePath != "newapi.db" {
		t.Errorf("expected default DatabasePath 'newapi.db', got '%s'", cfg.DatabasePath)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected default OutputFormat 'json', got '%s'", cfg.OutputFormat)
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("NEWAPI_DB_PATH", "/tmp/test.db")
	os.Setenv("NEWAPI_OUTPUT_FORMAT", "csv")
	defer func() {
		os.Unsetenv("NEWAPI_DB_PATH")
		os.Unsetenv("NEWAPI_OUTPUT_FORMAT")
	}()

	cfg := Load()

	if cfg.DatabasePath != "/tmp/test.db" {
		t.Errorf("expected DatabasePath '/tmp/test.db', got '%s'", cfg.DatabasePath)
	}
	if cfg.OutputFormat != "csv" {
		t.Errorf("expected OutputFormat 'csv', got '%s'", cfg.OutputFormat)
	}
}

func TestLoadPartialEnv(t *testing.T) {
	os.Setenv("NEWAPI_DB_PATH", "/tmp/partial.db")
	os.Unsetenv("NEWAPI_OUTPUT_FORMAT")
	defer os.Unsetenv("NEWAPI_DB_PATH")

	cfg := Load()

	if cfg.DatabasePath != "/tmp/partial.db" {
		t.Errorf("expected DatabasePath '/tmp/partial.db', got '%s'", cfg.DatabasePath)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected default OutputFormat 'json' when env not set, got '%s'", cfg.OutputFormat)
	}
}

func TestLoadWithDBPathEnv(t *testing.T) {
	os.Setenv("NEWAPI_DB_PATH", "my_data.db")
	defer os.Unsetenv("NEWAPI_DB_PATH")

	cfg := Load()

	if cfg.DatabasePath != "my_data.db" {
		t.Errorf("expected DatabasePath 'my_data.db', got '%s'", cfg.DatabasePath)
	}
}
