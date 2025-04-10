package main

import (
	"log"
	"os"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/minminseo/tipstar-chat-api/infra/db"
	"github.com/minminseo/tipstar-chat-api/router"
	"github.com/minminseo/tipstar-chat-api/usecase"
)

func main() {
	// 環境変数からDB接続情報を取得
	if err := godotenv.Load(); err != nil {
		log.Println(".envファイル読み込みエラー")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URLが設定されていません")
	}

	// PostgreSQL への接続プール作成
	pool, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer pool.Close()

	msgRepo := db.NewPgxMessageRepository(pool)

	editUC := usecase.NewEditMessageUsecase(msgRepo)
	deleteUC := usecase.NewDeleteMessageUsecase(msgRepo)

	r := router.NewRouter(editUC, deleteUC)

	log.Println("サーバー起動: http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("サーバー起動エラー: %v", err)
	}

}
