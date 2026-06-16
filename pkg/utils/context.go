package utils

import (
	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"
const ContextTenantIDKey = "tenantID"

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

// SetTenantID menyimpan tenant_id ke dalam Gin Context
func SetTenantID(c *gin.Context, tenantID int) {
	c.Set(ContextTenantIDKey, tenantID)
}

// GetTenantID mengambil tenant_id dari Gin Context
func GetTenantID(c *gin.Context) (int, bool) {
	id, exists := c.Get(ContextTenantIDKey)
	if !exists {
		return 0, false
	}

	tenantID, ok := id.(int)
	return tenantID, ok
}

