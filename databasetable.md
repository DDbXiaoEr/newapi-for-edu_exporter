# Database Tables Structure

## quota_data

### Table Structure
```sql
CREATE TABLE `quota_data` (
    `id` integer,
    `user_id` integer,
    `username` text DEFAULT "",
    `model_name` text DEFAULT "",
    `created_at` integer,
    `token_used` integer DEFAULT 0,
    `count` integer DEFAULT 0,
    `quota` integer DEFAULT 0,
    PRIMARY KEY (`id`)
);
```

### Indexes
```sql
CREATE INDEX `idx_qdt_model_user_name` ON `quota_data`(`model_name`,`username`);
CREATE INDEX `idx_quota_data_user_id` ON `quota_data`(`user_id`);
CREATE INDEX `idx_qdt_created_at` ON `quota_data`(`created_at`);
```

### Field Descriptions
| Field | Type | Description | Default |
|-------|------|-------------|---------|
| id | integer | Primary key identifier | - |
| user_id | integer | User identifier | - |
| username | text | Username | "" |
| model_name | text | Name of the model | "" |
| created_at | integer | Creation timestamp | - |
| token_used | integer | Number of tokens used | 0 |
| count | integer | Usage count | 0 |
| quota | integer | Quota amount | 0 |

---

## tokens

### Table Structure
```sql
CREATE TABLE `tokens` (
    `id` integer,
    `user_id` integer,
    `key` char(48),
    `status` integer DEFAULT 1,
    `name` text,
    `created_time` integer,
    `accessed_time` integer,
    `expired_time` integer DEFAULT -1,
    `remain_quota` integer DEFAULT 0,
    `unlimited_quota` numeric,
    `model_limits_enabled` numeric,
    `model_limits` text,
    `allow_ips` text DEFAULT "",
    `used_quota` integer DEFAULT 0,
    `group` text DEFAULT "",
    `cross_group_retry` numeric,
    `deleted_at` datetime,
    PRIMARY KEY (`id`)
);
```

### Indexes
```sql
CREATE UNIQUE INDEX `idx_tokens_key` ON `tokens`(`key`);
CREATE INDEX `idx_tokens_user_id` ON `tokens`(`user_id`);
CREATE INDEX `idx_tokens_deleted_at` ON `tokens`(`deleted_at`);
CREATE INDEX `idx_tokens_name` ON `tokens`(`name`);
```

### Field Descriptions
| Field | Type | Description | Default |
|-------|------|-------------|---------|
| id | integer | Primary key identifier | - |
| user_id | integer | User identifier | - |
| key | char(48) | Token key | - |
| status | integer | Token status | 1 |
| name | text | Token name | - |
| created_time | integer | Creation timestamp | - |
| accessed_time | integer | Last accessed timestamp | - |
| expired_time | integer | Expiration timestamp | -1 |
| remain_quota | integer | Remaining quota | 0 |
| unlimited_quota | numeric | Unlimited quota flag | - |
| model_limits_enabled | numeric | Model limits enabled flag | - |
| model_limits | text | Model limits configuration | - |
| allow_ips | text | Allowed IP addresses | "" |
| used_quota | integer | Used quota amount | 0 |
| group | text | Token group | "" |
| cross_group_retry | numeric | Cross group retry flag | - |
| deleted_at | datetime | Deletion timestamp | - |

---

## users

### Table Structure
```sql
CREATE TABLE IF NOT EXISTS "users" (
    `id` integer,
    `username` text UNIQUE,
    `password` text NOT NULL,
    `display_name` text,
    `role` integer DEFAULT 1,
    `status` integer DEFAULT 1,
    `email` text,
    `github_id` text,
    `discord_id` text,
    `oidc_id` text,
    `wechat_id` text,
    `telegram_id` text,
    `access_token` char(32),
    `quota` integer DEFAULT 0,
    `used_quota` integer DEFAULT 0,
    `request_count` integer DEFAULT 0,
    `group` varchar(64) DEFAULT "default",
    `aff_code` varchar(32),
    `aff_count` integer DEFAULT 0,
    `aff_quota` integer DEFAULT 0,
    `aff_history` integer DEFAULT 0,
    `inviter_id` integer,
    `deleted_at` datetime,
    `linux_do_id` text,
    `setting` text,
    `remark` varchar(255),
    `stripe_customer` varchar(64),
    PRIMARY KEY (`id`)
);
```

### Indexes
```sql
CREATE INDEX `idx_users_stripe_customer` ON `users`(`stripe_customer`);
CREATE INDEX `idx_users_display_name` ON `users`(`display_name`);
CREATE INDEX `idx_users_git_hub_id` ON `users`(`github_id`);
CREATE INDEX `idx_users_we_chat_id` ON `users`(`wechat_id`);
CREATE INDEX `idx_users_telegram_id` ON `users`(`telegram_id`);
CREATE UNIQUE INDEX `idx_users_access_token` ON `users`(`access_token`);
CREATE UNIQUE INDEX `idx_users_aff_code` ON `users`(`aff_code`);
CREATE INDEX `idx_users_deleted_at` ON `users`(`deleted_at`);
CREATE INDEX `idx_users_linux_do_id` ON `users`(`linux_do_id`);
CREATE INDEX `idx_users_username` ON `users`(`username`);
CREATE INDEX `idx_users_email` ON `users`(`email`);
CREATE INDEX `idx_users_discord_id` ON `users`(`discord_id`);
CREATE INDEX `idx_users_oidc_id` ON `users`(`oidc_id`);
CREATE INDEX `idx_users_inviter_id` ON `users`(`inviter_id`);
```

### Field Descriptions
| Field | Type | Description | Default |
|-------|------|-------------|---------|
| id | integer | Primary key identifier | - |
| username | text | Username (unique) | - |
| password | text | Password (not null) | - |
| display_name | text | Display name | - |
| role | integer | User role | 1 |
| status | integer | Account status | 1 |
| email | text | Email address | - |
| github_id | text | GitHub ID | - |
| discord_id | text | Discord ID | - |
| oidc_id | text | OpenID Connect ID | - |
| wechat_id | text | WeChat ID | - |
| telegram_id | text | Telegram ID | - |
| access_token | char(32) | Access token | - |
| quota | integer | User quota | 0 |
| used_quota | integer | Used quota amount | 0 |
| request_count | integer | Request count | 0 |
| group | varchar(64) | User group | "default" |
| aff_code | varchar(32) | Affiliate code | - |
| aff_count | integer | Affiliate count | 0 |
| aff_quota | integer | Affiliate quota | 0 |
| aff_history | integer | Affiliate history | 0 |
| inviter_id | integer | Inviter user ID | - |
| deleted_at | datetime | Deletion timestamp | - |
| linux_do_id | text | Linux Do ID | - |
| setting | text | User settings | - |
| remark | varchar(255) | User remark | - |
| stripe_customer | varchar(64) | Stripe customer ID | - |