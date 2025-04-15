package domain

import (
	"context"
)

// メッセージの永続化処理のメソッドを定義するインターフェース
// ユースケース層が依存する用のインターフェース
// 具体的な実装はインフラ層で行う
type MessageRepository interface {
	FetchMessageByID(ctx context.Context, id MessageID) (*Message, error) // クライアントからきたMessageIDを元にDBからメッセージを取得するメソッド
	SaveMessage(msg *Message) error                                       // メッセージをDBに挿入するメソッド
	Update(ctx context.Context, msg *Message) error                       // メッセージを編集するメソッド
	SoftDelete(ctx context.Context, msg *Message) error                   // メッセージを論理削除するメソッド
	GetAllMessages(tipID TipID) ([]*Message, error)                       // tipIDでに対応するチャット履歴を一覧取得する。
}
