# Database Tables Structure

Primary store: PostgreSQL (or MySQL). Request logs may live in ClickHouse.

## quota_data (primary DB)

```sql
CREATE TABLE quota_data (
    id integer,
    user_id integer,
    username text DEFAULT '',
    model_name text DEFAULT '',
    created_at bigint,
    use_group text DEFAULT '',
    token_id integer DEFAULT 0,
    channel_id integer DEFAULT 0,
    node_name text DEFAULT '',
    token_used integer DEFAULT 0,
    count integer DEFAULT 0,
    quota integer DEFAULT 0,
    PRIMARY KEY (id)
);
```

| Field | Description |
|-------|-------------|
| id | Primary key |
| user_id | User identifier |
| username | Username |
| model_name | Model name |
| created_at | Hour-truncated unix timestamp |
| use_group | Billing / routing group |
| token_id | Token identifier |
| channel_id | Channel identifier |
| node_name | Node name |
| token_used | Tokens used |
| count | Request count |
| quota | Quota consumed |

---

## tokens (primary DB)

PostgreSQL requires quoting `"key"` and `"group"`.

| Field | Description |
|-------|-------------|
| id | Primary key |
| user_id | User identifier |
| key | Token key |
| status | Token status |
| name | Token name |
| created_time / accessed_time / expired_time | Unix timestamps |
| remain_quota / used_quota / unlimited_quota | Quota fields |
| model_limits_enabled / model_limits | Model limit config |
| allow_ips | Allowed IPs |
| group | Token group |
| cross_group_retry | Cross-group retry flag |
| auto_groups | Auto group list (JSON text) |
| deleted_at | Soft delete timestamp |

---

## users (primary DB)

PostgreSQL requires quoting `"group"`.

| Field | Description |
|-------|-------------|
| id | Primary key |
| username | Unique username |
| password | Password hash (not exported in JSON via `json:"-"`) |
| display_name / role / status / email | Account fields |
| github_id / discord_id / oidc_id / wechat_id / telegram_id / linux_do_id / cas_id | Bindings |
| access_token / access_token_created_at | Dashboard PAT |
| quota / used_quota / request_count | Usage |
| group | User group |
| aff_code / aff_count / aff_quota / aff_history / inviter_id | Affiliate |
| deleted_at | Soft delete |
| setting / remark / stripe_customer | Extra |
| created_at / last_login_at | Unix timestamps |

---

## logs (ClickHouse, optional)

Default table name `logs` (override with `LOG_SQL_CLICKHOUSE_TABLE`).

```sql
CREATE TABLE IF NOT EXISTS logs (
    id Int64 DEFAULT 0,
    user_id Int32 DEFAULT 0,
    created_at Int64 DEFAULT 0,
    type Int32 DEFAULT 0,
    content String DEFAULT '',
    username String DEFAULT '',
    token_name String DEFAULT '',
    model_name String DEFAULT '',
    quota Int32 DEFAULT 0,
    prompt_tokens Int32 DEFAULT 0,
    completion_tokens Int32 DEFAULT 0,
    use_time Int32 DEFAULT 0,
    is_stream UInt8 DEFAULT 0,
    channel_id Int32 DEFAULT 0,
    token_id Int32 DEFAULT 0,
    `group` String DEFAULT '',
    ip String DEFAULT '',
    request_id String DEFAULT '',
    upstream_request_id String DEFAULT '',
    other String DEFAULT ''
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(toDateTime(created_at))
ORDER BY (created_at, request_id);
```
