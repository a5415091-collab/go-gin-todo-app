package handler

import (
	"net/http"

	"github.com/a5415091-collab/go-gin-todo-app/jwt"
	"github.com/a5415091-collab/go-gin-todo-app/logger"
	"github.com/a5415091-collab/go-gin-todo-app/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// POST /signup
func (h *AuthHandler) Signup(c *gin.Context) {
	logger.Logger.Info("request received", "handler", "Signup")

	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=64"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("signup validation failed", "reason", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	err := h.authService.Signup(req.Email, req.Password)
	if err != nil {
		logger.Logger.Error("signup failed", "email", req.Email, "reason", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Info("user signup", "email", req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "signup success"})
}

// POST /login
func (h *AuthHandler) Login(c *gin.Context) {
	logger.Logger.Info("request received", "handler", "Login")

	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=64"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("login validation failed", "reason", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// 1. service.Login で認証
	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		logger.Logger.Warn("login failed", "email", req.Email, "reason", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// 2. JWT の発行
	token, err := jwt.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	// 3. トークンを返す
	logger.Logger.Info("user login success", "email", req.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"token":   token,
	})
}
