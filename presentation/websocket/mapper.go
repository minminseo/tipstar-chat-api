package websocket

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/minminseo/tipstar-chat-api/domain"
)

// ユーザーIDは、WSRequestMessage 内のものではなく、接続時に取得した conn.UserID を使用する。
func ToSendDomainFromWSRequest(req *WSRequestMessage, connUserID string) (*domain.Message, error) {
	return domain.NewMessage(
		domain.MessageID(generateUUID()), // 新規送信なので新たに生成
		domain.TipID(req.TipID),
		domain.UserID(connUserID), // ヘッダーからのユーザーIDを使用
		req.Content,

		true, // 送信者は自分としてフラグを立てる
	)
}

// user_id は引数 connUserID から取得します。MessageID 必須。
func ToEditDomainFromWSRequest(req *WSRequestMessage, connUserID string) (*domain.Message, error) {
	if req.MessageID == "" {
		return nil, fmt.Errorf("message_idが編集リクエストに含まれていません")
	}
	return &domain.Message{
		ID:      domain.MessageID(req.MessageID),
		TipID:   domain.TipID(req.TipID),
		UserID:  domain.UserID(connUserID), // ヘッダーからのユーザーIDを使用
		Content: req.Content,               // 新しい内容
		// CreatedAt, UpdatedAt, DeletedAt は usecase で既存メッセージを取得して更新する前提
	}, nil
}

// user_id は引数 connUserID から取得します。MessageID 必須。
func ToDeleteDomainFromWSRequest(req *WSRequestMessage, connUserID string) (*domain.Message, error) {
	if req.MessageID == "" {
		return nil, fmt.Errorf("message_idが削除リクエストに含まれていません")
	}
	return &domain.Message{
		ID:     domain.MessageID(req.MessageID),
		TipID:  domain.TipID(req.TipID),
		UserID: domain.UserID(connUserID), // ヘッダーからのユーザーIDを使用
		// Content は編集と異なり不要
	}, nil
}

func ToBroadcastMessage(msg *domain.Message) *WSBroadcastMessage {
	var ts int64 = msg.CreatedAt.Unix()
	return &WSBroadcastMessage{
		Type:      "send",
		MessageID: string(msg.ID),
		TipID:     string(msg.TipID),
		Content:   msg.Content,
		Timestamp: ts,
	}
}

func ToEditBroadcastMessage(msg *domain.Message) *EditBroadcastMessage {
	return &EditBroadcastMessage{
		Type:       "edit",
		MessageID:  string(msg.ID),
		TipID:      string(msg.TipID),
		NewContent: msg.Content, // 編集後の内容。必要に応じて更新済みの値を利用
		EditedAt:   msg.UpdatedAt.Unix(),
	}
}

func ToDeleteBroadcastMessage(msg *domain.Message) *DeleteBroadcastMessage {
	var deletedAt int64
	if msg.DeletedAt != nil {
		deletedAt = msg.DeletedAt.Unix()
	}
	return &DeleteBroadcastMessage{
		Type:      "delete",
		MessageID: string(msg.ID),
		TipID:     string(msg.TipID),
		DeletedAt: deletedAt,
	}
}

func generateUUID() string {
	return uuid.New().String()
}
