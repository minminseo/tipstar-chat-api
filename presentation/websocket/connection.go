package websocket

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 各接続クライアントのWebsocket接続を管理する構造体
type Connection struct {
	Conn       *websocket.Conn // 実際のWebSocket接続オブジェクト
	UserID     string          // 接続クライアントを識別するためのユーザーID
	Send       chan []byte     // 接続先へのブロードキャスト用チャネル
	LastActive time.Time       // 最後にデータの送受信があった時刻
	Ctx        context.Context // HTTPリクエストのContextを継承するフィールド
	mu         sync.Mutex      // ブロードキャスト時の排他制御用
}

// Connectionに紐づくContextを取得するメソッド
func (c *Connection) Context() context.Context {
	return c.Ctx
}

// 接続クライアントからのメッセージを読取るための関数（ループを使ってこれを実現する）
// 外部から渡されたhandler（コールバック）を呼び出す
func (c *Connection) ReadPump(handler func(msg []byte, c *Connection)) {
	// 無限ループさせてクライアントからのメッセージを受信し続ける
	// クライアント側が切断した場合、err != nilはtrueになり、ループを抜ける
	// その後router.goでLeaveメソッドが実行される
	for {
		_, message, err := c.Conn.ReadMessage()

		if err != nil {
			log.Printf("ReadPump: メッセージ読取りエラー: %v", err)
			break
		}
		c.LastActive = time.Now() // アクティビティ更新
		handler(message, c)       // 受信したメッセージとConnectionを引数として渡して、ReadPumpに渡されているhandler（コールバック関数）を呼び出す
	}
}

// 送信ループでSendチャネルに流し込まれるメッセージを取り出す→ブロードキャスト
// Sendチャネルにブロードキャスト用メッセージが送信されたら（代入されたら）、接続先（c.Conn）に書き込みクライアントへブロードキャスト
func (c *Connection) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for msg := range c.Send {

		// 書き込み時は排他制御する
		// 排他制御しないと、同一の共有リソースに対して同時に書き込みをしてしまいデータ競合が起こる
		c.mu.Lock()
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		c.mu.Unlock()
		if err != nil {
			break
		}
	}
}
