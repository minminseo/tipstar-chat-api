package usecase_test

import (
	"testing"
	"time"

	"github.com/minminseo/tipstar-chat-api/domain"
	"github.com/minminseo/tipstar-chat-api/usecase"
)

type mockDeleteRepo struct {
	FindByIDFn func(domain.MessageID) (*domain.Message, error)
	DeleteFn   func(*domain.Message) error
}

func (m *mockDeleteRepo) FindByID(id domain.MessageID) (*domain.Message, error) {
	return m.FindByIDFn(id)
}
func (m *mockDeleteRepo) Insert(_ *domain.Message) error   { return nil }
func (m *mockDeleteRepo) Update(_ *domain.Message) error   { return nil }
func (m *mockDeleteRepo) Delete(msg *domain.Message) error { return m.DeleteFn(msg) }

func TestDeleteMessageUsecase_Execute(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		message     *domain.Message
		findErr     error
		deleteErr   error
		userID      string
		expectError bool
	}{
		{
			name: "正常系",
			message: &domain.Message{
				ID:        domain.MessageID("m1"),
				UserID:    domain.UserID("u1"),
				CreatedAt: now,
			},
			userID:      "u1",
			expectError: false,
		},
		{
			name:        "見つからない",
			message:     nil,
			findErr:     nil,
			userID:      "u1",
			expectError: true,
		},
		{
			name: "他人が削除",
			message: &domain.Message{
				ID:     domain.MessageID("m1"),
				UserID: domain.UserID("u1"),
			},
			userID:      "u2",
			expectError: true,
		},
		{
			name: "既に削除済み",
			message: &domain.Message{
				ID:        domain.MessageID("m1"),
				UserID:    domain.UserID("u1"),
				DeletedAt: &now,
			},
			userID:      "u1",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDeleteRepo{
				FindByIDFn: func(_ domain.MessageID) (*domain.Message, error) {
					return tt.message, tt.findErr
				},
				DeleteFn: func(_ *domain.Message) error {
					return tt.deleteErr
				},
			}

			uc := usecase.NewDeleteMessageUsecase(mock)
			err := uc.Execute(domain.MessageID("m1"), domain.UserID(tt.userID))

			if tt.expectError && err == nil {
				t.Errorf("期待したエラーが返されなかった")
			}
			if !tt.expectError && err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}
		})
	}
}
