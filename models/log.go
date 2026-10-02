package models

type Log struct {
	ID                int64  `json:"id"`
	UserID            int64  `json:"user_id"`
	CreatedAt         int64  `json:"created_at"`
	Type              int64  `json:"type"`
	Content           string `json:"content"`
	Username          string `json:"username"`
	TokenName         string `json:"token_name"`
	ModelName         string `json:"model_name"`
	Quota             int64  `json:"quota"`
	PromptTokens      int64  `json:"prompt_tokens"`
	CompletionTokens  int64  `json:"completion_tokens"`
	UseTime           int64  `json:"use_time"`
	IsStream          uint8  `json:"is_stream"`
	ChannelID         int64  `json:"channel_id"`
	TokenID           int64  `json:"token_id"`
	Group             string `json:"group"`
	IP                string `json:"ip"`
	RequestID         string `json:"request_id"`
	UpstreamRequestID string `json:"upstream_request_id"`
	Other             string `json:"other"`
}
