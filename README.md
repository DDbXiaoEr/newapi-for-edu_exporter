


# newapi_exporter
本仓库所有文件均由AI生成（因为本人是个偏瘫患者所以求轻吐槽代码质量）

One-API 数据 Prometheus 指标导出器

## 项目简介

newapi_exporter 是一个用于将 One-API 数据库中的数据导出为 Prometheus 指标的 Go 应用。它可以从 SQLite 数据库中读取用户、令牌和配额数据，并将这些数据转换为 Prometheus 监控指标，同时支持导出为 JSON 和 CSV 格式。

## 功能特性

- **Prometheus 指标导出**：将数据库数据导出为 Prometheus 可拉取的指标
- **多格式支持**：支持 JSON 和 CSV 格式的数据导出
- **统计数据查询**：提供系统统计数据，包括用户数、令牌数、配额使用等
- **HTTP 服务**：内置 HTTP 服务器用于 Prometheus 拉取指标
- **配置管理**：支持通过环境变量进行配置

## 技术栈

- Go
- Prometheus Client
- SQLite

## 快速开始

### 编译运行

```bash
# 编译
make build

# 或手动编译
go build -o bin/newapi-exporter

# 运行
./bin/newapi-exporter
```

### 配置说明

支持通过环境变量配置：

| 环境变量 | 说明 | 默认值 |
|---------|------|-------|
| PORT | HTTP 服务端口 | 9520 |
| DB_PATH | 数据库文件路径 | one-api.db |
| OUTPUT_PATH | 输出目录 | output |

### Prometheus 配置

参考 `prometheus.yml.example` 配置 Prometheus 来拉取指标：

```yaml
scrape_configs:
  - job_name: 'newapi_exporter'
    static_configs:
      - targets: ['localhost:9520']
```

## 数据模型

### 用户 (users)

- 用户基本信息
- 配额限制
- 使用统计

### 令牌 (tokens)

- API 令牌信息
- 关联用户
- 限额配置

### 配额数据 (quota_data)

- 配额使用记录
- 时间戳信息
- 使用量统计

### 统计数据 (statistics)

- 系统总用户数
- 总令牌数
- 活跃配额数

## 输出示例

导出的数据会保存在 `output/` 目录下：

- `users.json` - 用户数据
- `tokens.json` - 令牌数据
- `quota_data.json` - 配额数据
- `statistics.json` - 统计数据

## 许可证

MIT License