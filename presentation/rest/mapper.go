package rest

import (
	"github.com/minminseo/tipstar-chat-api/domain"
)

// ToChatMessageResponse converts a domain.Message to ChatMessageResponse.
func ToChatMessageResponse(msg *domain.Message) *ChatMessageResponse {
	return &ChatMessageResponse{
		MessageID: string(msg.ID),
		TipID:     string(msg.TipID),
		UserID:    string(msg.UserID),
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt.Unix(),
		UpdatedAt: msg.UpdatedAt.Unix(),
	}
}

// ToChatMessagesResponse converts a slice of domain.Message to a slice of ChatMessageResponse.
func ToChatMessagesResponse(msgs []*domain.Message) []*ChatMessageResponse {
	res := make([]*ChatMessageResponse, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, ToChatMessageResponse(m))
	}
	return res
}
