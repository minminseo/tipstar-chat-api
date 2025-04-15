package usecase

// WebSocket経由のリクエストに対するユースケース

/*
ここに実装されている3つのメソッドの処理の流れ
1. プレゼンテーション層の/websocketのパッケージでwebsocket経由で受信したメッセージを引数として受け取る
2. 必要な処理（ドメイン層で定義されているビジネスロジック）を施す
3. ドメイン層にある永続化処理系のインターフェースに定義されているメソッドを呼び出す

このユースケース層の依存先であるドメイン層の「永続化処理メソッドが定義されているインターフェース」の具体的な実装はインフラ層で行う。

*/

import (
	"context"
	"errors"
	"log"

	"github.com/minminseo/tipstar-chat-api/domain"
)

type onlyWSMessageUseCase struct {
	repo domain.MessageRepository
}

// 永続化処理のインターフェースのメソッドをユースケース層に依存注入するコンストラクタ関数
func NewOnlyWSMessageUseCase(repo domain.MessageRepository) OnlyWSUsecase {
	//明示的にフィールドrepoに引数repo（インターフェース）を代入して依存注入（ドメイン層の永続化処理専門のインターフェースのメソッドを渡す）
	return &onlyWSMessageUseCase{repo: repo}
}

// メッセージ送信のユースケース
func (uc *onlyWSMessageUseCase) ExecuteSendMessage(ctx context.Context, msg *domain.Message) error {

	return uc.repo.SaveMessage(msg)
}

// メッセージ編集のユースケース
func (uc *onlyWSMessageUseCase) EditMessage(ctx context.Context, messageID domain.MessageID, userID domain.UserID, newContent string) error {
	msg, err := uc.repo.FetchMessageByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("メッセージが見つかりません")
	}
	if err := msg.SetEditedContent(userID, newContent); err != nil {
		return err
	}
	log.Printf("ContentとUpdatedAtの実体書き換え成功（永続化前）")

	// ドメイン層の永続化処理系のインターフェースに定義されている編集系のメソッドを呼び出す。具体的な実装はインフラ層で行う。
	return uc.repo.Update(ctx, msg)
}

// メッセージ論理削除のユースケース
func (uc *onlyWSMessageUseCase) DeleteMessage(ctx context.Context, messageID domain.MessageID, userID domain.UserID) error {
	msg, err := uc.repo.FetchMessageByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("削除対象のメッセージが見つかりません")
	}
	if err := msg.SetDeletedMessage(userID); err != nil {
		return err
	}
	log.Printf("DeletedAtの実体書き換え成功（永続化前）")

	// ドメイン層の永続化処理系のインターフェースに定義されている論理削除メソッドを呼び出す。具体的な実装はインフラ層で行う。
	return uc.repo.SoftDelete(ctx, msg)
}
