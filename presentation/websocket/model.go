package websocket

// WSRequestMessage は、クライアントから送信されるWebSocketリクエストメッセージのモデルです。
// 新規送信、編集、削除いずれの場合も、この形式で受信します。
// 例：type "send", "edit", "delete"
type WSRequestMessage struct {
	Type      string `json:"type"`       // "send", "edit", "delete"
	MessageID string `json:"message_id"` // 新規の場合は空。編集・削除の場合は既存のID
	TipID     string `json:"tip_id"`     // 対象チャットルームのID
	Content   string `json:"content"`    // メッセージ内容（送信の場合はメッセージ全文、編集の場合は新しい内容。削除では無視）
	UserID    string `json:"user_id"`    // クライアントから送信されるユーザーID
}

// WSBroadcastMessage は、サーバーがクライアントに送信するWebSocketレスポンスの基本モデルです。
type WSBroadcastMessage struct {
	Type      string `json:"type"`       // "send", "catchup" など（新規送信やキャッチアップ用）
	MessageID string `json:"message_id"` // メッセージID
	TipID     string `json:"tip_id"`     // 対象チャットルームのID
	Content   string `json:"content"`    // メッセージ内容
	Timestamp int64  `json:"timestamp"`  // Unixタイムスタンプ（作成時刻）
}

// --- 以下、編集と削除のブロードキャスト用の構造体 ---

// EditBroadcastMessage は、編集結果を WebSocket ブロードキャストする際に使用するモデルです。
type EditBroadcastMessage struct {
	Type       string `json:"type"`        // 固定で "edit"
	MessageID  string `json:"message_id"`  // 編集対象のメッセージID
	TipID      string `json:"tip_id"`      // チャットルームのID
	NewContent string `json:"new_content"` // 編集後の新しい内容
	EditedAt   int64  `json:"edited_at"`   // Unix タイムスタンプ（更新時刻）
}

// DeleteBroadcastMessage は、削除結果を WebSocket ブロードキャストする際に使用するモデルです。
type DeleteBroadcastMessage struct {
	Type      string `json:"type"`       // 固定で "delete"
	MessageID string `json:"message_id"` // 削除対象のメッセージID
	TipID     string `json:"tip_id"`     // チャットルームのID
	DeletedAt int64  `json:"deleted_at"` // Unix タイムスタンプ（削除時刻）
}
