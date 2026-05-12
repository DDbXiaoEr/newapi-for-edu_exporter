# newapi_exporter

Golang project that exports usage statistics from a [new-api](https://github.com/Calcium-Ion/new-api) (LLM API gateway) SQLite database.

## Project Structure

```
.
├── main.go              # Entry point (export files or serve Prometheus metrics)
├── config/
│   ├── config.go        # Configuration via env vars / flags
│   └── config_test.go   # Config tests
├── db/
│   └── db.go            # SQLite database connection
├── models/
│   ├── quota_data.go    # QuotaData struct
│   ├── token.go         # Token struct
│   ├── user.go          # User struct
│   └── statistics.go    # Statistics struct (user_count, total_tokens, total_requests)
├── exporter/
│   └── exporter.go      # Query + file export logic (json/csv)
├── metrics/
│   └── collector.go     # Prometheus collector (gauge metrics)
├── Makefile             # Build / test / lint targets
├── prometheus.yml       # Demo Prometheus scrape config
└── databasetable.md     # Database schema documentation
```

## Database

SQLite database (new-api stores data in SQLite by default). Expects three tables:

- **quota_data** - Per-user/per-model token usage: `user_id`, `username`, `model_name`, `token_used`, `count`, `quota`
- **tokens** - API key tokens: `key`, `status`, `remain_quota`, `used_quota`, `group`, `model_limits`
- **users** - User accounts: `username`, `role`, `status`, `quota`, `used_quota`, `request_count`, `group`

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
| `-config` | (auto-detect) | Path to YAML configuration file |
| `-out` | `output` | Output directory (file export mode) |
| `-format` | `json` | Output format: `json` or `csv` |
| `-serve` | (none) | Start Prometheus HTTP server on port |

**Configuration file (`config.yaml`):**
```yaml
db_type: sqlite          # sqlite or mysql
database_path: newapi.db # path to SQLite database file
output_format: json      # json or csv
mysql_host: localhost
mysql_port: "3306"
mysql_user: root
mysql_password: ""
mysql_database: newapi
```

**Environment variables (override config file):**
| Variable | Default |
|---|---|
| `NEWAPI_CONFIG_FILE` | (auto-detect `config.yaml`/`config.yml`) |
| `NEWAPI_DB_PATH` | `newapi.db` |
| `NEWAPI_OUTPUT_FORMAT` | `json` |
| `NEWAPI_DB_TYPE` | `sqlite` |
| `NEWAPI_MYSQL_HOST` | `localhost` |
| `NEWAPI_MYSQL_PORT` | `3306` |
| `NEWAPI_MYSQL_USER` | `root` |
| `NEWAPI_MYSQL_PASSWORD` | (empty) |
| `NEWAPI_MYSQL_DATABASE` | `newapi` |

## Prometheus Metrics

When running with `-serve`, exposes at `/metrics`:

| Metric | Type | Description |
|---|---|---|
| `newapi_user_count` | Gauge | Active (non-deleted) users |
| `newapi_total_tokens` | Gauge | Sum of `token_used` in quota_data |
| `newapi_total_requests` | Gauge | Sum of `count` in quota_data |

## Development

```bash
make build    # Build binary to bin/newapi-exporter
make test     # Run tests
make lint     # Run go vet
make clean    # Remove bin/ and output/
```

## Dependencies

- `github.com/mattn/go-sqlite3` - SQLite driver (CGO required)
- `github.com/prometheus/client_golang` - Prometheus client
