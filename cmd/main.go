package main

// サーバー全体の初期化と起動処理を担当します。

/*
処理の流れ
1. 環境変数の取得
2. データベース接続プールの初期化（インフラ層の実装を利用）
3. リポジトリ層、ユースケース層、ハンドラー層のインスタンス化と依存注入
4. WebSocketのルーム管理のためのHubのインスタンス化と、ルーム管理ループの起動
5. ルーターの初期化（依存性注入済みのハンドラーを渡す）
6. サーバーの起動（指定されたポートでHTTPサーバーを起動）

*/

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/minminseo/tipstar-chat-api/infra/db"
	"github.com/minminseo/tipstar-chat-api/presentation/rest"
	"github.com/minminseo/tipstar-chat-api/presentation/websocket"
	"github.com/minminseo/tipstar-chat-api/router"
	"github.com/minminseo/tipstar-chat-api/usecase"
)

func main() {
	// .envファイル読み込み（開発環境用）
	if err := godotenv.Load(); err != nil {
		log.Println(".envファイル読み込みエラー")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URLが設定されていません")
	}

	// データベース接続プールの作成
	pool, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer pool.Close()

	// インスタンス化と注入
	// コンストラクタを起動、外側でインスタンス化したDB接続プール注入、永続化処理のインターフェースのメソッドの具象実装をインスタンス化
	msgRepo := db.NewPgxMessageRepository(pool)

	// コンストラクタを起動、外側でインスタンス化した永続化処理を注入、ユースケースのインターフェースのメソッドの具象実装をインスタンス化
	onlyRestUC := usecase.NewOnlyRestMessageUseCase(msgRepo)
	onlyWSCUC := usecase.NewOnlyWSMessageUseCase(msgRepo)

	// コンストラクタを起動、外側でインスタンス化したユースケースを注入、ハンドラーのインターフェースのメソッドの具象実装をインスタンス化
	restHandler := rest.NewOnlyRestMessageHandler(onlyRestUC)
	wsHandler := websocket.NewOnlyWSMessageHandler(onlyWSCUC, nil) // hubは後でセットするのでnilを渡す

	// WebSocketのハブ生成とルーム管理ループの起動
	hub := websocket.NewHub()
	wsHandler.SetHub(hub) // wsHandler 内で Hub を利用する場合の setter を実装しておく
	go hub.Run()

	// 依存注入済みのハンドラーを渡す
	r := router.NewRouter(restHandler, wsHandler, hub)

	// サーバー起動
	port := os.Getenv("PORT")

	addr := ":" + port
	log.Printf("サーバー起動: http://localhost:%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}
