package domain

// メッセージまわりのDB操作のインターフェース定義
// ユースケース層が「どんな操作ができるか」だけを知るためのインターフェースを定義
// インフラ層でこのインターフェースの関数たちを実装して、ユースケース層はこのインターフェースに依存する
type MessageRepository interface {
	FindByID(id MessageID) (*Message, error) // クライアントからきたMessageIDを元にDBからメッセージを取得する関数
	Insert(msg *Message) error               // メッセージをDBに挿入する関数（送信系は未実装）
	Update(msg *Message) error               // メッセージを編集する関数
	Delete(msg *Message) error               // メッセージを論理削除する関数
}
