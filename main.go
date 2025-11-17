package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/a5415091-collab/go-gin-todo-app/handler" // ← ここを自分のモジュール名に変更
)

func main() {
	r := gin.Default()

	// 動作確認
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// 認証系
	r.POST("/signup", handler.Signup)
	r.POST("/login", handler.Login)

	// Todo系
	r.GET("/todos", handler.GetTodos)

	r.Run(":8080")
}
