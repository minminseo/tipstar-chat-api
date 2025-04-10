package usecase

import "github.com/minminseo/tipstar-chat-api/domain"

// ユースケースの実装をインターフェースとして定義し、ハンドラ層はこのインターフェースだけに依存してユースケース層の実装に依存しないようにする。

type EditMessageInput interface {
	Execute(messageID domain.MessageID, userID domain.UserID, newContent string) error
}

type DeleteMessageInput interface {
	Execute(messageID domain.MessageID, userID domain.UserID) error
}
