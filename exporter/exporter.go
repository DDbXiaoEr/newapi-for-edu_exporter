package exporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/quanxiaoer/newapi_exporter/config"
	"github.com/quanxiaoer/newapi_exporter/db"
	"github.com/quanxiaoer/newapi_exporter/models"
)

type Exporter struct {
	db  *db.Database
	cfg *config.Config
}

func New(database *db.Database, cfg *config.Config) *Exporter {
	return &Exporter{db: database, cfg: cfg}
}

func (e *Exporter) ExportQuotaData() ([]models.QuotaData, error) {
	rows, err := e.db.Conn().Query(`
		SELECT id, user_id, username, model_name, created_at, use_group, token_id, channel_id, node_name, token_used, count, quota
		FROM quota_data`)
	if err != nil {
		return nil, fmt.Errorf("query quota_data: %w", err)
	}
	defer rows.Close()

	var results []models.QuotaData
	for rows.Next() {
		var r models.QuotaData
		if err := rows.Scan(
			&r.ID, &r.UserID, &r.Username, &r.ModelName, &r.CreatedAt,
			&r.UseGroup, &r.TokenID, &r.ChannelID, &r.NodeName,
			&r.TokenUsed, &r.Count, &r.Quota,
		); err != nil {
			return nil, fmt.Errorf("scan quota_data: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (e *Exporter) ExportTokens() ([]models.Token, error) {
	q := fmt.Sprintf(`
		SELECT id, user_id, %s, status, name, created_time, accessed_time, expired_time,
		       remain_quota, unlimited_quota, model_limits_enabled, model_limits,
		       allow_ips, used_quota, %s, cross_group_retry, auto_groups, deleted_at
		FROM tokens`, e.db.Quote("key"), e.db.Quote("group"))
	rows, err := e.db.Conn().Query(q)
	if err != nil {
		return nil, fmt.Errorf("query tokens: %w", err)
	}
	defer rows.Close()

	var results []models.Token
	for rows.Next() {
		var r models.Token
		if err := rows.Scan(
			&r.ID, &r.UserID, &r.Key, &r.Status, &r.Name, &r.CreatedTime,
			&r.AccessedTime, &r.ExpiredTime, &r.RemainQuota, &r.UnlimitedQuota,
			&r.ModelLimitsEnabled, &r.ModelLimits, &r.AllowIPs, &r.UsedQuota,
			&r.Group, &r.CrossGroupRetry, &r.AutoGroups, &r.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan tokens: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (e *Exporter) ExportUsers() ([]models.User, error) {
	q := fmt.Sprintf(`
		SELECT id, username, password, display_name, role, status, email,
		       github_id, discord_id, oidc_id, wechat_id, telegram_id,
		       access_token, access_token_created_at, quota, used_quota, request_count, %s,
		       aff_code, aff_count, aff_quota, aff_history, inviter_id,
		       deleted_at, linux_do_id, cas_id, setting, remark, stripe_customer,
		       created_at, last_login_at
		FROM users`, e.db.Quote("group"))
	rows, err := e.db.Conn().Query(q)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var results []models.User
	for rows.Next() {
		var r models.User
		if err := rows.Scan(
			&r.ID, &r.Username, &r.Password, &r.DisplayName, &r.Role,
			&r.Status, &r.Email, &r.GitHubID, &r.DiscordID, &r.OidcID,
			&r.WechatID, &r.TelegramID, &r.AccessToken, &r.AccessTokenCreatedAt,
			&r.Quota, &r.UsedQuota, &r.RequestCount, &r.Group, &r.AffCode,
			&r.AffCount, &r.AffQuota, &r.AffHistory, &r.InviterID, &r.DeletedAt,
			&r.LinuxDoID, &r.CasID, &r.Setting, &r.Remark, &r.StripeCustomer,
			&r.CreatedAt, &r.LastLoginAt,
		); err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (e *Exporter) ExportLogs() ([]models.Log, error) {
	if !e.db.HasLogs() {
		return nil, fmt.Errorf("clickhouse is not configured")
	}

	q := fmt.Sprintf(`
		SELECT id, user_id, created_at, type, content, username, token_name, model_name,
		       quota, prompt_tokens, completion_tokens, use_time, is_stream,
		       channel_id, token_id, %s, ip, request_id, upstream_request_id, other
		FROM %s`, e.db.LogQuote("group"), e.db.LogTable())
	rows, err := e.db.LogConn().Query(q)
	if err != nil {
		return nil, fmt.Errorf("query logs: %w", err)
	}
	defer rows.Close()

	var results []models.Log
	for rows.Next() {
		var r models.Log
		if err := rows.Scan(
			&r.ID, &r.UserID, &r.CreatedAt, &r.Type, &r.Content, &r.Username,
			&r.TokenName, &r.ModelName, &r.Quota, &r.PromptTokens, &r.CompletionTokens,
			&r.UseTime, &r.IsStream, &r.ChannelID, &r.TokenID, &r.Group, &r.IP,
			&r.RequestID, &r.UpstreamRequestID, &r.Other,
		); err != nil {
			return nil, fmt.Errorf("scan logs: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (e *Exporter) ExportStatistics() (models.Statistics, error) {
	var s models.Statistics

	if err := e.db.Conn().QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&s.UserCount); err != nil {
		return s, fmt.Errorf("query user_count: %w", err)
	}

	if err := e.db.Conn().QueryRow("SELECT COALESCE(SUM(token_used), 0) FROM quota_data").Scan(&s.TotalTokens); err != nil {
		return s, fmt.Errorf("query total_tokens: %w", err)
	}

	if err := e.db.Conn().QueryRow("SELECT COALESCE(SUM(count), 0) FROM quota_data").Scan(&s.TotalRequests); err != nil {
		return s, fmt.Errorf("query total_requests: %w", err)
	}

	if e.db.HasLogs() {
		table := e.db.LogTable()
		tps, err := e.db.TokenTPS()
		if err != nil {
			return s, fmt.Errorf("query token_tps: %w", err)
		}
		s.TokenTPS = tps
		channelTPS, err := e.db.TokenTPSByChannel()
		if err != nil {
			return s, fmt.Errorf("query channel_tps: %w", err)
		}
		s.ChannelTPS = channelTPS
		if err := e.db.LogConn().QueryRow(fmt.Sprintf("SELECT COALESCE(SUM(prompt_tokens + completion_tokens), 0) FROM %s", table)).Scan(&s.LogTokens); err != nil {
			return s, fmt.Errorf("query log_tokens: %w", err)
		}
		if err := e.db.LogConn().QueryRow(fmt.Sprintf("SELECT COALESCE(SUM(quota), 0) FROM %s", table)).Scan(&s.LogQuota); err != nil {
			return s, fmt.Errorf("query log_quota: %w", err)
		}
	}

	return s, nil
}

func (e *Exporter) ExportChannelTPS() ([]models.ChannelTPS, error) {
	if !e.db.HasLogs() {
		return nil, fmt.Errorf("clickhouse is not configured")
	}
	return e.db.TokenTPSByChannel()
}

func (e *Exporter) WriteOutput(data interface{}, path string, format string) error {
	switch format {
	case "json":
		return e.writeJSON(data, path)
	case "csv":
		return e.writeCSV(data, path)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func (e *Exporter) writeJSON(data interface{}, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (e *Exporter) writeCSV(data interface{}, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	var records []map[string]interface{}
	if err := json.Unmarshal(b, &records); err != nil {
		var single map[string]interface{}
		if err2 := json.Unmarshal(b, &single); err2 != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		records = []map[string]interface{}{single}
	}

	if len(records) == 0 {
		return nil
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := make([]string, 0, len(records[0]))
	for k := range records[0] {
		headers = append(headers, k)
	}
	writer.Write(headers)

	for _, record := range records {
		row := make([]string, len(headers))
		for i, h := range headers {
			row[i] = fmt.Sprintf("%v", record[h])
		}
		writer.Write(row)
	}
	return nil
}
