package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/minminseo/tipstar-chat-api/presentation/websocket"
)

func NewRouter(
	restHandler http.Handler, // 一覧取得ハンドラー
	wsHandler *websocket.OnlyWSMessageHandler, // Websocket系の処理のハンドラー
	hub *websocket.Hub, // WebSocketのハブ（ルーム管理用）
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// 認証は各リクエストのヘッダーからuser_idを受け取る前提（今後JWT認証に変更する）

	r.Get("/messages/{tipID}", restHandler.ServeHTTP)

	r.Get("/ws/{tipID}", func(w http.ResponseWriter, r *http.Request) {
		tipID := chi.URLParam(r, "tipID")

		// Websocketへの昇格処理。Websocketは双方向通信のためのプロトコル。
		// connは接続情報、userIDはユーザーID
		conn, userID, err := websocket.UpgradeHTTP(w, r)
		if err != nil {
			return
		}

		// Connection構造体をインスタンス化
		wsConn := &websocket.Connection{
			Conn:       conn,
			UserID:     userID,
			Send:       make(chan []byte, 256),
			LastActive: time.Now(),
			Ctx:        r.Context(),
		}

		//取得したtipIDに紐づくRoom（実質のチャットルーム）を取得
		room := hub.GetRoom(tipID)

		// 取得したRoomに対して、roomのポインタ型をレシーバーとして持つJoinメソッドにインスタンス化したConnectionオブジェクト（wsConn）を引数として渡す
		// Joinでは該当roomのclientsフィールドにwsConnが追加される
		room.Join(wsConn)

		// ゴルーチンで非同期でclientsに存在するクライアントにメッセージを送信する（書き込みは復数ユーザーへのブロードキャストという形になるためゴルーチンを使う（排他制御必須））
		go wsConn.WritePump()

		// クライアントからのメッセージを受信し、ハンドラーに渡す
		wsConn.ReadPump(wsHandler.HandleWSMessage)

		// 接続クライアントが切断されたら実行
		room.Leave(wsConn)
	})

	return r
}
