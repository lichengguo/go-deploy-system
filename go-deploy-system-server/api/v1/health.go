package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Health 心跳检测
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
	})
}
