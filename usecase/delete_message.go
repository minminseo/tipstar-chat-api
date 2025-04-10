package usecase

import (
	"errors"

	"github.com/minminseo/tipstar-chat-api/domain"
)

type DeleteMessageUseCase struct {
	repo domain.MessageRepository
}

// インフラ層のインターフェースを満たすDB操作に関する実装をユースケース層に依存注入するコンストラクタ関数（インフラ層に依存せずここで使えるようにする）
func NewDeleteMessageUsecase(repo domain.MessageRepository) *DeleteMessageUseCase {

	// 明示的にフィールドrepoに引数repo（インターフェース）を代入して依存注入（インフラ層の実装をユースケースに層に渡す。）
	return &DeleteMessageUseCase{repo: repo}
}

// インフラ層のインターフェース（実際はドメイン層にインターフェースは定義）越しでDB操作をする関数
func (uc *DeleteMessageUseCase) Execute(messageID domain.MessageID, userID domain.UserID) error {

	// クライアントから渡されたmessageIDをドメイン型（domain.MessageID）にキャストし、メッセージを取得するFindByID関数（抽象）に引数として渡す
	msg, err := uc.repo.FindByID(messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("削除対象のメッセージが見つかりません")
	}

	// クライアントから渡されたuserIDをドメイン型（domain.UserID）にキャストし、メッセージのDeletedAtを更新する（論理削除）Delete関数に引数として渡す
	// SupabaseのDBと分離してて、Renderの方にはmessagesテーブルしか作っていないので整合性は全くなし状態
	if err := msg.Delete(userID); err != nil {
		return err
	}

	// EditContent関数によって内容が更新されたMessageのインスタンス（msg）を、更新を反映するUpdate関数（抽象）に引数として渡し永続化処理を依頼。
	if err := uc.repo.Delete(msg); err != nil {
		return err
	}

	return nil
}
