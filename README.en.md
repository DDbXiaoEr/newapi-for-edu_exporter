# newapi_exporter

Prometheus metrics exporter for new-api data.

Reads users / tokens / quota_data from PostgreSQL (or MySQL). Request logs may be read from ClickHouse. SQLite is not supported.

## Features

- **Prometheus Metrics Export**: Exports database data as metrics consumable by Prometheus
- **Multi-format Support**: JSON and CSV file export
- **Statistics Query**: User count, token usage, request count, and optional ClickHouse log stats
- **HTTP Service**: Built-in HTTP server for Prometheus to scrape metrics
- **Configuration Management**: YAML config file plus environment variables (`SQL_DSN` / `LOG_SQL_DSN` compatible with new-api)

## Technology Stack

- Go
- Prometheus Client
- PostgreSQL / MySQL
- ClickHouse (optional, logs)

## Quick Start

```bash
make build
./bin/newapi-exporter -out output -format json
./bin/newapi-exporter -serve 9090
```

### Configuration

```yaml
db_type: postgres
postgres_host: localhost
postgres_port: "5432"
postgres_user: root
postgres_password: "123456"
postgres_database: newapi
postgres_sslmode: disable
clickhouse_dsn: clickhouse://default:@localhost:9000/newapi_audit
clickhouse_table: logs
output_format: json
```

`SQL_DSN` and `LOG_SQL_DSN` from new-api are also accepted.

## Data Model

### Users / Tokens / Quota Data

Primary database tables.

### Logs

ClickHouse consume logs (`prompt_tokens`, `completion_tokens`, `quota`, ...).

### Statistics

- Active user count
- Total tokens / requests from `quota_data`
- Optional token TPS / tokens / quota from ClickHouse

## License

MIT License
