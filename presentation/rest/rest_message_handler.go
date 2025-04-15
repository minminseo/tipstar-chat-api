package rest

// ここではHTTP経由（Rest API）のリクエストのハンドリングを行う（最終的にJSON形式にしてwに書き込む）

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/minminseo/tipstar-chat-api/usecase"
)

type OnlyRestMessageHandler struct {
	uc usecase.OnlyRestUsecase
}

// ユースケースのインターフェースを満たす実装プレゼンテーションそうに注入するコンストラクタ関数（ユースケース内部の処理を隠してここで使えるようにする）
func NewOnlyRestMessageHandler(uc usecase.OnlyRestUsecase) *OnlyRestMessageHandler {
	return &OnlyRestMessageHandler{uc: uc}
}

// チャット履歴一覧取得のハンドラー
func (h *OnlyRestMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tipID := chi.URLParam(r, "tipID")
	if tipID == "" {
		http.Error(w, "tipIDが必要です", http.StatusBadRequest)
		return
	}
	messages, err := h.uc.GetAllMessages(r.Context(), tipID)
	if err != nil {
		http.Error(w, "メッセージの取得に失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}
	response := ToChatMessagesResponse(messages)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
