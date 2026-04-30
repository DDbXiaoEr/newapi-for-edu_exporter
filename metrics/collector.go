package metrics

import (
	"database/sql"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	db *sql.DB

	userCount     *prometheus.Desc
	totalTokens   *prometheus.Desc
	totalRequests *prometheus.Desc
}

func NewCollector(db *sql.DB) *Collector {
	return &Collector{
		db: db,
		userCount: prometheus.NewDesc(
			"newapi_user_count",
			"Total number of active users",
			nil, nil,
		),
		totalTokens: prometheus.NewDesc(
			"newapi_total_tokens",
			"Total tokens used across all models",
			nil, nil,
		),
		totalRequests: prometheus.NewDesc(
			"newapi_total_requests",
			"Total API requests across all models",
			nil, nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.userCount
	ch <- c.totalTokens
	ch <- c.totalRequests
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var userCount, totalTokens, totalRequests float64

	if err := c.db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&userCount); err != nil {
		log.Printf("metrics: query user_count: %v", err)
		userCount = 0
	}

	if err := c.db.QueryRow("SELECT COALESCE(SUM(token_used), 0) FROM quota_data").Scan(&totalTokens); err != nil {
		log.Printf("metrics: query total_tokens: %v", err)
		totalTokens = 0
	}

	if err := c.db.QueryRow("SELECT COALESCE(SUM(count), 0) FROM quota_data").Scan(&totalRequests); err != nil {
		log.Printf("metrics: query total_requests: %v", err)
		totalRequests = 0
	}

	ch <- prometheus.MustNewConstMetric(c.userCount, prometheus.GaugeValue, userCount)
	ch <- prometheus.MustNewConstMetric(c.totalTokens, prometheus.GaugeValue, totalTokens)
	ch <- prometheus.MustNewConstMetric(c.totalRequests, prometheus.GaugeValue, totalRequests)
}
