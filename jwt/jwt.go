package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("super_secret_key_123") // ← 好きなキーにしてOK（外部に出さない）

// -----------------------------
// JWTを作る関数（login時に使う）
// -----------------------------
func CreateToken(userID uint) (string, error) {
	// トークンに入れる情報（Claims）
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 有効期限 24h
	}

	// 署名アルゴリズム HS256 を使ってトークンを作る
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// secretKey で署名して実体の文字列にする
	return token.SignedString(secretKey)
}

// -----------------------------
// JWTを検証する関数（Middlewareで使う）
// -----------------------------
func VerifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}
