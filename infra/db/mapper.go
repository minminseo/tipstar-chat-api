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
func ToDbModel(m *domain.Message) *MessageModel {
	return &MessageModel{
		ID:        string(m.ID),
		TipID:     string(m.TipID),
		UserID:    string(m.UserID),
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}
