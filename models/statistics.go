package models

type Statistics struct {
	UserCount     int64 `json:"user_count"`
	TotalTokens   int64 `json:"total_tokens"`
	TotalRequests int64 `json:"total_requests"`
}
