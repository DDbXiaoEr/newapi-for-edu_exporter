# newapi_exporter

Prometheus metrics exporter for One-API data

## Project Overview

newapi_exporter is a Go application designed to export data from One-API databases as Prometheus metrics. It reads user, token, and quota data from an SQLite database and converts it into Prometheus monitoring metrics, while also supporting export in JSON and CSV formats.

## Features

- **Prometheus Metrics Export**: Exports database data as metrics consumable by Prometheus
- **Multi-format Support**: Supports exporting data in JSON and CSV formats
- **Statistics Query**: Provides system statistics including user count, token count, and quota usage
- **HTTP Service**: Built-in HTTP server for Prometheus to scrape metrics
- **Configuration Management**: Supports configuration via environment variables

## Technology Stack

- Go
- Prometheus Client
- SQLite

## Quick Start

### Build and Run

```bash
# Build
make build

# Or build manually
go build -o bin/newapi-exporter

# Run
./bin/newapi-exporter
```

### Configuration

Configuration is supported via environment variables:

| Environment Variable | Description | Default Value |
|----------------------|-------------|---------------|
| PORT | HTTP server port | 9520 |
| DB_PATH | Database file path | one-api.db |
| OUTPUT_PATH | Output directory | output |

### Prometheus Configuration

Refer to `prometheus.yml.example` to configure Prometheus for metric scraping:

```yaml
scrape_configs:
  - job_name: 'newapi_exporter'
    static_configs:
      - targets: ['localhost:9520']
```

## Data Model

### Users

- User basic information
- Quota limits
- Usage statistics

### Tokens

- API token information
- Associated user
- Rate limit configuration

### Quota Data

- Quota usage records
- Timestamp information
- Usage statistics

### Statistics

- Total number of users
- Total number of tokens
- Number of active quotas

## Output Examples

Exported data is saved in the `output/` directory:

- `users.json` - User data
- `tokens.json` - Token data
- `quota_data.json` - Quota data
- `statistics.json` - Statistics data

## License

MIT License