package db

// ドメイン層で定義した、インフラ層のメソッド（ここ）を持つインターフェースをここで実装

import (
	"context"
	"errors"
	"log"

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
func (r *PgxMessageRepository) FetchMessageByID(ctx context.Context, id domain.MessageID) (*domain.Message, error) {
	const query = `
		SELECT id, tip_id, user_id, content, created_at, updated_at, deleted_at
		FROM messages
		WHERE id = $1
	`

	var m MessageModel
	err := r.DB.QueryRow(ctx, query, string(id)).Scan(
		&m.ID,
		&m.TipID,
		&m.UserID,
		&m.Content,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return ToDomainModel(&m, false), nil // 同時にドメインモデル構造体に変換
}

// 以下永続化処理
// DBモデル構造体→ドメインモデル構造体へのマッピング、その逆のマッピングは変換関数を使用（/infra/db/mapper.goに定義）

// メッセージの挿入（ユースケース的にはメッセージ送信）。これを使うコードの実装は後回しだけどインターフェースの条件満たすために先に定義。
func (r *PgxMessageRepository) SaveMessage(msg *domain.Message) error {
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
		dbMsg.UpdatedAt)
	log.Printf("メッセージ送信（永続化）")
	return err
}

// メッセージの編集。編集対象のメッセージがなければエラー返す
func (r *PgxMessageRepository) Update(ctx context.Context, msg *domain.Message) error {
	const query = `
	UPDATE messages
	SET content = $1, updated_at = $2
	WHERE id = $3
	`
	dbMsg := ToDbModel(msg)
	tag, err := r.DB.Exec(ctx, query, dbMsg.Content, dbMsg.UpdatedAt, dbMsg.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("対象メッセージが見つかりません")
	}
	log.Printf("メッセージ編集（永続化）")
	return nil
}

// メッセージの論理削除（deleted_atを設定）。削除対象のメッセージなければエラー返す
func (r *PgxMessageRepository) SoftDelete(ctx context.Context, msg *domain.Message) error {
	const query = `
	UPDATE messages
	SET deleted_at = $1
	WHERE id = $2
	`
	dbMsg := ToDbModel(msg)
	tag, err := r.DB.Exec(ctx, query, dbMsg.DeletedAt, dbMsg.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("削除対象のメッセージが見つかりません")
	}
	log.Printf("メッセージ削除（永続化）")
	return nil
}

// tip_idに紐づくメッセージの一覧をcreatedAtの昇順で取得。
func (r *PgxMessageRepository) GetAllMessages(tipID domain.TipID) ([]*domain.Message, error) {
	const query = `
	SELECT id, tip_id, user_id, content, created_at, updated_at, deleted_at
	FROM messages
	WHERE tip_id = $1
	ORDER BY created_at ASC
	`
	rows, err := r.DB.Query(context.Background(), query, string(tipID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []*domain.Message
	for rows.Next() {
		var m MessageModel
		if err := rows.Scan(&m.ID, &m.TipID, &m.UserID, &m.Content, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt); err != nil {
			return nil, err
		}
		message := ToDomainModel(&m, false)
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Printf("メッセージ一覧取得")
	return messages, nil
}
