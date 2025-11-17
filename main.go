package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/a5415091-collab/go-gin-todo-app/db"
	"github.com/a5415091-collab/go-gin-todo-app/handler"
	"github.com/a5415091-collab/go-gin-todo-app/model"
)

func main() {

	r := gin.Default()

	db.Init()
	db.DB.AutoMigrate(&model.User{})
	db.DB.AutoMigrate(&model.Todo{})

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
