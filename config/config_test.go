package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Unsetenv("NEWAPI_DB_TYPE")
	os.Unsetenv("NEWAPI_DB_PATH")
	os.Unsetenv("NEWAPI_OUTPUT_FORMAT")
	os.Unsetenv("NEWAPI_MYSQL_HOST")
	os.Unsetenv("NEWAPI_MYSQL_PORT")
	os.Unsetenv("NEWAPI_MYSQL_USER")
	os.Unsetenv("NEWAPI_MYSQL_PASSWORD")
	os.Unsetenv("NEWAPI_MYSQL_DATABASE")

	cfg := Load()

	if cfg.DBType != "sqlite" {
		t.Errorf("expected default DBType 'sqlite', got '%s'", cfg.DBType)
	}
	if cfg.DatabasePath != "newapi.db" {
		t.Errorf("expected default DatabasePath 'newapi.db', got '%s'", cfg.DatabasePath)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected default OutputFormat 'json', got '%s'", cfg.OutputFormat)
	}
	if cfg.MySQLHost != "localhost" {
		t.Errorf("expected default MySQLHost 'localhost', got '%s'", cfg.MySQLHost)
	}
	if cfg.MySQLPort != "3306" {
		t.Errorf("expected default MySQLPort '3306', got '%s'", cfg.MySQLPort)
	}
	if cfg.MySQLUser != "root" {
		t.Errorf("expected default MySQLUser 'root', got '%s'", cfg.MySQLUser)
	}
	if cfg.MySQLPassword != "" {
		t.Errorf("expected default MySQLPassword '', got '%s'", cfg.MySQLPassword)
	}
	if cfg.MySQLDatabase != "newapi" {
		t.Errorf("expected default MySQLDatabase 'newapi', got '%s'", cfg.MySQLDatabase)
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

func TestLoadMySQLEnv(t *testing.T) {
	os.Setenv("NEWAPI_DB_TYPE", "mysql")
	os.Setenv("NEWAPI_MYSQL_HOST", "10.0.0.1")
	os.Setenv("NEWAPI_MYSQL_PORT", "3307")
	os.Setenv("NEWAPI_MYSQL_USER", "admin")
	os.Setenv("NEWAPI_MYSQL_PASSWORD", "secret")
	os.Setenv("NEWAPI_MYSQL_DATABASE", "mydb")
	defer func() {
		os.Unsetenv("NEWAPI_DB_TYPE")
		os.Unsetenv("NEWAPI_MYSQL_HOST")
		os.Unsetenv("NEWAPI_MYSQL_PORT")
		os.Unsetenv("NEWAPI_MYSQL_USER")
		os.Unsetenv("NEWAPI_MYSQL_PASSWORD")
		os.Unsetenv("NEWAPI_MYSQL_DATABASE")
	}()

	cfg := Load()

	if cfg.DBType != "mysql" {
		t.Errorf("expected DBType 'mysql', got '%s'", cfg.DBType)
	}
	if cfg.MySQLHost != "10.0.0.1" {
		t.Errorf("expected MySQLHost '10.0.0.1', got '%s'", cfg.MySQLHost)
	}
	if cfg.MySQLPort != "3307" {
		t.Errorf("expected MySQLPort '3307', got '%s'", cfg.MySQLPort)
	}
	if cfg.MySQLUser != "admin" {
		t.Errorf("expected MySQLUser 'admin', got '%s'", cfg.MySQLUser)
	}
	if cfg.MySQLPassword != "secret" {
		t.Errorf("expected MySQLPassword 'secret', got '%s'", cfg.MySQLPassword)
	}
	if cfg.MySQLDatabase != "mydb" {
		t.Errorf("expected MySQLDatabase 'mydb', got '%s'", cfg.MySQLDatabase)
	}
}
