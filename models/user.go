package models

import (
	"database/sql"
)

type User struct {
	ID             int64          `json:"id"`
	Username       string         `json:"username"`
	Password       string         `json:"-"`
	DisplayName    sql.NullString `json:"display_name"`
	Role           int64          `json:"role"`
	Status         int64          `json:"status"`
	Email          sql.NullString `json:"email"`
	GitHubID       sql.NullString `json:"github_id"`
	DiscordID      sql.NullString `json:"discord_id"`
	OidcID         sql.NullString `json:"oidc_id"`
	WechatID       sql.NullString `json:"wechat_id"`
	TelegramID     sql.NullString `json:"telegram_id"`
	AccessToken    sql.NullString `json:"access_token"`
	Quota          int64          `json:"quota"`
	UsedQuota      int64          `json:"used_quota"`
	RequestCount   int64          `json:"request_count"`
	Group          string         `json:"group"`
	AffCode        sql.NullString `json:"aff_code"`
	AffCount       int64          `json:"aff_count"`
	AffQuota       int64          `json:"aff_quota"`
	AffHistory     int64          `json:"aff_history"`
	InviterID      sql.NullInt64  `json:"inviter_id"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
	LinuxDoID      sql.NullString `json:"linux_do_id"`
	Setting        sql.NullString `json:"setting"`
	Remark         sql.NullString `json:"remark"`
	StripeCustomer sql.NullString `json:"stripe_customer"`
}
