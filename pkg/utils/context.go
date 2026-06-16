package utils

import (
	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"

// SetUserID menyimpan user_id ke dalam Gin Context
func SetUserID(c *gin.Context, userID string) {
	c.Set(ContextUserIDKey, userID)
}

// GetUserID mengambil user_id dari Gin Context
func GetUserID(c *gin.Context) (string, bool) {
	id, exists := c.Get(ContextUserIDKey)
	if !exists {
		return "", false
	}

	userID, ok := id.(string)
	// println("GetUserID berhasil:", userID)
	return userID, ok
}
