package rest

import "github.com/minminseo/tipstar-chat-api/domain"

// リクエスト構造体をドメインモデル（編集用）に変換する関数
func ToEditDomainModel(req *EditMessageRequest, userID string) *domain.Message {
	return &domain.Message{
		ID:      domain.MessageID(req.MessageID),
		UserID:  domain.UserID(userID),
		Content: req.NewContent,
	}
}

// リクエスト構造体をドメインモデル（削除用）に変換する関数
func ToDeleteDomainModel(req *DeleteMessageRequest, userID string) *domain.Message {
	return &domain.Message{
		ID:     domain.MessageID(req.MessageID),
		UserID: domain.UserID(userID),
	}
}

// ドメインモデルをレスポンス構造体（編集用）に変換する関数
func ToEditMessageResponse(msgID domain.MessageID) *EditMessageResponse {
	return &EditMessageResponse{
		Message:   "メッセージを編集しました",
		MessageID: string(msgID),
	}
}

// ドメインモデルをレスポンス構造体（削除用）に変換する関数
func ToDeleteMessageResponse(msgID domain.MessageID) *DeleteMessageResponse {
	return &DeleteMessageResponse{
		Message:   "メッセージを削除しました",
		MessageID: string(msgID),
	}
}
