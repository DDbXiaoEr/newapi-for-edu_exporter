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
go run . -db one-api.db -out output -format json
```

**Prometheus HTTP server:**
```bash
go run . -db one-api.db -serve 9090
```

**Flags:**
| Flag | Default | Description |
|---|---|---|
| `-db` | `newapi.db` | Path to SQLite database |
| `-out` | `output` | Output directory (file export mode) |
| `-format` | `json` | Output format: `json` or `csv` |
| `-serve` | (none) | Start Prometheus HTTP server on port |

**Environment variables (fallback for flags):**
| Variable | Default |
|---|---|
| `NEWAPI_DB_PATH` | `newapi.db` |
| `NEWAPI_OUTPUT_FORMAT` | `json` |

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
