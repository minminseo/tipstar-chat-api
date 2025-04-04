package usecase

import (
	"errors"

	"github.com/minminseo/tipstar-chat-api/domain"
)

type EditMessageUseCase struct {
	repo domain.MessageRepository
}

// インフラ層のインターフェースを満たすDB操作に関する実装をユースケース層に依存注入するコンストラクタ関数
func NewEditMessageUseCase(repo domain.MessageRepository) *EditMessageUseCase {

	// 明示的にフィールドrepoに引数repo（インターフェース）を代入して依存注入（インフラ層の実装をユースケースに層に渡す。）
	return &EditMessageUseCase{repo: repo}
}

// インフラ層のインターフェース越しでDB操作をする関数
func (uc *EditMessageUseCase) Execute(messageID string, userID string, newContent string) error {

	// クライアントから渡されたmessageIDをドメイン型（domain.MessageID）にキャストし、メッセージを取得するFindByID関数（抽象）に引数として渡す
	msg, err := uc.repo.FindByID(domain.MessageID(messageID))
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("メッセージが見つかりません")
	}

	// クライアントから渡されたuserIDをドメイン型（domain.UserID）にキャストし、メッセージを編集する権限があるか判別するCanEdit関数に引数として渡す
	if !msg.CanEdit(domain.UserID(userID)) {
		return errors.New("このメッセージを編集する権限がありません")
	}

	// クライアントから渡されたnewContentをメッセージ編集処理のみを行うEditContent関数に引数として渡す
	if err := msg.EditContent(newContent); err != nil {
		return err
	}

	// EditContent関数によって内容が更新されたMessageのインスタンス（msg）を、更新を反映するUpdate関数（抽象）に引数として渡し永続化処理を依頼。
	if err := uc.repo.Update(msg); err != nil {
		return err
	}

	return nil
}
