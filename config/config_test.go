package config

import (
	"os"
	"path/filepath"
	"testing"
)

func unsetConfigEnv(t *testing.T) {
	t.Helper()
	keys := []string{
		"NEWAPI_CONFIG_FILE",
		"NEWAPI_DB_TYPE",
		"NEWAPI_OUTPUT_FORMAT",
		"NEWAPI_POSTGRES_DSN",
		"SQL_DSN",
		"NEWAPI_POSTGRES_HOST",
		"NEWAPI_POSTGRES_PORT",
		"NEWAPI_POSTGRES_USER",
		"NEWAPI_POSTGRES_PASSWORD",
		"NEWAPI_POSTGRES_DATABASE",
		"NEWAPI_POSTGRES_SSLMODE",
		"NEWAPI_MYSQL_HOST",
		"NEWAPI_MYSQL_PORT",
		"NEWAPI_MYSQL_USER",
		"NEWAPI_MYSQL_PASSWORD",
		"NEWAPI_MYSQL_DATABASE",
		"NEWAPI_CLICKHOUSE_DSN",
		"LOG_SQL_DSN",
		"NEWAPI_CLICKHOUSE_TABLE",
		"LOG_SQL_CLICKHOUSE_TABLE",
		"NEWAPI_DB_PATH",
	}
	for _, key := range keys {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
}

func TestLoadDefaults(t *testing.T) {
	unsetConfigEnv(t)
	t.Chdir(t.TempDir())

	cfg := Load()

	if cfg.DBType != "postgres" {
		t.Errorf("expected default DBType 'postgres', got '%s'", cfg.DBType)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected default OutputFormat 'json', got '%s'", cfg.OutputFormat)
	}
	if cfg.PostgresHost != "localhost" {
		t.Errorf("expected default PostgresHost 'localhost', got '%s'", cfg.PostgresHost)
	}
	if cfg.PostgresPort != "5432" {
		t.Errorf("expected default PostgresPort '5432', got '%s'", cfg.PostgresPort)
	}
	if cfg.PostgresUser != "root" {
		t.Errorf("expected default PostgresUser 'root', got '%s'", cfg.PostgresUser)
	}
	if cfg.PostgresDatabase != "newapi" {
		t.Errorf("expected default PostgresDatabase 'newapi', got '%s'", cfg.PostgresDatabase)
	}
	if cfg.PostgresSSLMode != "disable" {
		t.Errorf("expected default PostgresSSLMode 'disable', got '%s'", cfg.PostgresSSLMode)
	}
	if cfg.ClickHouseTable != "logs" {
		t.Errorf("expected default ClickHouseTable 'logs', got '%s'", cfg.ClickHouseTable)
	}
	if cfg.ClickHouseDSN != "" {
		t.Errorf("expected default ClickHouseDSN empty, got '%s'", cfg.ClickHouseDSN)
	}
}

func TestLoadFromEnv(t *testing.T) {
	unsetConfigEnv(t)
	t.Chdir(t.TempDir())
	t.Setenv("NEWAPI_OUTPUT_FORMAT", "csv")
	t.Setenv("NEWAPI_DB_TYPE", "mysql")

	cfg := Load()

	if cfg.OutputFormat != "csv" {
		t.Errorf("expected OutputFormat 'csv', got '%s'", cfg.OutputFormat)
	}
	if cfg.DBType != "mysql" {
		t.Errorf("expected DBType 'mysql', got '%s'", cfg.DBType)
	}
}

func TestLoadPartialEnv(t *testing.T) {
	unsetConfigEnv(t)
	t.Chdir(t.TempDir())
	t.Setenv("NEWAPI_POSTGRES_HOST", "pg.internal")

	cfg := Load()

	if cfg.PostgresHost != "pg.internal" {
		t.Errorf("expected PostgresHost 'pg.internal', got '%s'", cfg.PostgresHost)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected default OutputFormat 'json' when env not set, got '%s'", cfg.OutputFormat)
	}
}

func TestLoadPostgresEnv(t *testing.T) {
	unsetConfigEnv(t)
	t.Chdir(t.TempDir())
	t.Setenv("NEWAPI_DB_TYPE", "postgres")
	t.Setenv("NEWAPI_POSTGRES_HOST", "10.0.0.1")
	t.Setenv("NEWAPI_POSTGRES_PORT", "5433")
	t.Setenv("NEWAPI_POSTGRES_USER", "admin")
	t.Setenv("NEWAPI_POSTGRES_PASSWORD", "secret")
	t.Setenv("NEWAPI_POSTGRES_DATABASE", "mydb")
	t.Setenv("NEWAPI_POSTGRES_SSLMODE", "require")

	cfg := Load()

	if cfg.DBType != "postgres" {
		t.Errorf("expected DBType 'postgres', got '%s'", cfg.DBType)
	}
	if cfg.PostgresHost != "10.0.0.1" {
		t.Errorf("expected PostgresHost '10.0.0.1', got '%s'", cfg.PostgresHost)
	}
	if cfg.PostgresPort != "5433" {
		t.Errorf("expected PostgresPort '5433', got '%s'", cfg.PostgresPort)
	}
	if cfg.PostgresUser != "admin" {
		t.Errorf("expected PostgresUser 'admin', got '%s'", cfg.PostgresUser)
	}
	if cfg.PostgresPassword != "secret" {
		t.Errorf("expected PostgresPassword 'secret', got '%s'", cfg.PostgresPassword)
	}
	if cfg.PostgresDatabase != "mydb" {
		t.Errorf("expected PostgresDatabase 'mydb', got '%s'", cfg.PostgresDatabase)
	}
	if cfg.PostgresSSLMode != "require" {
		t.Errorf("expected PostgresSSLMode 'require', got '%s'", cfg.PostgresSSLMode)
	}

	dsn := cfg.PostgresDSNString()
	if dsn != "postgres://admin:secret@10.0.0.1:5433/mydb?sslmode=require" {
		t.Errorf("unexpected PostgresDSNString: %s", dsn)
	}
}

func TestLoadPostgresDSNEnv(t *testing.T) {
	unsetConfigEnv(t)
	t.Chdir(t.TempDir())
	t.Setenv("SQL_DSN", "postgresql://root:123456@localhost:5432/newapi?sslmode=disable")

	cfg := Load()

	if cfg.PostgresDSN != "postgresql://root:123456@localhost:5432/newapi?sslmode=disable" {
		t.Errorf("expected SQL_DSN to populate PostgresDSN, got '%s'", cfg.PostgresDSN)
	}
	if cfg.PostgresDSNString() != cfg.PostgresDSN {
		t.Errorf("expected PostgresDSNString to return explicit DSN")
	}
}

func TestLoadClickHouseEnv(t *testing.T) {
	unsetConfigEnv(t)
	t.Chdir(t.TempDir())
	t.Setenv("LOG_SQL_DSN", "clickhouse://default:@localhost:9000/newapi_audit")
	t.Setenv("LOG_SQL_CLICKHOUSE_TABLE", "token_usage_logs")

	cfg := Load()

	if cfg.ClickHouseDSN != "clickhouse://default:@localhost:9000/newapi_audit" {
		t.Errorf("expected ClickHouseDSN from LOG_SQL_DSN, got '%s'", cfg.ClickHouseDSN)
	}
	if cfg.ClickHouseTable != "token_usage_logs" {
		t.Errorf("expected ClickHouseTable 'token_usage_logs', got '%s'", cfg.ClickHouseTable)
	}
}

func TestLoadMySQLEnv(t *testing.T) {
	unsetConfigEnv(t)
	t.Chdir(t.TempDir())
	t.Setenv("NEWAPI_DB_TYPE", "mysql")
	t.Setenv("NEWAPI_MYSQL_HOST", "10.0.0.1")
	t.Setenv("NEWAPI_MYSQL_PORT", "3307")
	t.Setenv("NEWAPI_MYSQL_USER", "admin")
	t.Setenv("NEWAPI_MYSQL_PASSWORD", "secret")
	t.Setenv("NEWAPI_MYSQL_DATABASE", "mydb")

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
	unsetConfigEnv(t)

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `
db_type: postgres
output_format: csv
postgres_host: 192.168.1.1
postgres_port: "5433"
postgres_user: fileuser
postgres_password: filepass
postgres_database: filedb
postgres_sslmode: require
clickhouse_dsn: clickhouse://default:@localhost:9000/newapi_audit
clickhouse_table: logs
channel_names:
  1: OpenAI
  2: Claude
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := LoadWithConfigFile(path)

	if cfg.DBType != "postgres" {
		t.Errorf("expected DBType 'postgres', got '%s'", cfg.DBType)
	}
	if cfg.OutputFormat != "csv" {
		t.Errorf("expected OutputFormat 'csv', got '%s'", cfg.OutputFormat)
	}
	if cfg.PostgresHost != "192.168.1.1" {
		t.Errorf("expected PostgresHost '192.168.1.1', got '%s'", cfg.PostgresHost)
	}
	if cfg.PostgresPort != "5433" {
		t.Errorf("expected PostgresPort '5433', got '%s'", cfg.PostgresPort)
	}
	if cfg.PostgresUser != "fileuser" {
		t.Errorf("expected PostgresUser 'fileuser', got '%s'", cfg.PostgresUser)
	}
	if cfg.PostgresPassword != "filepass" {
		t.Errorf("expected PostgresPassword 'filepass', got '%s'", cfg.PostgresPassword)
	}
	if cfg.PostgresDatabase != "filedb" {
		t.Errorf("expected PostgresDatabase 'filedb', got '%s'", cfg.PostgresDatabase)
	}
	if cfg.ClickHouseDSN != "clickhouse://default:@localhost:9000/newapi_audit" {
		t.Errorf("expected ClickHouseDSN from file, got '%s'", cfg.ClickHouseDSN)
	}
	if cfg.ChannelNames["1"] != "OpenAI" || cfg.ChannelNames["2"] != "Claude" {
		t.Errorf("expected channel_names from file, got %#v", cfg.ChannelNames)
	}
}

func TestConfigFileEnvOverride(t *testing.T) {
	unsetConfigEnv(t)

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `
db_type: mysql
output_format: json
postgres_host: 10.0.0.10
clickhouse_dsn: clickhouse://file@localhost:9000/filedb
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("NEWAPI_DB_TYPE", "postgres")
	t.Setenv("NEWAPI_POSTGRES_HOST", "env-override.com")
	t.Setenv("LOG_SQL_DSN", "clickhouse://env@localhost:9000/envdb")

	cfg := LoadWithConfigFile(path)

	if cfg.DBType != "postgres" {
		t.Errorf("expected env override DBType 'postgres', got '%s'", cfg.DBType)
	}
	if cfg.PostgresHost != "env-override.com" {
		t.Errorf("expected env override PostgresHost 'env-override.com', got '%s'", cfg.PostgresHost)
	}
	if cfg.ClickHouseDSN != "clickhouse://env@localhost:9000/envdb" {
		t.Errorf("expected env override ClickHouseDSN, got '%s'", cfg.ClickHouseDSN)
	}
}

func TestConfigFilePartial(t *testing.T) {
	unsetConfigEnv(t)

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
	if cfg.PostgresPort != "5432" {
		t.Errorf("expected default PostgresPort '5432', got '%s'", cfg.PostgresPort)
	}
	if cfg.MySQLPort != "3306" {
		t.Errorf("expected default MySQLPort '3306', got '%s'", cfg.MySQLPort)
	}
}

func TestLoadFromCwdConfigFile(t *testing.T) {
	unsetConfigEnv(t)

	dir := t.TempDir()
	t.Chdir(dir)

	content := `
db_type: mysql
mysql_host: cwd.db.host
output_format: csv
`
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Load()

	if cfg.DBType != "mysql" {
		t.Errorf("expected DBType 'mysql' from cwd config.yaml, got '%s'", cfg.DBType)
	}
	if cfg.MySQLHost != "cwd.db.host" {
		t.Errorf("expected MySQLHost 'cwd.db.host' from cwd config.yaml, got '%s'", cfg.MySQLHost)
	}
	if cfg.OutputFormat != "csv" {
		t.Errorf("expected OutputFormat 'csv' from cwd config.yaml, got '%s'", cfg.OutputFormat)
	}
}

func TestLoadFromCwdConfigYmlFallback(t *testing.T) {
	unsetConfigEnv(t)

	dir := t.TempDir()
	t.Chdir(dir)

	content := `
db_type: mysql
mysql_host: yml.db.host
`
	if err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := LoadWithConfigFile("config.yaml")

	if cfg.MySQLHost != "yml.db.host" {
		t.Errorf("expected MySQLHost 'yml.db.host' from cwd config.yml fallback, got '%s'", cfg.MySQLHost)
	}
}

func TestResolveConfigPath(t *testing.T) {
	unsetConfigEnv(t)
	dir := t.TempDir()
	t.Chdir(dir)

	if p := resolveConfigPath(""); p != "" {
		t.Errorf("expected empty path when no config in cwd, got '%s'", p)
	}
	if p := resolveConfigPath("/nonexistent/path.yaml"); p != "/nonexistent/path.yaml" {
		t.Errorf("expected explicit path returned, got '%s'", p)
	}

	yamlPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(yamlPath, []byte("db_type: postgres\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if p := resolveConfigPath(""); p != yamlPath {
		t.Errorf("expected cwd config.yaml, got '%s'", p)
	}
	if p := resolveConfigPath("config.yaml"); p != yamlPath {
		t.Errorf("expected cwd-relative config.yaml, got '%s'", p)
	}
}

func TestLoadFromFileNotFound(t *testing.T) {
	_, err := loadFromFile("/nonexistent/file.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
