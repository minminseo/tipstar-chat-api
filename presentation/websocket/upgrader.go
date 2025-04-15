package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// HTTP接続をWebSocket接続（双方向通信）に昇格させるためのUpgrater（構造体）を定義。
// 昇格処理は下の方で実装（JWTを検証したあとに実行）。
// wsUpgraderはwebsocketパッケージのUpgrader型の構造体リテラルによって直接インスタンス化
var wsUpgrader = websocket.Upgrader{

	// CheckOriginフィールドにオリジンチェックの結果を代入しUpgraderをインスタンス化する
	// Websocket接続のオリジンチェック（クロスサイトWebSocketハイジャック対策というやつらしい）
	CheckOrigin: func(r *http.Request) bool {

		/*
			// 環境変数からURLを取得
			// Vercelデプロイ時に設定する必要ある。
			// 開発してるうちはとりあえずhttp://localhost:3000とかにする。使用済みなら適宜変更
			feURL := os.Getenv("FE_URL")

			// HTTPリクエストのヘッダーからOriginを取得し、環境変数で指定されたURLと一致するか確認
			return r.Header.Get("Origin") == feURL
		*/
		return true // 開発中はtrue（全オリジン許可）
	},
}

func UpgradeHTTP(w http.ResponseWriter, r *http.Request) (*websocket.Conn, string, error) { // 返り値として接続情報、ユーザーID、エラーを返す

	// JWT認証なし
	// ユーザーIDは "X-User-Id" ヘッダーから取得
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		http.Error(w, "X-User-Id ヘッダーがありません", http.StatusUnauthorized)
		return nil, "", http.ErrNoCookie
	}

	// HTTP接続をWebSocket接続へ昇格
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, "", err
	}
	return conn, userID, err
}
