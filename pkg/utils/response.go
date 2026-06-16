package utils

import (
	"github.com/gin-gonic/gin"
)

// RespondJSON mengirim response JSON sukses
func RespondJSON(c *gin.Context, status int, payload interface{}) {
	c.JSON(status, gin.H{
		"status":  "success",
		"data":    payload,
		"message": nil,
	})
}

// RespondError mengirim response JSON error
func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":  "error",
		"data":    nil,
		"message": message,
	})
}

// RespondMessage mengirim response dengan message custom
func RespondMessage(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":  "success",
		"data":    nil,
		"message": message,
	})
}
