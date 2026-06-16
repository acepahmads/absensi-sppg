package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetClaims(c *gin.Context) (*Claims, bool) {
	claimsRaw, exists := c.Get("claims")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
		return nil, false
	}

	claims, ok := claimsRaw.(*Claims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return nil, false
	}

	return claims, true
}
