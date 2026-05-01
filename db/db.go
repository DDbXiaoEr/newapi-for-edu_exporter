package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/quanxiaoer/newapi_exporter/config"
)

type Database struct {
	conn *sql.DB
}

func Open(cfg *config.Config) (*Database, error) {
	var driverName, dataSourceName string

	switch cfg.DBType {
	case "sqlite":
		driverName = "sqlite3"
		dataSourceName = cfg.DatabasePath
	case "mysql":
		driverName = "mysql"
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
			cfg.MySQLUser, cfg.MySQLPassword, cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDatabase)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Database{conn: conn}, nil
}

func (d *Database) Close() error {
	return d.conn.Close()
}

func (d *Database) Conn() *sql.DB {
	return d.conn
}
