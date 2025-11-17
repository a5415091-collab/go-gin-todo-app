package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"todos": []string{},
	})
}
