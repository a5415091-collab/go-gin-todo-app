package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/a5415091-collab/go-gin-todo-app/db"
	"github.com/a5415091-collab/go-gin-todo-app/handler"
	"github.com/a5415091-collab/go-gin-todo-app/repository"
	"github.com/a5415091-collab/go-gin-todo-app/service"
)

func main() {
	r := gin.Default()

	// DB 初期化
	db.Init()

	// Repository 作成
	userRepo := repository.NewUserRepository()
	todoRepo := repository.NewTodoRepository()

	// Service 作成
	authService := service.NewAuthService(userRepo)
	todoService := service.NewTodoService(todoRepo)

	// Handler に service を渡す
	authHandler := handler.NewAuthHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)

	// 動作確認用
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 認証系
	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)

	// TODO系
	r.GET("/todos", todoHandler.GetTodos)
	r.GET("/todos/:id", todoHandler.GetTodo)
	r.POST("/todos", todoHandler.CreateTodo)
	r.PUT("/todos/:id", todoHandler.UpdateTodo)
	r.DELETE("/todos/:id", todoHandler.DeleteTodo)

	r.Run(":8080")

}
