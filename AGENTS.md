# newapi_exporter

Golang project that exports usage statistics from a [new-api](https://github.com/Calcium-Ion/new-api) (LLM API gateway) database.

Primary data lives in PostgreSQL (or MySQL). Request/consume logs may live in ClickHouse (`LOG_SQL_DSN`). SQLite is not supported.

## Project Structure

```
.
├── main.go              # Entry point (export files or serve Prometheus metrics)
├── config/
│   ├── config.go        # Configuration via env vars / flags
│   └── config_test.go   # Config tests
├── db/
│   ├── db.go            # PostgreSQL / MySQL + optional ClickHouse connections
│   └── db_test.go
├── models/
│   ├── quota_data.go    # QuotaData struct
│   ├── token.go         # Token struct
│   ├── user.go          # User struct
│   ├── log.go           # ClickHouse log struct
│   └── statistics.go    # Statistics / ChannelTPS structs
├── exporter/
│   └── exporter.go      # Query + file export logic (json/csv)
├── metrics/
│   └── collector.go     # Prometheus collector (gauge metrics)
├── Makefile             # Build / test / lint targets
├── prometheus.yml       # Demo Prometheus scrape config
└── databasetable.md     # Database schema documentation
```

## Database

**Primary database** (`db_type: postgres` or `mysql`):

- **quota_data** - Per-user/per-model token usage: `user_id`, `username`, `model_name`, `token_used`, `count`, `quota`, plus `use_group`, `token_id`, `channel_id`, `node_name`
- **tokens** - API key tokens: `key`, `status`, `remain_quota`, `used_quota`, `group`, `model_limits`
- **users** - User accounts: `username`, `role`, `status`, `quota`, `used_quota`, `request_count`, `group`

**Log database** (optional ClickHouse via `clickhouse_dsn` / `LOG_SQL_DSN`):

- **logs** - Consume/audit-style request logs: `prompt_tokens`, `completion_tokens`, `quota`, `model_name`, `username`, `request_id`

See `databasetable.md` for full schema.

## Usage

**File export:**
```bash
go run . -out output -format json
```

**Prometheus HTTP server:**
```bash
go run . -serve 9090
```

**Flags:**
| Flag | Default | Description |
|---|---|---|
| `-config` | `config.yaml` | Path to YAML configuration file (current directory) |
| `-out` | `output` | Output directory (file export mode) |
| `-format` | `json` | Output format: `json` or `csv` |
| `-serve` | (none) | Start Prometheus HTTP server on port |

**Configuration file (`config.yaml`):**
```yaml
db_type: postgres        # postgres or mysql
output_format: json      # json or csv
postgres_host: localhost
postgres_port: "5432"
postgres_user: root
postgres_password: "123456"
postgres_database: newapi
postgres_sslmode: disable
clickhouse_dsn: clickhouse://default:@localhost:9000/newapi_audit
clickhouse_table: logs
channel_names:
  1: OpenAI
  2: Claude
```

Alternatively set `postgres_dsn` / `SQL_DSN` to a full URL, matching new-api.

**Environment variables (override config file):**
| Variable | Default |
|---|---|
| `NEWAPI_CONFIG_FILE` | `config.yaml` in current directory |
| `NEWAPI_OUTPUT_FORMAT` | `json` |
| `NEWAPI_DB_TYPE` | `postgres` |
| `NEWAPI_POSTGRES_DSN` / `SQL_DSN` | (empty; built from host/user/...) |
| `NEWAPI_POSTGRES_HOST` | `localhost` |
| `NEWAPI_POSTGRES_PORT` | `5432` |
| `NEWAPI_POSTGRES_USER` | `root` |
| `NEWAPI_POSTGRES_PASSWORD` | (empty) |
| `NEWAPI_POSTGRES_DATABASE` | `newapi` |
| `NEWAPI_POSTGRES_SSLMODE` | `disable` |
| `NEWAPI_MYSQL_HOST` | `localhost` |
| `NEWAPI_MYSQL_PORT` | `3306` |
| `NEWAPI_MYSQL_USER` | `root` |
| `NEWAPI_MYSQL_PASSWORD` | (empty) |
| `NEWAPI_MYSQL_DATABASE` | `newapi` |
| `NEWAPI_CLICKHOUSE_DSN` / `LOG_SQL_DSN` | (empty; optional) |
| `NEWAPI_CLICKHOUSE_TABLE` / `LOG_SQL_CLICKHOUSE_TABLE` | `logs` |

## Prometheus Metrics

When running with `-serve`, exposes at `/metrics`:

| Metric | Type | Description |
|---|---|---|
| `newapi_user_count` | Gauge | Active (non-deleted) users |
| `newapi_total_tokens` | Gauge | Sum of `token_used` in quota_data |
| `newapi_total_requests` | Gauge | Sum of `count` in quota_data |
| `newapi_token_tps` | Gauge | Tokens/sec from consume logs over last 60s (if ClickHouse configured) |
| `newapi_channel_token_tps` | Gauge | Per-channel tokens/sec over last 60s (`channel_id`, `channel_name` from `channel_names`) |
| `newapi_log_tokens` | Gauge | Sum of prompt+completion tokens in logs |
| `newapi_log_quota` | Gauge | Sum of quota in logs |

## Development

```bash
make build    # Build binary to bin/newapi-exporter
make test     # Run tests
make lint     # Run go vet
make clean    # Remove bin/ and output/
```

## Dependencies

- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/ClickHouse/clickhouse-go/v2` - ClickHouse driver
- `github.com/prometheus/client_golang` - Prometheus client
