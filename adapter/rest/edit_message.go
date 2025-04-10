package rest

import (
	"encoding/json"
	"net/http"

	"github.com/minminseo/tipstar-chat-api/usecase"
)

// メッセージ編集のリクエストを受取、編集系のユースケース（インターフェース）を実行するハンドラ。
type EditMessageHandler struct {
	usecase usecase.EditMessageInput // usecase.EditMessageInput（インターフェース）を型として持つusecaseフィールドを定義
}

// ユースケース層のインターフェースを満たす実装を受け取るコンストラクタ関数（ユースケースに依存せずここで使えるようにする）
func NewEditMessageHandler(uc usecase.EditMessageInput) *EditMessageHandler {
	return &EditMessageHandler{usecase: uc} // ユースケースの実装をuc.usecaseに代入して依存注入
}

// JWT検証、認証、リクエスト構造体をドメインモデル構造体にマッピング、ユースケースの実行、レスポンスの生成、ブロードキャストの実行を行うハンドラ
func (h *EditMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// クライアントからのJSONリクエストを構造体にデコード
	var req EditMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエスト形式が不正です", http.StatusBadRequest)
		return
	}

	// JWTから取り出して埋め込まれたcontextからuserIDをとりだす
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "ユーザーIDが不正です", http.StatusUnauthorized)
		return
	}

	// リクエスト構造体からドメインモデル構造体にマッピング
	domainMsg := ToEditDomainModel(&req, userID)

	// ユースケース実行
	if err := h.usecase.Execute(domainMsg.ID, domainMsg.UserID, domainMsg.Content); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 編集成功したらレスポンス構造体にマッピング（編集成功メッセージ含む）
	res := ToEditMessageResponse(domainMsg.ID)

	// JSON形式でレスポンスを返却。
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

	// メッセージ編集処理の成功をレスポンスしたもののその結果もブロードキャストする必要があるからここにブロードキャストする関数を呼び出す処理が必要になりそう
}
