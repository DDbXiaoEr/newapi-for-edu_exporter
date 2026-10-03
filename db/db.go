package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/quanxiaoer/newapi_exporter/config"
	"github.com/quanxiaoer/newapi_exporter/models"
)

const (
	DriverPostgres   = "postgres"
	DriverMySQL      = "mysql"
	DriverClickHouse = "clickhouse"
)

type Database struct {
	conn         *sql.DB
	logConn      *sql.DB
	driver       string
	logTable     string
	channelNames map[int64]string
}

func Open(cfg *config.Config) (*Database, error) {
	driver, dsn, err := primaryDSN(cfg)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open(sqlDriverName(driver), dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db := &Database{
		conn:         conn,
		driver:       driver,
		logTable:     cfg.ClickHouseTable,
		channelNames: parseChannelNames(cfg.ChannelNames),
	}
	if db.logTable == "" {
		db.logTable = "logs"
	}

	if dsn := strings.TrimSpace(cfg.ClickHouseDSN); isClickHouseDSN(dsn) {
		if !isValidIdent(db.logTable) {
			conn.Close()
			return nil, fmt.Errorf("invalid clickhouse table name %q", db.logTable)
		}
		logConn, err := sql.Open(DriverClickHouse, dsn)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("open clickhouse: %w", err)
		}
		if err := logConn.Ping(); err != nil {
			logConn.Close()
			conn.Close()
			return nil, fmt.Errorf("ping clickhouse: %w", err)
		}
		db.logConn = logConn
	}

	return db, nil
}

func isClickHouseDSN(dsn string) bool {
	return strings.HasPrefix(dsn, "clickhouse://") ||
		strings.HasPrefix(dsn, "tcp://") ||
		strings.HasPrefix(dsn, "http://") ||
		strings.HasPrefix(dsn, "https://")
}

func isValidIdent(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r == '_':
		case r >= '0' && r <= '9' && i > 0:
		default:
			return false
		}
	}
	return true
}

func primaryDSN(cfg *config.Config) (string, string, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.DBType)) {
	case "postgres", "postgresql", "pgsql":
		return DriverPostgres, cfg.PostgresDSNString(), nil
	case "mysql":
		return DriverMySQL, cfg.MySQLDSNString(), nil
	case "sqlite", "sqlite3":
		return "", "", fmt.Errorf("sqlite is no longer supported; use postgres or mysql")
	default:
		return "", "", fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}
}

func sqlDriverName(driver string) string {
	if driver == DriverPostgres {
		return "pgx"
	}
	return driver
}

func (d *Database) Close() error {
	var first error
	if d.logConn != nil {
		first = d.logConn.Close()
	}
	if err := d.conn.Close(); err != nil && first == nil {
		first = err
	}
	return first
}

func (d *Database) Conn() *sql.DB {
	return d.conn
}

func (d *Database) LogConn() *sql.DB {
	return d.logConn
}

func (d *Database) HasLogs() bool {
	return d.logConn != nil
}

func (d *Database) Driver() string {
	return d.driver
}

func (d *Database) LogTable() string {
	return d.logTable
}

func (d *Database) Quote(ident string) string {
	return quoteIdent(d.driver, ident)
}

func (d *Database) LogQuote(ident string) string {
	return quoteIdent(DriverClickHouse, ident)
}

const (
	LogTypeConsume = 2
	TokenTPSWindow = 60
)

func (d *Database) TokenTPS() (float64, error) {
	if d.logConn == nil {
		return 0, fmt.Errorf("clickhouse is not configured")
	}
	q := fmt.Sprintf(`
		SELECT COALESCE(SUM(prompt_tokens + completion_tokens), 0) / %d
		FROM %s
		WHERE type = %d AND created_at >= toUnixTimestamp(now()) - %d`,
		TokenTPSWindow, d.logTable, LogTypeConsume, TokenTPSWindow)
	var tps float64
	if err := d.logConn.QueryRow(q).Scan(&tps); err != nil {
		return 0, err
	}
	return tps, nil
}

func (d *Database) TokenTPSByChannel() ([]models.ChannelTPS, error) {
	if d.logConn == nil {
		return nil, fmt.Errorf("clickhouse is not configured")
	}
	q := fmt.Sprintf(`
		SELECT channel_id, COALESCE(SUM(prompt_tokens + completion_tokens), 0) / %d
		FROM %s
		WHERE type = %d AND created_at >= toUnixTimestamp(now()) - %d
		GROUP BY channel_id
		ORDER BY channel_id`,
		TokenTPSWindow, d.logTable, LogTypeConsume, TokenTPSWindow)
	rows, err := d.logConn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ChannelTPS
	for rows.Next() {
		var r models.ChannelTPS
		if err := rows.Scan(&r.ChannelID, &r.TPS); err != nil {
			return nil, err
		}
		r.ChannelName = d.channelNames[r.ChannelID]
		results = append(results, r)
	}
	return results, rows.Err()
}

func parseChannelNames(raw map[string]string) map[int64]string {
	names := make(map[int64]string, len(raw))
	for idStr, name := range raw {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		names[id] = name
	}
	return names
}

func quoteIdent(driver, ident string) string {
	if driver == DriverPostgres {
		return `"` + ident + `"`
	}
	return "`" + ident + "`"
}
