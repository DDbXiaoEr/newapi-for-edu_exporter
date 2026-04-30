package models

import (
	"database/sql"
)

type Token struct {
	ID                 int64          `json:"id"`
	UserID             int64          `json:"user_id"`
	Key                string         `json:"key"`
	Status             int64          `json:"status"`
	Name               sql.NullString `json:"name"`
	CreatedTime        int64          `json:"created_time"`
	AccessedTime       int64          `json:"accessed_time"`
	ExpiredTime        int64          `json:"expired_time"`
	RemainQuota        int64          `json:"remain_quota"`
	UnlimitedQuota     sql.NullBool   `json:"unlimited_quota"`
	ModelLimitsEnabled sql.NullBool   `json:"model_limits_enabled"`
	ModelLimits        sql.NullString `json:"model_limits"`
	AllowIPs           string         `json:"allow_ips"`
	UsedQuota          int64          `json:"used_quota"`
	Group              string         `json:"group"`
	CrossGroupRetry    sql.NullBool   `json:"cross_group_retry"`
	DeletedAt          sql.NullTime   `json:"deleted_at"`
}
