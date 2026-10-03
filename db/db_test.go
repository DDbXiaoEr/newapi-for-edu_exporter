package db

import (
	"testing"

	"github.com/quanxiaoer/newapi_exporter/config"
)

func TestPrimaryDSNPostgres(t *testing.T) {
	cfg := &config.Config{
		DBType:           "pgsql",
		PostgresHost:     "localhost",
		PostgresPort:     "5432",
		PostgresUser:     "root",
		PostgresPassword: "123456",
		PostgresDatabase: "newapi",
		PostgresSSLMode:  "disable",
	}

	driver, dsn, err := primaryDSN(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if driver != DriverPostgres {
		t.Errorf("expected driver postgres, got %s", driver)
	}
	if dsn != "postgres://root:123456@localhost:5432/newapi?sslmode=disable" {
		t.Errorf("unexpected dsn: %s", dsn)
	}
}

func TestPrimaryDSNRejectsSQLite(t *testing.T) {
	_, _, err := primaryDSN(&config.Config{DBType: "sqlite"})
	if err == nil {
		t.Fatal("expected sqlite to be rejected")
	}
}

func TestQuoteIdent(t *testing.T) {
	if got := quoteIdent(DriverPostgres, "group"); got != `"group"` {
		t.Errorf("postgres quote: %s", got)
	}
	if got := quoteIdent(DriverMySQL, "group"); got != "`group`" {
		t.Errorf("mysql quote: %s", got)
	}
	if got := quoteIdent(DriverClickHouse, "group"); got != "`group`" {
		t.Errorf("clickhouse quote: %s", got)
	}
}

func TestSQLDriverName(t *testing.T) {
	if got := sqlDriverName(DriverPostgres); got != "pgx" {
		t.Errorf("expected pgx, got %s", got)
	}
	if got := sqlDriverName(DriverMySQL); got != "mysql" {
		t.Errorf("expected mysql, got %s", got)
	}
}

func TestIsClickHouseDSN(t *testing.T) {
	if !isClickHouseDSN("clickhouse://default:@localhost:9000/newapi_audit") {
		t.Fatal("expected clickhouse DSN")
	}
	if isClickHouseDSN("postgresql://root:123456@localhost:5432/newapi") {
		t.Fatal("postgres DSN should not be treated as clickhouse")
	}
}

func TestIsValidIdent(t *testing.T) {
	if !isValidIdent("logs") || !isValidIdent("token_usage_logs") {
		t.Fatal("expected valid identifiers")
	}
	if isValidIdent("logs; drop table x") || isValidIdent("1logs") || isValidIdent("") {
		t.Fatal("expected invalid identifiers to be rejected")
	}
}

func TestParseChannelNames(t *testing.T) {
	got := parseChannelNames(map[string]string{"1": "OpenAI", "2": "Claude", "bad": "skip"})
	if got[1] != "OpenAI" || got[2] != "Claude" {
		t.Fatalf("unexpected names: %#v", got)
	}
	if _, ok := got[0]; ok {
		t.Fatal("invalid id should be skipped")
	}
}
