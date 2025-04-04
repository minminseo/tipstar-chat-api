package domain

import (
	"time"
)

type MessageID string
type TipID string
type UserID string

type Message struct {
	ID        MessageID
	TipID     TipID
	UserID    UserID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	IsAuthor  bool
}
