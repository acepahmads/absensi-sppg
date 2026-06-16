package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"absensi-sppg/pkg/utils"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("TOken tidak ada")
			c.Redirect(http.StatusFound, "/")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.Redirect(http.StatusFound, "/")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}

		// Simpan user ID dan tenant ID ke context
		utils.SetUserID(c, claims.UserID)
		utils.SetTenantID(c, claims.TenantID)
		c.Set("claims", claims)

		// Inject into standard Go request context so repositories/services can extract it
		ctx := context.WithValue(c.Request.Context(), "tenantID", claims.TenantID)
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func JWTAuthDashboard1() gin.HandlerFunc {
	return func(c *gin.Context) {
		// authHeader := c.GetHeader("Authorization")
		tokenString, err := c.Cookie("token")
		// fmt.Println("Auth Header:", authHeader)
		// if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		if tokenString == "" {
			c.Redirect(http.StatusFound, "/")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			return
		}

		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Simpan user ID ke context
		// c.Set("userID", claims.UserID)
		utils.SetUserID(c, claims.UserID)
		c.Next()
	}
}

func JWTAuthDashboard() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("authToken")
		fmt.Println("Cek token", tokenStr)
		// fmt.Println("Token String:", tokenStr)
		if err != nil {
			fmt.Println("err0", err)
			c.SetCookie(
				"authToken",
				"",
				-1, // expire now
				"/",
				"",
				false,
				true,
			)
			c.Redirect(http.StatusFound, "/")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			return
		}

		claims, err := utils.ValidateJWT(tokenStr)
		// fmt.Println("Claims:", claims, "Error:", err)
		if err != nil {
			fmt.Println("err1", err)
			c.SetCookie(
				"authToken",
				"",
				-1, // expire now
				"/",
				"",
				false,
				true,
			)
			c.Redirect(http.StatusFound, "/")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}

		c.Set("claims", claims)
		utils.SetUserID(c, claims.UserID)
		utils.SetTenantID(c, claims.TenantID)

		// Inject into standard Go request context so repositories/services can extract it
		ctx := context.WithValue(c.Request.Context(), "tenantID", claims.TenantID)
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
