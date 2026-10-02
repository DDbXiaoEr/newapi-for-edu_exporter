package models

type QuotaData struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	ModelName string `json:"model_name"`
	CreatedAt int64  `json:"created_at"`
	UseGroup  string `json:"use_group"`
	TokenID   int64  `json:"token_id"`
	ChannelID int64  `json:"channel_id"`
	NodeName  string `json:"node_name"`
	TokenUsed int64  `json:"token_used"`
	Count     int64  `json:"count"`
	Quota     int64  `json:"quota"`
}
