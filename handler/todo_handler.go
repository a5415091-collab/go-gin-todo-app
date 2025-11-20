package handler

import (
	"net/http"
	"strconv"

	"github.com/a5415091-collab/go-gin-todo-app/logger"
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
	userIDAny, exists := c.Get("userID")
	if !exists {
		logger.Logger.Warn(
			"userID not found in context",
			"handler", "GetTodos",
			"error", "missing userID",
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userID := userIDAny.(uint)
	logger.Logger.Info("request received", "handler", "GetTodos", "userID", userID)

	todos, err := h.todoService.FindAll(userID)
	if err != nil {
		logger.Logger.Error("failed to get todos", "userID", userID, "reason", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Info("get todos success", "userID", userID, "count", len(todos))
	c.JSON(http.StatusOK, todos)
}

// --- GET /todos/:id (詳細) ---
func (h *TodoHandler) GetTodo(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		logger.Logger.Warn(
			"userID not found in context",
			"handler", "GetTodo",
			"error", "missing userID",
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userID := userIDAny.(uint)
	id, _ := strconv.Atoi(c.Param("id"))
	logger.Logger.Info("request received", "handler", "GetTodo", "userID", userID, "todoID", uint(id))

	todo, err := h.todoService.FindByID(userID, uint(id))
	if err != nil {
		logger.Logger.Warn("todo not found", "todoID", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	logger.Logger.Info("get todo success", "todoID", id)
	c.JSON(http.StatusOK, todo)
}

// --- POST /todos (新規作成) ---
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		logger.Logger.Warn(
			"userID not found in context",
			"handler", "CreateTodo",
			"error", "missing userID",
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userID := userIDAny.(uint)
	logger.Logger.Info("request received", "handler", "CreateTodo", "userID", userID)

	var req struct {
		Title string `json:"title" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("create todo validation failed", "reason", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required and must be 1-100 characters"})
		return
	}

	todo, err := h.todoService.Create(userID, req.Title)
	if err != nil {
		logger.Logger.Error("failed to create todo", "userID", userID, "reason", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Info("create todo success", "todoID", todo.ID, "userID", userID)
	c.JSON(http.StatusOK, todo)
}

// --- PUT /todos/:id (更新) ---
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		logger.Logger.Warn(
			"userID not found in context",
			"handler", "UpdateTodo",
			"error", "missing userID",
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userID := userIDAny.(uint)
	id, _ := strconv.Atoi(c.Param("id"))
	logger.Logger.Info("request received", "handler", "UpdateTodo", "userID", userID, "todoID", id)

	var req struct {
		Title string `json:"title" binding:"required,min=1,max=100"`
		Done  *bool  `json:"done" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("update validation failed", "reason", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: title must be 1-100 chars and done must be true/false",
		})
		return
	}

	todo, err := h.todoService.Update(userID, uint(id), req.Title, req.Done)
	if err != nil {
		logger.Logger.Error("update failed", "todoID", id, "reason", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Info("update todo success", "todoID", id)
	c.JSON(http.StatusOK, todo)
}

// --- DELETE /todos/:id ---
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		logger.Logger.Warn(
			"userID not found in context",
			"handler", "DeleteTodo",
			"error", "missing userID",
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userID := userIDAny.(uint)
	id, _ := strconv.Atoi(c.Param("id"))
	logger.Logger.Info("request received", "handler", "DeleteTodo", "userID", userID, "todoID", id)

	err := h.todoService.Delete(userID, uint(id))
	if err != nil {
		logger.Logger.Error("delete failed", "todoID", id, "reason", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Info("delete todo success", "todoID", id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
