package rest

// ChatMessageResponse は、REST APIで返すチャットメッセージのレスポンス形式です。
type ChatMessageResponse struct {
	MessageID string `json:"message_id"`
	TipID     string `json:"tip_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"` // Unix timestamp
	UpdatedAt int64  `json:"updated_at"` // Unix timestamp
}
