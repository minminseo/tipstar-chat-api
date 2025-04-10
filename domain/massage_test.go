package domain_test

import (
	"testing"
	"time"

	"github.com/minminseo/tipstar-chat-api/domain"
)

func TestCanEdit(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		userID   domain.UserID
		msgUser  domain.UserID
		deleted  bool
		expected bool
	}{
		{"本人で削除前", "user1", "user1", false, true},
		{"本人で削除後", "user1", "user1", true, false},
		{"他人で削除前", "user2", "user1", false, false},
		{"他人で削除後", "user2", "user1", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var deletedAt *time.Time
			if tt.deleted {
				deletedAt = &now
			}

			msg := &domain.Message{
				UserID:    tt.msgUser,
				DeletedAt: deletedAt,
			}

			got := msg.CanEdit(tt.userID)
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestEditContent(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		deleted     bool
		newContent  string
		expectError bool
	}{
		{"正常系", false, "新しい内容", false},
		{"削除済み", true, "新しい内容", true},
		{"空文字", false, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var deletedAt *time.Time
			if tt.deleted {
				deletedAt = &now
			}

			msg := &domain.Message{
				Content:   "元の内容",
				UpdatedAt: now,
				DeletedAt: deletedAt,
			}

			err := msg.EditContent(tt.newContent)
			if tt.expectError && err == nil {
				t.Errorf("期待したエラーが返されなかった")
			}
			if !tt.expectError && err != nil {
				t.Errorf("想定外のエラー: %v", err)
			}
			if !tt.expectError && msg.Content != tt.newContent {
				t.Errorf("内容が更新されていない: got %v, want %v", msg.Content, tt.newContent)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		msgUserID   domain.UserID
		requestUser domain.UserID
		deleted     bool
		expectError bool
	}{
		{"正常系：本人で未削除", "user1", "user1", false, false},
		{"異常系：他人のメッセージ", "user1", "user2", false, true},
		{"異常系：すでに削除済み", "user1", "user1", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var deletedAt *time.Time
			if tt.deleted {
				deletedAt = &now
			}

			msg := &domain.Message{
				UserID:    tt.msgUserID,
				DeletedAt: deletedAt,
			}

			err := msg.Delete(tt.requestUser)

			if tt.expectError && err == nil {
				t.Errorf("期待したエラーが返されなかった")
			}
			if !tt.expectError && err != nil {
				t.Errorf("想定外のエラーが返された: %v", err)
			}

			// 削除成功時のみ、DeletedAt が更新されていることを確認
			if !tt.expectError && msg.DeletedAt == nil {
				t.Errorf("削除が成功したはずなのに DeletedAt が更新されていません")
			}
		})
	}
}
