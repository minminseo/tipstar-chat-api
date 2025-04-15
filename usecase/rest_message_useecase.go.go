package usecase

// HTTP経由（Rest API）のリクエストに対するユースケース

/*
ここに実装されているメソッドの処理の流れ
1. プレゼンテーション層の/restのパッケージでHTTP（RoomRest API）経由のリクエストで受け取ったtipIDを引数として受け取る
2. ドメイン層にある永続化処理系のインターフェースに定義されているメソッドを呼び出す（現状GetAllMessagesはただの取得処理だからインターフェース分離した方が良いかもしれない）

このユースケース層の依存先であるドメイン層の「永続化処理メソッドが定義されているインターフェース」の具体的な実装はインフラ層で行う。

*/
import (
	"context"

	"github.com/minminseo/tipstar-chat-api/domain"
)

type onlyRestMessageUseCase struct {
	repo domain.MessageRepository
}

// 永続化処理のインターフェースのメソッドをユースケース層に依存注入するコンストラクタ関数
func NewOnlyRestMessageUseCase(repo domain.MessageRepository) OnlyRestUsecase {
	return &onlyRestMessageUseCase{repo: repo}
}

// メッセージ一覧取得のユースケース
func (uc *onlyRestMessageUseCase) GetAllMessages(ctx context.Context, tipID string) ([]*domain.Message, error) {
	return uc.repo.GetAllMessages(domain.TipID(tipID))
}
