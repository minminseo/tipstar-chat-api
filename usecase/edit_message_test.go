package usecase_test

import (
	"testing"
	"time"

	"github.com/minminseo/tipstar-chat-api/domain"
	"github.com/minminseo/tipstar-chat-api/usecase"
)

type mockEditRepo struct {
	FindByIDFn func(domain.MessageID) (*domain.Message, error)
	UpdateFn   func(*domain.Message) error
}

func (m *mockEditRepo) FindByID(id domain.MessageID) (*domain.Message, error) {
	return m.FindByIDFn(id)
}

func (m *mockEditRepo) Insert(_ *domain.Message) error { return nil }
func (m *mockEditRepo) Delete(_ *domain.Message) error { return nil }
func (m *mockEditRepo) Update(msg *domain.Message) error {
	return m.UpdateFn(msg)
}

func TestEditMessageUseCase_Execute(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		message     *domain.Message
		findErr     error
		updateErr   error
		userID      string
		newContent  string
		expectError bool
	}{
		{
			name: "正常系",
			message: &domain.Message{
				ID:        "m1",
				UserID:    "u1",
				Content:   "old",
				CreatedAt: now,
				UpdatedAt: now,
			},
			userID:      "u1",
			newContent:  "new",
			expectError: false,
		},
		{
			name:        "見つからない",
			message:     nil,
			findErr:     nil,
			userID:      "u1",
			newContent:  "new",
			expectError: true,
		},
		{
			name: "他人による編集",
			message: &domain.Message{
				ID:     "m1",
				UserID: "u1",
			},
			userID:      "u2",
			newContent:  "new",
			expectError: true,
		},
		{
			name: "空文字編集",
			message: &domain.Message{
				ID:     "m1",
				UserID: "u1",
			},
			userID:      "u1",
			newContent:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockEditRepo{
				FindByIDFn: func(_ domain.MessageID) (*domain.Message, error) {
					return tt.message, tt.findErr
				},
				UpdateFn: func(_ *domain.Message) error {
					return tt.updateErr
				},
			}

			uc := usecase.NewEditMessageUseCase(mock)
			err := uc.Execute("m1", tt.userID, tt.newContent)

			if tt.expectError && err == nil {
				t.Errorf("期待したエラーが返されなかった")
			}
			if !tt.expectError && err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}
		})
	}
}
