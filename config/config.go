package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBType        string `yaml:"db_type"`
	DatabasePath  string `yaml:"database_path"`
	OutputFormat  string `yaml:"output_format"`
	MySQLHost     string `yaml:"mysql_host"`
	MySQLPort     string `yaml:"mysql_port"`
	MySQLUser     string `yaml:"mysql_user"`
	MySQLPassword string `yaml:"mysql_password"`
	MySQLDatabase string `yaml:"mysql_database"`
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
		DBType:        "sqlite",
		DatabasePath:  "newapi.db",
		OutputFormat:  "json",
		MySQLHost:     "localhost",
		MySQLPort:     "3306",
		MySQLUser:     "root",
		MySQLPassword: "",
		MySQLDatabase: "newapi",
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
	if fileCfg.DatabasePath != "" {
		c.DatabasePath = fileCfg.DatabasePath
	}
	if fileCfg.OutputFormat != "" {
		c.OutputFormat = fileCfg.OutputFormat
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
}

func (c *Config) mergeEnv() {
	if v := os.Getenv("NEWAPI_DB_TYPE"); v != "" {
		c.DBType = v
	}
	if v := os.Getenv("NEWAPI_DB_PATH"); v != "" {
		c.DatabasePath = v
	}
	if v := os.Getenv("NEWAPI_OUTPUT_FORMAT"); v != "" {
		c.OutputFormat = v
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
}