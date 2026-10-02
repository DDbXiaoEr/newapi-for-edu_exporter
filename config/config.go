package config

import (
	"fmt"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBType       string `yaml:"db_type"`
	OutputFormat string `yaml:"output_format"`

	PostgresDSN      string `yaml:"postgres_dsn"`
	PostgresHost     string `yaml:"postgres_host"`
	PostgresPort     string `yaml:"postgres_port"`
	PostgresUser     string `yaml:"postgres_user"`
	PostgresPassword string `yaml:"postgres_password"`
	PostgresDatabase string `yaml:"postgres_database"`
	PostgresSSLMode  string `yaml:"postgres_sslmode"`

	MySQLHost     string `yaml:"mysql_host"`
	MySQLPort     string `yaml:"mysql_port"`
	MySQLUser     string `yaml:"mysql_user"`
	MySQLPassword string `yaml:"mysql_password"`
	MySQLDatabase string `yaml:"mysql_database"`

	ClickHouseDSN   string `yaml:"clickhouse_dsn"`
	ClickHouseTable string `yaml:"clickhouse_table"`
}

func Load() *Config {
	return LoadWithConfigFile("")
}

func LoadWithConfigFile(configPath string) *Config {
	cfg := defaults()

	if path := resolveConfigPath(configPath); path != "" {
		if fileCfg, err := loadFromFile(path); err == nil {
			cfg.mergeFile(fileCfg)
		}
	}

	cfg.mergeEnv()

	return cfg
}

func resolveConfigPath(explicitPath string) string {
	if explicitPath != "" {
		return explicitPath
	}
	if p := os.Getenv("NEWAPI_CONFIG_FILE"); p != "" {
		return p
	}
	for _, name := range []string{"config.yaml", "config.yml"} {
		if _, err := os.Stat(name); err == nil {
			return name
		}
	}
	return ""
}

func defaults() *Config {
	return &Config{
		DBType:           "postgres",
		OutputFormat:     "json",
		PostgresHost:     "localhost",
		PostgresPort:     "5432",
		PostgresUser:     "root",
		PostgresPassword: "",
		PostgresDatabase: "newapi",
		PostgresSSLMode:  "disable",
		MySQLHost:        "localhost",
		MySQLPort:        "3306",
		MySQLUser:        "root",
		MySQLPassword:    "",
		MySQLDatabase:    "newapi",
		ClickHouseTable:  "logs",
	}
}

func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return &cfg, nil
}

func (c *Config) mergeFile(fileCfg *Config) {
	if fileCfg.DBType != "" {
		c.DBType = fileCfg.DBType
	}
	if fileCfg.OutputFormat != "" {
		c.OutputFormat = fileCfg.OutputFormat
	}
	if fileCfg.PostgresDSN != "" {
		c.PostgresDSN = fileCfg.PostgresDSN
	}
	if fileCfg.PostgresHost != "" {
		c.PostgresHost = fileCfg.PostgresHost
	}
	if fileCfg.PostgresPort != "" {
		c.PostgresPort = fileCfg.PostgresPort
	}
	if fileCfg.PostgresUser != "" {
		c.PostgresUser = fileCfg.PostgresUser
	}
	if fileCfg.PostgresPassword != "" {
		c.PostgresPassword = fileCfg.PostgresPassword
	}
	if fileCfg.PostgresDatabase != "" {
		c.PostgresDatabase = fileCfg.PostgresDatabase
	}
	if fileCfg.PostgresSSLMode != "" {
		c.PostgresSSLMode = fileCfg.PostgresSSLMode
	}
	if fileCfg.MySQLHost != "" {
		c.MySQLHost = fileCfg.MySQLHost
	}
	if fileCfg.MySQLPort != "" {
		c.MySQLPort = fileCfg.MySQLPort
	}
	if fileCfg.MySQLUser != "" {
		c.MySQLUser = fileCfg.MySQLUser
	}
	if fileCfg.MySQLPassword != "" {
		c.MySQLPassword = fileCfg.MySQLPassword
	}
	if fileCfg.MySQLDatabase != "" {
		c.MySQLDatabase = fileCfg.MySQLDatabase
	}
	if fileCfg.ClickHouseDSN != "" {
		c.ClickHouseDSN = fileCfg.ClickHouseDSN
	}
	if fileCfg.ClickHouseTable != "" {
		c.ClickHouseTable = fileCfg.ClickHouseTable
	}
}

func (c *Config) mergeEnv() {
	if v := os.Getenv("NEWAPI_DB_TYPE"); v != "" {
		c.DBType = v
	}
	if v := os.Getenv("NEWAPI_OUTPUT_FORMAT"); v != "" {
		c.OutputFormat = v
	}
	if v := firstEnv("NEWAPI_POSTGRES_DSN", "SQL_DSN"); v != "" {
		c.PostgresDSN = v
	}
	if v := os.Getenv("NEWAPI_POSTGRES_HOST"); v != "" {
		c.PostgresHost = v
	}
	if v := os.Getenv("NEWAPI_POSTGRES_PORT"); v != "" {
		c.PostgresPort = v
	}
	if v := os.Getenv("NEWAPI_POSTGRES_USER"); v != "" {
		c.PostgresUser = v
	}
	if v := os.Getenv("NEWAPI_POSTGRES_PASSWORD"); v != "" {
		c.PostgresPassword = v
	}
	if v := os.Getenv("NEWAPI_POSTGRES_DATABASE"); v != "" {
		c.PostgresDatabase = v
	}
	if v := os.Getenv("NEWAPI_POSTGRES_SSLMODE"); v != "" {
		c.PostgresSSLMode = v
	}
	if v := os.Getenv("NEWAPI_MYSQL_HOST"); v != "" {
		c.MySQLHost = v
	}
	if v := os.Getenv("NEWAPI_MYSQL_PORT"); v != "" {
		c.MySQLPort = v
	}
	if v := os.Getenv("NEWAPI_MYSQL_USER"); v != "" {
		c.MySQLUser = v
	}
	if v := os.Getenv("NEWAPI_MYSQL_PASSWORD"); v != "" {
		c.MySQLPassword = v
	}
	if v := os.Getenv("NEWAPI_MYSQL_DATABASE"); v != "" {
		c.MySQLDatabase = v
	}
	if v := firstEnv("NEWAPI_CLICKHOUSE_DSN", "LOG_SQL_DSN"); v != "" {
		c.ClickHouseDSN = v
	}
	if v := firstEnv("NEWAPI_CLICKHOUSE_TABLE", "LOG_SQL_CLICKHOUSE_TABLE"); v != "" {
		c.ClickHouseTable = v
	}
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return ""
}

func (c *Config) PostgresDSNString() string {
	if c.PostgresDSN != "" {
		return c.PostgresDSN
	}
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.PostgresUser, c.PostgresPassword),
		Host:   fmt.Sprintf("%s:%s", c.PostgresHost, c.PostgresPort),
		Path:   "/" + c.PostgresDatabase,
	}
	q := url.Values{}
	if c.PostgresSSLMode != "" {
		q.Set("sslmode", c.PostgresSSLMode)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Config) MySQLDSNString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		c.MySQLUser, c.MySQLPassword, c.MySQLHost, c.MySQLPort, c.MySQLDatabase)
}
