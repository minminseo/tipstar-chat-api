package domain

import (
	"errors"
	"time"
)

// チャットではなくメッセージと呼んでいく
// 型エイリアスを使って、可読性を上げる
type MessageID string
type TipID string
type UserID string

type Message struct {
	ID        MessageID  // メッセージ全部を識別する用途
	TipID     TipID      // 各メッセージがどのTipID（実質チャットルーム）に属するか識別する用
	UserID    UserID     // メッセージの送信主識別する用
	Content   string     // メッセージの文章
	CreatedAt time.Time  // メッセージの送信日時
	UpdatedAt time.Time  // CreatedAtと比較して未編集かは判定できるのと、nil持たせてもあんまり意味ないのでポインタ型にはしない
	DeletedAt *time.Time // 削除されてないという状態を分かりやすくしたい（nil使いたい）のでポインタ型
	IsAuthor  bool       // メッセージが投稿主のものかどうかUI制御するためのフラグ（永続化はしない）
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

// TODO:メッセージの所有権を検証する関数を共通化して切り出すかどうか決める
/*
// メッセージの所有権を持っているかどうかの判定
// TODO:大差ないけどあとで値レシーバーとポインタレシーバーのパフォの違いを検証してみる（ベンチマーク、メモリプロファイリング、pprof等）
// 理由：比較しかしない＆呼び出し頻度も低いと思うけど、コピーコスト高め（120バイトくらい）なのと一貫性を考慮して一応ポインタレシーバーにする
// 参考：ただ、特定の型に対して、実装されているメソッドのうち一つでもポインタ使ってるなら、比較のみだとしても一貫して全てのメソッドをポインタレシーバーにするのが推奨らしいからこれでいいかも
func (m *Message) ValidateOwnership(userID UserID) bool {

	// 編集対象にしているメッセージが自分のメッセージか判定
	return m.UserID == userID
}
*/

// メッセージの編集処理（ヒープメモリ上のMessageの実体に対する書き換え）
func (m *Message) SetEditedContent(userID UserID, newContent string) error {

	// TODO:共通化候補
	// 所有権の検証
	if m.UserID != userID {
		return errors.New("このメッセージを編集する権限がありません")
	}

	// 論理削除済み、つまり削除日が存在するメッセージの編集をできないようにする
	if m.DeletedAt != nil {
		return errors.New("このメッセージはすでに削除されています")
	}

	// 編集するメッセージが空の文字列の場合はエラーを返す
	if newContent == "" {
		return errors.New("メッセージ内容が空です")
	}

	// if文全て通過したら、mのポインタが指すメモリ上のMessageインスタンスのContent、UpdatedAtフィールドをそれぞれ新しい値で書き換える。
	m.Content = newContent
	m.UpdatedAt = time.Now()
	return nil
}

// メッセージの削除処理（ヒープメモリ上のMessageの実体に対する書き換え）
// DB的には論理削除
func (m *Message) SetDeletedMessage(userID UserID) error {

	// TODO:共通化候補
	// 所有権の検証
	if m.UserID != userID {
		return errors.New("このメッセージを削除する権限がありません")
	}
	if m.DeletedAt != nil {
		return errors.New("このメッセージはすでに削除済みです")
	}

	// if文全て通過したら、mのポインタが指すメモリ上のMessageインスタンスのDeletedAtフィールドを新しい値（削除日時）で書き換える。
	now := time.Now()
	m.DeletedAt = &now
	return nil
}
