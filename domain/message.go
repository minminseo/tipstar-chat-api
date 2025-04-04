package domain

import (
	"errors"
	"time"
)

// チャットではなくメッセージと呼んでいく

type MessageID string
type TipID string
type UserID string

type Message struct {
	ID        MessageID  // メッセージ全部を識別する用途
	TipID     TipID      // 各メッセージ（子）がどのTipsID（実質チャットルーム）に属するか識別する用
	UserID    UserID     // メッセージの送信主識別する用 + チャットの枠の上あたりに名前表示する用
	Content   string     // メッセージの文章
	CreatedAt time.Time  // メッセージの送信日時
	UpdatedAt time.Time  // CreatedAtと比較して未編集かは判定できるのと、nil持たせてもあんまり意味ないのでポインタ型にはしない
	DeletedAt *time.Time // 削除されてないという状態を分かりやすくしたい（nil使いたい）のでポインタ型
	IsAuthor  bool       // メッセージが投稿主のものかどうかUI制御するためのフラグ（DB保存はしない）
}

// メッセージのファクトリ関数定義
func NewMessage(id MessageID, tipID TipID, userID UserID, content string, isAuthor bool) (*Message, error) {
	if content == "" {
		return nil, errors.New("メッセージが空")
	}

	// メッセージ作成日と更新日はこのアプリのドメインモデルの一部（）にするので、この２つの値の初期化もファクトリ関数内で初期化する。
	// IDはこのアプリでは意味を持たず単なる識別用でしか使わないので、ファクトリ関数内で初期化しない。
	now := time.Now()

	return &Message{
		ID:        id,
		TipID:     tipID,
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil, // 削除済み等のUI表示をするというドメインモデルの一部になるのでファクトリ関数内でnilで初期化する
		IsAuthor:  isAuthor,
	}, nil
}

// メッセージの編集権限を持っているかどうかの判定
func (m *Message) CanEdit(by UserID) bool {

	// 編集対象にしているメッセージが自分のメッセージか、削除済みじゃないかどうか判定
	return m.UserID == by && m.DeletedAt == nil
}

// メッセージの編集処理
func (m *Message) EditContent(newContent string) error {

	// 削除日が存在する（論理削除）メッセージの編集をできないようにする
	if m.DeletedAt != nil {
		return errors.New("このメッセージはすでに削除されています")
	}

	// 編集するメッセージが空の文字列の場合はエラーを返す
	if newContent == "" {
		return errors.New("メッセージ内容が空です")
	}

	// if文全て通過したら、メッセージのcontentと更新日時を更新
	m.Content = newContent
	m.UpdatedAt = time.Now()
	return nil
}

// メッセージの削除処理
func (m *Message) Delete(by UserID) error {
	if m.UserID != by {
		return errors.New("このメッセージを削除する権限がありません")
	}
	if m.DeletedAt != nil {
		return errors.New("このメッセージはすでに削除済みです")
	}

	// if全て通過したら、論理削除
	now := time.Now()
	m.DeletedAt = &now
	return nil
}
