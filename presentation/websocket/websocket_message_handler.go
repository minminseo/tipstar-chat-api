package websocket

// ここではWebSocket経由で受信したリクエストのハンドリングを行う（最終的にJSON形式にしてブロードキャスト関数を呼び出す）
// ブロードキャストやルームの取得はHubを参照して行う
// コンストラクタの引数hubはまずnilを受け取り、その後SetHubメソッドで外部でインスタンス化されたHubを注入させる。

import (
	"encoding/json"
	"log"

	"github.com/minminseo/tipstar-chat-api/usecase"
)

type OnlyWSMessageHandler struct {
	uc  usecase.OnlyWSUsecase // usecase.OnlyWSUsecase（インターフェース）を型として持つucフィールドを定義
	hub *Hub                  // ルーム管理用のHub
}

// ユースケースのインターフェースを満たすメソッドをプレゼンテーション層に注入するコンストラクタ関数（ユースケース内部の処理を隠してここで使えるようにする）
// この時点ではHubにnilを渡す。まだインスタンス化されていないから。
func NewOnlyWSMessageHandler(uc usecase.OnlyWSUsecase, hub *Hub) *OnlyWSMessageHandler {
	return &OnlyWSMessageHandler{
		uc:  uc,
		hub: hub,
	}
}

// 外部で生成されたHubインスタンスを、OnlyWSMessageHandler内部のhubフィールドに注入する
// Hubのインスタンス化はmain.go。この時点でnil状態のhubフィールドにインスタンス化されたHubを注入する
// その後、ハンドラーは受信したWebsocketメッセージを処理し、必要に応じてHubを介してそのメッセージに対応するチャットルーム（Room（=tip毎のチャットルーム））を取得し、ブロードキャスト
func (h *OnlyWSMessageHandler) SetHub(hub *Hub) {
	h.hub = hub
}

// Websocket経由のリクエストのボディに含まれるTypeフィールドの値毎に処理を分岐。
// 1個の接続クライアントには基本1個の読み取りループを回すので、読み取りループ実行の関数（このアプリではReadPump）ではこの関数を呼び出してTypeフィールドの値毎に処理を分岐する
func (h *OnlyWSMessageHandler) HandleWSMessage(rawMsg []byte, conn *Connection) {
	var req WSRequestMessage
	if err := json.Unmarshal(rawMsg, &req); err != nil {
		log.Printf("HandleWSMessage: WSリクエストのJSONのデコードに失敗: %v", err)
		return
	}
	switch req.Type {
	case "send":
		h.SendMessageHandler(rawMsg, conn)
	case "edit":
		h.EditMessageHandler(rawMsg, conn)
	case "delete":
		h.DeleteMessageHandler(rawMsg, conn)
	default:
		log.Printf("HandleWSMessage: 予期しないリクエストのTypeが含まれています: %s", req.Type)
	}
}

// メッセージ送信のハンドラー
func (h *OnlyWSMessageHandler) SendMessageHandler(rawMsg []byte, conn *Connection) {
	var req WSRequestMessage
	if err := json.Unmarshal(rawMsg, &req); err != nil {
		log.Printf("SendMessageHandler: WSリクエストのJSONのデコードに失敗: %v", err)
		return
	}
	if req.Type != "send" {
		log.Printf("SendMessageHandler: 予期しないリクエストのTypeが含まれています: %s", req.Type)
		return
	}

	// ユーザーIDは接続時に取得した conn.UserID を使用する
	msg, err := ToSendDomainFromWSRequest(&req, conn.UserID)
	if err != nil {
		log.Printf("SendMessageHandler: ドメインモデルへの変換に失敗: %v", err)
		return
	}
	if err := h.uc.ExecuteSendMessage(conn.Context(), msg); err != nil {
		log.Printf("SendMessageHandler: メッセージの永続化に失敗: %v", err)
		return
	}
	wsResp := ToBroadcastMessage(msg)
	bMsg, err := json.Marshal(wsResp)
	if err != nil {
		log.Printf("SendMessageHandler: ブロードキャスト用メッセージのJSONエンコードに失敗: %v", err)
		return
	}
	room := h.hub.GetRoom(req.TipID)
	room.Broadcast(bMsg)
}

// メッセージ編集のハンドラー
func (h *OnlyWSMessageHandler) EditMessageHandler(rawMsg []byte, conn *Connection) {
	var req WSRequestMessage
	if err := json.Unmarshal(rawMsg, &req); err != nil {
		log.Printf("EditMessageHandler: WSリクエストのJSONのデコードに失敗: %v", err)
		return
	}
	if req.Type != "edit" {
		log.Printf("EditMessageHandler: 予期しないリクエストのTypeが含まれています: %s", req.Type)
		return
	}

	// ユーザーIDは接続時に取得した conn.UserID を使用する
	msg, err := ToEditDomainFromWSRequest(&req, conn.UserID)
	if err != nil {
		log.Printf("EditMessageHandler: ドメインモデルへの変換に失敗: %v", err)
		return
	}
	if err := h.uc.EditMessage(conn.Context(), msg.ID, msg.UserID, msg.Content); err != nil {
		log.Printf("EditMessageHandler: メッセージの編集に失敗: %v", err)
		return
	}
	wsResp := ToEditBroadcastMessage(msg)
	bMsg, err := json.Marshal(wsResp)
	if err != nil {
		log.Printf("EditMessageHandler: ブロードキャスト用メッセージのJSONエンコードに失敗: %v", err)
		return
	}
	room := h.hub.GetRoom(req.TipID)
	room.Broadcast(bMsg)
}

// メッセージ削除のハンドラー
func (h *OnlyWSMessageHandler) DeleteMessageHandler(rawMsg []byte, conn *Connection) {
	var req WSRequestMessage
	if err := json.Unmarshal(rawMsg, &req); err != nil {
		log.Printf("DeleteMessageHandler: WSリクエストのJSONのデコードに失敗: %v", err)
		return
	}
	if req.Type != "delete" {
		log.Printf("DeleteMessageHandler: 予期しないリクエストのTypeが含まれています: %s", req.Type)
		return
	}

	// ユーザーIDは接続時に取得した conn.UserID を使用する
	msg, err := ToDeleteDomainFromWSRequest(&req, conn.UserID)
	if err != nil {
		log.Printf("DeleteMessageHandler: ドメインモデルへの変換に失敗: %v", err)
		return
	}
	if err := h.uc.DeleteMessage(conn.Context(), msg.ID, msg.UserID); err != nil {
		log.Printf("DeleteMessageHandler: メッセージの削除に失敗: %v", err)
		return
	}
	wsResp := ToDeleteBroadcastMessage(msg)
	bMsg, err := json.Marshal(wsResp)
	if err != nil {
		log.Printf("DeleteMessageHandler: ブロードキャスト用メッセージのJSONエンコードに失敗: %v", err)
		return
	}
	room := h.hub.GetRoom(req.TipID)
	room.Broadcast(bMsg)
}
