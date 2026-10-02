package metrics

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/quanxiaoer/newapi_exporter/db"
)

type Collector struct {
	db *db.Database

	userCount     *prometheus.Desc
	totalTokens   *prometheus.Desc
	totalRequests *prometheus.Desc
	tokenTPS      *prometheus.Desc
	logTokens     *prometheus.Desc
	logQuota      *prometheus.Desc
}

func NewCollector(database *db.Database) *Collector {
	return &Collector{
		db: database,
		userCount: prometheus.NewDesc(
			"newapi_user_count",
			"Total number of active users",
			nil, nil,
		),
		totalTokens: prometheus.NewDesc(
			"newapi_total_tokens",
			"Total tokens used across all models (quota_data)",
			nil, nil,
		),
		totalRequests: prometheus.NewDesc(
			"newapi_total_requests",
			"Total API requests across all models (quota_data)",
			nil, nil,
		),
		tokenTPS: prometheus.NewDesc(
			"newapi_token_tps",
			"Tokens per second from consume logs over the last 60 seconds",
			nil, nil,
		),
		logTokens: prometheus.NewDesc(
			"newapi_log_tokens",
			"Sum of prompt_tokens + completion_tokens in ClickHouse logs",
			nil, nil,
		),
		logQuota: prometheus.NewDesc(
			"newapi_log_quota",
			"Sum of quota in ClickHouse logs",
			nil, nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.userCount
	ch <- c.totalTokens
	ch <- c.totalRequests
	if c.db.HasLogs() {
		ch <- c.tokenTPS
		ch <- c.logTokens
		ch <- c.logQuota
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var userCount, totalTokens, totalRequests float64

	if err := c.db.Conn().QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&userCount); err != nil {
		log.Printf("metrics: query user_count: %v", err)
		userCount = 0
	}

	if err := c.db.Conn().QueryRow("SELECT COALESCE(SUM(token_used), 0) FROM quota_data").Scan(&totalTokens); err != nil {
		log.Printf("metrics: query total_tokens: %v", err)
		totalTokens = 0
	}

	if err := c.db.Conn().QueryRow("SELECT COALESCE(SUM(count), 0) FROM quota_data").Scan(&totalRequests); err != nil {
		log.Printf("metrics: query total_requests: %v", err)
		totalRequests = 0
	}

	ch <- prometheus.MustNewConstMetric(c.userCount, prometheus.GaugeValue, userCount)
	ch <- prometheus.MustNewConstMetric(c.totalTokens, prometheus.GaugeValue, totalTokens)
	ch <- prometheus.MustNewConstMetric(c.totalRequests, prometheus.GaugeValue, totalRequests)

	if !c.db.HasLogs() {
		return
	}

	var logTokens, logQuota float64
	table := c.db.LogTable()

	tps, err := c.db.TokenTPS()
	if err != nil {
		log.Printf("metrics: query token_tps: %v", err)
		tps = 0
	}
	if err := c.db.LogConn().QueryRow(fmt.Sprintf("SELECT COALESCE(SUM(prompt_tokens + completion_tokens), 0) FROM %s", table)).Scan(&logTokens); err != nil {
		log.Printf("metrics: query log_tokens: %v", err)
		logTokens = 0
	}
	if err := c.db.LogConn().QueryRow(fmt.Sprintf("SELECT COALESCE(SUM(quota), 0) FROM %s", table)).Scan(&logQuota); err != nil {
		log.Printf("metrics: query log_quota: %v", err)
		logQuota = 0
	}

	ch <- prometheus.MustNewConstMetric(c.tokenTPS, prometheus.GaugeValue, tps)
	ch <- prometheus.MustNewConstMetric(c.logTokens, prometheus.GaugeValue, logTokens)
	ch <- prometheus.MustNewConstMetric(c.logQuota, prometheus.GaugeValue, logQuota)
}
