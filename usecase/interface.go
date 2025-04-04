package usecase

// ユースケースの実装をインターフェースとして定義し、ハンドラ層はこのインターフェースだけに依存してユースケース層の実装に依存しないようにする。

type EditMessageInput interface {
	Execute(messageID string, userID string, newContent string) error
}

type DeleteMessageInput interface {
	Execute(messageID string, userID string) error
}
