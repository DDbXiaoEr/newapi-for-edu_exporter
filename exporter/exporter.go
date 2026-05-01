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
	rows, err := e.db.Conn().Query("SELECT id, user_id, username, model_name, created_at, token_used, count, quota FROM quota_data")
	if err != nil {
		return nil, fmt.Errorf("query quota_data: %w", err)
	}
	defer rows.Close()

	var results []models.QuotaData
	for rows.Next() {
		var r models.QuotaData
		if err := rows.Scan(&r.ID, &r.UserID, &r.Username, &r.ModelName, &r.CreatedAt, &r.TokenUsed, &r.Count, &r.Quota); err != nil {
			return nil, fmt.Errorf("scan quota_data: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (e *Exporter) ExportTokens() ([]models.Token, error) {
	rows, err := e.db.Conn().Query(`
		SELECT id, user_id, ` + "`key`" + `, status, ` + "`name`" + `, created_time, accessed_time, expired_time,
		       remain_quota, unlimited_quota, model_limits_enabled, model_limits,
		       allow_ips, used_quota, ` + "`group`" + `, cross_group_retry, deleted_at
		FROM tokens`)
	if err != nil {
		return nil, fmt.Errorf("query tokens: %w", err)
	}
	defer rows.Close()

	var results []models.Token
	for rows.Next() {
		var r models.Token
		if err := rows.Scan(&r.ID, &r.UserID, &r.Key, &r.Status, &r.Name, &r.CreatedTime,
			&r.AccessedTime, &r.ExpiredTime, &r.RemainQuota, &r.UnlimitedQuota,
			&r.ModelLimitsEnabled, &r.ModelLimits, &r.AllowIPs, &r.UsedQuota,
			&r.Group, &r.CrossGroupRetry, &r.DeletedAt); err != nil {
			return nil, fmt.Errorf("scan tokens: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (e *Exporter) ExportUsers() ([]models.User, error) {
	rows, err := e.db.Conn().Query(`
		SELECT id, username, password, display_name, role, status, email,
		       github_id, discord_id, oidc_id, wechat_id, telegram_id,
		       access_token, quota, used_quota, request_count, ` + "`group`" + `,
		       aff_code, aff_count, aff_quota, aff_history, inviter_id,
		       deleted_at, linux_do_id, setting, remark, stripe_customer
		FROM users`)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var results []models.User
	for rows.Next() {
		var r models.User
		if err := rows.Scan(&r.ID, &r.Username, &r.Password, &r.DisplayName, &r.Role,
			&r.Status, &r.Email, &r.GitHubID, &r.DiscordID, &r.OidcID,
			&r.WechatID, &r.TelegramID, &r.AccessToken, &r.Quota, &r.UsedQuota,
			&r.RequestCount, &r.Group, &r.AffCode, &r.AffCount, &r.AffQuota,
			&r.AffHistory, &r.InviterID, &r.DeletedAt, &r.LinuxDoID,
			&r.Setting, &r.Remark, &r.StripeCustomer); err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
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

	return s, nil
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
		return fmt.Errorf("unmarshal: %w", err)
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
