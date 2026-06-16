package utils

import (
	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":  "error",
		"message": message,
	})
}

func JSONSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}
