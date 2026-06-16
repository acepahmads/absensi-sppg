package handler

import (
	"bytes"
	"absensi-sppg/internal/model"
	"absensi-sppg/internal/service"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	// Baca body mentah
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca body"})
		return
	}

	// fmt.Println("Raw JSON:", string(bodyBytes))

	// reset body supaya masih bisa dibaca oleh ShouldBindJSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	var req model.RegisterAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Call the service to register the user
	email, err := h.userService.RegisterUser(req)
	if err != nil {
		fmt.Println("Error registering user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully with code " + email})
}

func (h *UserHandler) GetUserInfoByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	ID, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUserInfoByID(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	var jabatan string
	var name string
	if user.Jabatan.Valid {
		jabatan = user.Jabatan.String
	}

	if user.Name.Valid {
		name = user.Name.String
	}

	resp := model.UserInfoAccount{
		Name:     name,
		Position: jabatan,     // *string → aman
		Photo:    user.Photos, // []string
		Role:     user.Role,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Registered(c *gin.Context) {
	names, err := h.userService.Registered()
	if err != nil {
		fmt.Println("Error get names user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get name of user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"names": names})
}

func (h *UserHandler) GetLeaders(c *gin.Context) {
	leaders, err := h.userService.GetLeaders()
	if err != nil {
		fmt.Println("Error get names user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get name of user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"names": leaders})
}
