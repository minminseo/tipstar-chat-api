package rest

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	userIDKey = contextKey("userID")
)

// AuthorizationヘッダーからJWTを取り出し環境変数の秘密鍵で検証するミドルウェア
// 検証成功したらユーザーIDをContextに埋め込む
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorizationヘッダーを取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// "Bearer<token>形式かチェック
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			http.Error(w, "JWT secret not configured", http.StatusInternalServerError)
			return
		}

		// Base64デコード（SupabaseのシークレットキーはBase64エンコードされているので）
		decodedSecret, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode JWT secret: %v", err), http.StatusInternalServerError)
			return
		}
		// トークン検証
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return decodedSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// クレームからユーザーIDを取り出す
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			http.Error(w, "User ID not found in token", http.StatusUnauthorized)
			return
		}

		// ユーザーIDをContextに埋め込んでで次のハンドラに渡す
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// contextからユーザーIDを取り出す関数
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user id not found in context")
	}
	return userID, nil
}
