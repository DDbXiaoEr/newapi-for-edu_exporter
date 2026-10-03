package models

type Statistics struct {
	UserCount     int64        `json:"user_count"`
	TotalTokens   int64        `json:"total_tokens"`
	TotalRequests int64        `json:"total_requests"`
	TokenTPS      float64      `json:"token_tps"`
	LogTokens     int64        `json:"log_tokens"`
	LogQuota      int64        `json:"log_quota"`
	ChannelTPS    []ChannelTPS `json:"channel_tps,omitempty"`
}

type ChannelTPS struct {
	ChannelID   int64   `json:"channel_id"`
	ChannelName string  `json:"channel_name"`
	TPS         float64 `json:"tps"`
}
