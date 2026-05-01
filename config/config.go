package config

import (
	"os"
)

type Config struct {
	DBType        string
	DatabasePath  string
	OutputFormat  string
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string
}

func Load() *Config {
	dbType := os.Getenv("NEWAPI_DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	dbPath := os.Getenv("NEWAPI_DB_PATH")
	if dbPath == "" {
		dbPath = "newapi.db"
	}

	outputFormat := os.Getenv("NEWAPI_OUTPUT_FORMAT")
	if outputFormat == "" {
		outputFormat = "json"
	}

	mysqlHost := os.Getenv("NEWAPI_MYSQL_HOST")
	if mysqlHost == "" {
		mysqlHost = "localhost"
	}

	mysqlPort := os.Getenv("NEWAPI_MYSQL_PORT")
	if mysqlPort == "" {
		mysqlPort = "3306"
	}

	mysqlUser := os.Getenv("NEWAPI_MYSQL_USER")
	if mysqlUser == "" {
		mysqlUser = "root"
	}

	mysqlPassword := os.Getenv("NEWAPI_MYSQL_PASSWORD")
	if mysqlPassword == "" {
		mysqlPassword = ""
	}

	mysqlDatabase := os.Getenv("NEWAPI_MYSQL_DATABASE")
	if mysqlDatabase == "" {
		mysqlDatabase = "newapi"
	}

	return &Config{
		DBType:        dbType,
		DatabasePath:  dbPath,
		OutputFormat:  outputFormat,
		MySQLHost:     mysqlHost,
		MySQLPort:     mysqlPort,
		MySQLUser:     mysqlUser,
		MySQLPassword: mysqlPassword,
		MySQLDatabase: mysqlDatabase,
	}
}
