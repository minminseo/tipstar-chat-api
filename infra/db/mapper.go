package db

import (
	"github.com/minminseo/tipstar-chat-api/domain"
)

// DB構造体をドメインモデル構造体に変換する関数
func ToDomainModel(m *MessageModel, isAuthor bool) *domain.Message {
	return &domain.Message{
		ID:        domain.MessageID(m.ID),
		TipID:     domain.TipID(m.TipID),
		UserID:    domain.UserID(m.UserID),
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
		IsAuthor:  isAuthor,
	}
}

// ドメインモデル構造体をDBモデル構造体に変換する関数
func ToDbModel(msg *domain.Message) *MessageModel {
	return &MessageModel{
		ID:        string(msg.ID),
		TipID:     string(msg.TipID),
		UserID:    string(msg.UserID),
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		DeletedAt: msg.DeletedAt,
	}
}
