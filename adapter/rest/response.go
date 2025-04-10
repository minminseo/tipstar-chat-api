package rest

// レスポンス構造体を定義

// メッセージ編集成功時のレスポンス
type EditMessageResponse struct {
	Message   string `json:"message"`    // 編集成功メッセージ
	MessageID string `json:"message_id"` // 編集したメッセージID
}

// メッセージ削除成功時のレスポンス
type DeleteMessageResponse struct {
	Message   string `json:"message"`    // 削除成功メッセージ
	MessageID string `json:"message_id"` // 削除したメッセージID
}
