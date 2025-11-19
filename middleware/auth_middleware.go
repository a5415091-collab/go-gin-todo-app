package middleware

import (
	"net/http"
	"strings"

	myjwt "github.com/a5415091-collab/go-gin-todo-app/jwt"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Authorization ヘッダ
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		// "Bearer xxx" を分割
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// トークン検証
		token, err := myjwt.VerifyToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Claims を型アサーション
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			c.Abort()
			return
		}

		// user_id を Float64 → uint
		userIDFloat := claims["user_id"].(float64)
		userID := uint(userIDFloat)

		// context に保存
		c.Set("userID", userID)

		// 次へ
		c.Next()
	}
}
