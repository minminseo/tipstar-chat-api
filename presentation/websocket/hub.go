package websocket

import (
	"log"
	"sync"
	"time"
)

// 全てのRoomを管理するHub構造体
// 各RoomはtipIDをキーとして持ち、Hub経由で動的に送信と受信を行う
type Hub struct {
	Rooms map[string]*Room // キーは各Roomに対応するtipID
	mu    sync.RWMutex
}

// Hubをインスタンス化する関数
func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

// tipIDに対応するRoomを取得し、そのRoomが存在しなければ新しくインスタンス化しHubの管理下（Roomsマップ）に登録
func (h *Hub) GetRoom(tipID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	room, ok := h.Rooms[tipID]
	if !ok {
		room = NewRoom(tipID)
		h.Rooms[tipID] = room
	}
	return room
}

// CheckIdleConnections関数を呼び出し、5分毎にアイドリング状態のRoomをチェックしConnectionが0なったRoomを削除
func (h *Hub) Run() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		h.mu.Lock()
		for tipID, room := range h.Rooms {
			room.CheckIdleConnections()
			if room.IsEmpty() {
				delete(h.Rooms, tipID)
			}
		}
		h.mu.Unlock()
		log.Println("Hub: 5分間アイドリングしたRoomを削除")
	}
}
