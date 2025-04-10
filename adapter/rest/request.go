package rest

// メッセージ編集用リクエスト構造体
type EditMessageRequest struct {
	MessageID  string `json:"message_id"`  // 編集対象のメッセージID
	NewContent string `json:"new_content"` // 新しいメッセージ内容
}

// メッセージ削除用リクエスト構造体
type DeleteMessageRequest struct {
	MessageID string `json:"message_id"`

	// user_id（将来的にはJWTトークンはヘッダーからとりだすのでUserIDはなし）
}
