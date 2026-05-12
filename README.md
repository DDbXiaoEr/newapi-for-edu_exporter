# newapi_exporter

本仓库所有文件均由AI生成（因为本人是个偏瘫患者所以求轻吐槽代码质量）

[new-api](https://github.com/Calcium-Ion/new-api) (LLM API gateway) 数据导出器 —— 支持 Prometheus 指标暴露与 JSON/CSV 文件导出。

## 功能特性

- **双模式运行**：文件导出或 Prometheus HTTP 指标服务
- **多数据库支持**：SQLite 和 MySQL
- **多格式输出**：JSON 和 CSV
- **Prometheus 指标**：活跃用户数、总 Token 使用量、总请求数

## 快速开始

```bash
# 编译
make build

# 文件导出模式（SQLite）
./bin/newapi-exporter -out output -format json

# Prometheus 指标服务
./bin/newapi-exporter -serve 9090
```

## 配置

数据库连接通过配置文件指定（默认 `config.yaml` / `config.yml`），也可以通过环境变量覆盖。

### 配置文件 (`config.yaml`)

```yaml
db_type: sqlite          # sqlite 或 mysql
database_path: newapi.db # SQLite 数据库文件路径
output_format: json      # json 或 csv
mysql_host: localhost
mysql_port: "3306"
mysql_user: root
mysql_password: ""
mysql_database: newapi
```

### 配置优先级

命令行参数 > 环境变量 > 配置文件 > 默认值

## 命令行参数

| 参数 | 默认值 | 说明 |
|---|---|---|
| `-config` | 自动检测 | YAML 配置文件路径 |
| `-out` | `output` | 文件导出输出目录 |
| `-format` | `json` | 输出格式：`json` / `csv` |
| `-serve` | (空) | 启动 Prometheus HTTP 服务，值为端口号 |
| `-version` | - | 打印版本信息并退出 |
| `-help` | - | 打印帮助信息并退出 |

## 环境变量

| 环境变量 | 默认值 | 说明 |
|---|---|---|
| `NEWAPI_CONFIG_FILE` | 自动检测 `config.yaml` / `config.yml` | 配置文件路径 |
| `NEWAPI_DB_TYPE` | `sqlite` | 数据库类型 |
| `NEWAPI_DB_PATH` | `newapi.db` | SQLite 数据库文件路径 |
| `NEWAPI_OUTPUT_FORMAT` | `json` | 输出格式 |
| `NEWAPI_MYSQL_HOST` | `localhost` | MySQL 主机地址 |
| `NEWAPI_MYSQL_PORT` | `3306` | MySQL 端口 |
| `NEWAPI_MYSQL_USER` | `root` | MySQL 用户名 |
| `NEWAPI_MYSQL_PASSWORD` | (空) | MySQL 密码 |
| `NEWAPI_MYSQL_DATABASE` | `newapi` | MySQL 数据库名 |

## Prometheus 指标

启动 `-serve` 后，访问 `/metrics` 暴露以下指标：

| 指标 | 类型 | 说明 |
|---|---|---|
| `newapi_user_count` | Gauge | 活跃（未删除）用户数 |
| `newapi_total_tokens` | Gauge | quota_data 中 token_used 总和 |
| `newapi_total_requests` | Gauge | quota_data 中 count 总和 |

### Prometheus 配置

```yaml
scrape_configs:
  - job_name: 'newapi_exporter'
    static_configs:
      - targets: ['localhost:9090']
```

## 导出文件

文件导出模式会在输出目录生成以下文件：

- `statistics.{format}` - 统计概览（用户数、Token、请求量）
- `quota_data.{format}` - 配额使用记录
- `tokens.{format}` - API 令牌信息
- `users.{format}` - 用户信息

## 开发

```bash
make build    # 编译到 bin/newapi-exporter
make test     # 运行测试
make lint     # 代码检查 (go vet)
make clean    # 清理构建产物和输出目录
```

## 依赖

- `github.com/mattn/go-sqlite3` - SQLite 驱动（需 CGO）
- `github.com/go-sql-driver/mysql` - MySQL 驱动
- `github.com/prometheus/client_golang` - Prometheus 客户端

## 许可证

MIT License
