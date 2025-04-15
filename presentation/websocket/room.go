package websocket

import (
	"sync"
	"time"
)

// 各Tipに対応するチャットルームを管理する構造体
type Room struct {
	TipID        string
	Clients      map[*Connection]bool
	mu           sync.RWMutex
	LastActivity time.Time     // 最後のアクティビティ時刻
	idleDuration time.Duration // 5分間何もアクティビティがないかどうか判定するためのフィールド
}

// tipIDに対応するRoomインスタンスを生成
func NewRoom(tipId string) *Room {
	return &Room{
		TipID:        tipId,
		Clients:      make(map[*Connection]bool),
		LastActivity: time.Now(),
		idleDuration: 5 * time.Minute,
	}
}

// 引数で渡されたConnectionをRoomに追加し、参加時刻をLastActivityに記録
func (r *Room) Join(c *Connection) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clients[c] = true // ConnectionをRoomのclientsマップに追加
	r.LastActivity = time.Now()
}

// 引数で渡されたConnectionをRoomのClientsマップから削除しClose。この時の最後のアクティビティ時刻をLastActivityに記録
func (r *Room) Leave(c *Connection) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.Clients[c]; ok {
		delete(r.Clients, c)
		c.Conn.Close()
	}
	r.LastActivity = time.Now()
}

// Roomに属する全クライアント（Connection）のSendチャネルにメッセージを送信する（代入する）。
func (r *Room) Broadcast(message []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.LastActivity = time.Now() // 最後のアクティビティ時刻を更新
	for client := range r.Clients {
		select {
		case client.Send <- message:
		default:
			// チャネルがブロックしている場合はスキップ
		}
	}
}

// Roomが5分間何もアクティビティが無い場合（LastActivityからの経過時間がidleDurationを超えている場合）、全ConnectionをCloseし、Clientsマップから削除
func (r *Room) CheckIdleConnections() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if time.Since(r.LastActivity) > r.idleDuration {
		for client := range r.Clients {
			client.Conn.Close()
			delete(r.Clients, client)
		}
	}
}

// ルーム内にConnectionが無い場合はtrueを返す（ルームないにクライアントがいないかどうかを確認する）
func (r *Room) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Clients) == 0
}
