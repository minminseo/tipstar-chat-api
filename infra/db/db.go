package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// pgxpoolライブラリで実装されているPoolという構造体？をポインタ型にしてインスタンス化
func NewDB(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// コネクションプール作成
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("DB接続プールの初期化に失敗しました: %w", err)
	}

	// 接続確認
	if err := pool.Ping(ctx); err != nil {
		pool.Close() // 失敗したらClose
		return nil, fmt.Errorf("DBへの疎通確認に失敗しました: %w", err)
	}

	return pool, nil
}
