package config

import (
	"os"
	"path/filepath"
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
	os.Unsetenv("NEWAPI_CONFIG_FILE")

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

func TestLoadFromConfigFile(t *testing.T) {
	os.Unsetenv("NEWAPI_DB_TYPE")
	os.Unsetenv("NEWAPI_DB_PATH")
	os.Unsetenv("NEWAPI_OUTPUT_FORMAT")
	os.Unsetenv("NEWAPI_MYSQL_HOST")
	os.Unsetenv("NEWAPI_MYSQL_PORT")
	os.Unsetenv("NEWAPI_MYSQL_USER")
	os.Unsetenv("NEWAPI_MYSQL_PASSWORD")
	os.Unsetenv("NEWAPI_MYSQL_DATABASE")
	os.Unsetenv("NEWAPI_CONFIG_FILE")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `
db_type: mysql
database_path: /tmp/fromfile.db
output_format: csv
mysql_host: 192.168.1.1
mysql_port: "3307"
mysql_user: fileuser
mysql_password: filepass
mysql_database: filedb
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := LoadWithConfigFile(path)

	if cfg.DBType != "mysql" {
		t.Errorf("expected DBType 'mysql', got '%s'", cfg.DBType)
	}
	if cfg.DatabasePath != "/tmp/fromfile.db" {
		t.Errorf("expected DatabasePath '/tmp/fromfile.db', got '%s'", cfg.DatabasePath)
	}
	if cfg.OutputFormat != "csv" {
		t.Errorf("expected OutputFormat 'csv', got '%s'", cfg.OutputFormat)
	}
	if cfg.MySQLHost != "192.168.1.1" {
		t.Errorf("expected MySQLHost '192.168.1.1', got '%s'", cfg.MySQLHost)
	}
	if cfg.MySQLPort != "3307" {
		t.Errorf("expected MySQLPort '3307', got '%s'", cfg.MySQLPort)
	}
	if cfg.MySQLUser != "fileuser" {
		t.Errorf("expected MySQLUser 'fileuser', got '%s'", cfg.MySQLUser)
	}
	if cfg.MySQLPassword != "filepass" {
		t.Errorf("expected MySQLPassword 'filepass', got '%s'", cfg.MySQLPassword)
	}
	if cfg.MySQLDatabase != "filedb" {
		t.Errorf("expected MySQLDatabase 'filedb', got '%s'", cfg.MySQLDatabase)
	}
}

func TestConfigFileEnvOverride(t *testing.T) {
	os.Unsetenv("NEWAPI_DB_TYPE")
	os.Unsetenv("NEWAPI_DB_PATH")
	os.Unsetenv("NEWAPI_OUTPUT_FORMAT")
	os.Unsetenv("NEWAPI_MYSQL_HOST")
	os.Unsetenv("NEWAPI_MYSQL_PORT")
	os.Unsetenv("NEWAPI_MYSQL_USER")
	os.Unsetenv("NEWAPI_MYSQL_PASSWORD")
	os.Unsetenv("NEWAPI_MYSQL_DATABASE")
	os.Unsetenv("NEWAPI_CONFIG_FILE")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `
db_type: sqlite
database_path: /tmp/fromfile.db
output_format: json
mysql_host: 10.0.0.10
mysql_port: "3306"
mysql_user: fileuser
mysql_password: filepass
mysql_database: filedb
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	os.Setenv("NEWAPI_DB_TYPE", "mysql")
	os.Setenv("NEWAPI_MYSQL_HOST", "env-override.com")
	defer func() {
		os.Unsetenv("NEWAPI_DB_TYPE")
		os.Unsetenv("NEWAPI_MYSQL_HOST")
	}()

	cfg := LoadWithConfigFile(path)

	if cfg.DBType != "mysql" {
		t.Errorf("expected env override DBType 'mysql', got '%s'", cfg.DBType)
	}
	if cfg.MySQLHost != "env-override.com" {
		t.Errorf("expected env override MySQLHost 'env-override.com', got '%s'", cfg.MySQLHost)
	}
	if cfg.DatabasePath != "/tmp/fromfile.db" {
		t.Errorf("expected file DatabasePath '/tmp/fromfile.db', got '%s'", cfg.DatabasePath)
	}
}

func TestConfigFilePartial(t *testing.T) {
	os.Unsetenv("NEWAPI_DB_TYPE")
	os.Unsetenv("NEWAPI_DB_PATH")
	os.Unsetenv("NEWAPI_OUTPUT_FORMAT")
	os.Unsetenv("NEWAPI_MYSQL_HOST")
	os.Unsetenv("NEWAPI_MYSQL_PORT")
	os.Unsetenv("NEWAPI_MYSQL_USER")
	os.Unsetenv("NEWAPI_MYSQL_PASSWORD")
	os.Unsetenv("NEWAPI_MYSQL_DATABASE")
	os.Unsetenv("NEWAPI_CONFIG_FILE")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `
db_type: mysql
mysql_host: partial.db.host
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := LoadWithConfigFile(path)

	if cfg.DBType != "mysql" {
		t.Errorf("expected DBType 'mysql', got '%s'", cfg.DBType)
	}
	if cfg.MySQLHost != "partial.db.host" {
		t.Errorf("expected MySQLHost 'partial.db.host', got '%s'", cfg.MySQLHost)
	}
	if cfg.DatabasePath != "newapi.db" {
		t.Errorf("expected default DatabasePath 'newapi.db', got '%s'", cfg.DatabasePath)
	}
	if cfg.MySQLPort != "3306" {
		t.Errorf("expected default MySQLPort '3306', got '%s'", cfg.MySQLPort)
	}
}

func TestResolveConfigPath(t *testing.T) {
	if p := resolveConfigPath(""); p != "" {
		t.Logf("found config file: %s", p)
	}
	if p := resolveConfigPath("/nonexistent/path.yaml"); p != "/nonexistent/path.yaml" {
		t.Errorf("expected explicit path returned, got '%s'", p)
	}
}

func TestLoadFromFileNotFound(t *testing.T) {
	_, err := loadFromFile("/nonexistent/file.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}