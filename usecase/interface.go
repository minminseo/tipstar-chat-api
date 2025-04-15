package usecase

import (
	"context"

	"github.com/minminseo/tipstar-chat-api/domain"
)

// ユースケースの実装をインターフェースとして定義し、プレゼンテーション層はこのインターフェースだけに依存しユースケース内部の実装を隠す。

// HTTP経由（Rest API）のリクエスト用のユースケース
type OnlyRestUsecase interface {
	GetAllMessages(ctx context.Context, tipID string) ([]*domain.Message, error)
}

// Websocket経由のリクエストのユースケース
type OnlyWSUsecase interface {
	ExecuteSendMessage(ctx context.Context, msg *domain.Message) error
	EditMessage(ctx context.Context, messageID domain.MessageID, userID domain.UserID, newContent string) error
	DeleteMessage(ctx context.Context, messageID domain.MessageID, userID domain.UserID) error
}
