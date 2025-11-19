package handler

import (
	"net/http"
	"strconv"

	"github.com/a5415091-collab/go-gin-todo-app/service"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoService service.TodoService
}

func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{todoService}
}

// --- GET /todos (一覧) ---
func (h *TodoHandler) GetTodos(c *gin.Context) {
	// Middleware でセットした userID を取り出す
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userID := userIDAny.(uint)

	// Service 経由で取得
	todos, err := h.todoService.FindAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// --- GET /todos/:id (詳細) ---
func (h *TodoHandler) GetTodo(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}
	userID := userIDAny.(uint)

	id, _ := strconv.Atoi(c.Param("id"))

	todo, err := h.todoService.FindByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// --- POST /todos (新規作成) ---
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}
	userID := userIDAny.(uint)

	var req struct {
		Title string `json:"title" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required and must be 1-100 characters"})
		return
	}

	todo, err := h.todoService.Create(userID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// --- PUT /todos/:id (更新) ---
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}
	userID := userIDAny.(uint)

	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Title string `json:"title" binding:"required,min=1,max=100"`
		Done  *bool  `json:"done" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: title must be 1-100 chars and done must be true/false",
		})
		return
	}

	todo, err := h.todoService.Update(userID, uint(id), req.Title, req.Done)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// --- DELETE /todos/:id ---
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}
	userID := userIDAny.(uint)

	id, _ := strconv.Atoi(c.Param("id"))

	err := h.todoService.Delete(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
