package db

// ドメイン層で定義した、インフラ層のメソッド（ここ）を持つインターフェースをここで実装

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minminseo/tipstar-chat-api/domain"
)

type PgxMessageRepository struct {
	DB *pgxpool.Pool
}

func NewPgxMessageRepository(db *pgxpool.Pool) domain.MessageRepository {
	return &PgxMessageRepository{DB: db}
}

// メッセージをIDで取得する（論理削除も含めて）
func (r *PgxMessageRepository) FindByID(id domain.MessageID) (*domain.Message, error) {
	const query = `
		SELECT id, tip_id, user_id, content, created_at, updated_at, deleted_at
		FROM messages
		WHERE id = $1
	`

	var m MessageModel
	err := r.DB.QueryRow(context.Background(), query, string(id)).Scan(
		&m.ID, &m.TipID, &m.UserID, &m.Content, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return ToDomainModel(&m, false), nil // 同時にドメインモデル構造体に変換
}

// 以下永続化処理
// DBモデル構造体→ドメインモデル構造体へのマッピング、その逆のマッピングは変換関数を使用（/infra/db/mapper.goに定義）

// メッセージの挿入（ユースケース的にはメッセージ送信）。これを使うコードの実装は後回しだけどインターフェースの条件満たすために先に定義。
func (r *PgxMessageRepository) Insert(msg *domain.Message) error {
	const query = `
		INSERT INTO messages (id, tip_id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	dbMsg := ToDbModel(msg)

	_, err := r.DB.Exec(context.Background(), query,
		dbMsg.ID,
		dbMsg.TipID,
		dbMsg.UserID,
		dbMsg.Content,
		dbMsg.CreatedAt,
		dbMsg.UpdatedAt,
	)
	return err
}

// メッセージの編集（PUT）。編集対象のメッセージがなければエラー返す
func (r *PgxMessageRepository) Update(msg *domain.Message) error {
	const query = `
		UPDATE messages
		SET content = $1, updated_at = $2
		WHERE id = $3
	`
	dbMsg := ToDbModel(msg)

	tag, err := r.DB.Exec(context.Background(), query,
		dbMsg.Content,
		dbMsg.UpdatedAt,
		dbMsg.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("対象メッセージが見つかりません")
	}
	return nil
}

// メッセージの論理削除（deleted_at を設定）。削除対象のメッセージなければエラー返す
func (r *PgxMessageRepository) Delete(msg *domain.Message) error {
	const query = `
		UPDATE messages
		SET deleted_at = $1
		WHERE id = $2
	`

	dbMsg := ToDbModel(msg)

	tag, err := r.DB.Exec(context.Background(), query,
		dbMsg.DeletedAt,
		dbMsg.ID,
	)

	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("削除対象のメッセージが見つかりません")
	}

	fmt.Println("確認用→ DeletedAt:", dbMsg.DeletedAt)
	return nil
}
