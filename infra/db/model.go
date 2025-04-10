package db

import "time"

// DBモデル構造体定義
type MessageModel struct {
	ID        string     // messages.id（UUID）←PK
	TipID     string     // messages.tip_id（UUID）←NOT NULL制約
	UserID    string     // messages.user_id（UUID）←NOT NULL制約
	Content   string     // messages.content（TEXT）←NOT NULL制約
	CreatedAt time.Time  // messages.created_at（TIMESTAMP） ←NOT NULL制約
	UpdatedAt time.Time  // messages.updated_at（TIMESTAMP） ←NOT NULL制約（初期値はcreated_atと同じにする）
	DeletedAt *time.Time // messages.deleted_at（TIMESTAMP） ←NULL許容（論理削除したいから）
}
