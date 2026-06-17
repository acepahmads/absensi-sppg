package handler

import (
	"bytes"
	"absensi-sppg/internal/model"
	"absensi-sppg/internal/service"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	tenantIDStr := c.Query("tenant_id")
	var tenantID int
	if tenantIDStr != "" {
		fmt.Sscanf(tenantIDStr, "%d", &tenantID)
	}
	if tenantID <= 0 {
		tenantID = 1
	}
	names, err := h.userService.Registered(tenantID)
	if err != nil {
		fmt.Println("Error get names user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get name of user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"names": names})
}

func (h *UserHandler) GetLeaders(c *gin.Context) {
	tenantIDStr := c.Query("tenant_id")
	var tenantID int
	if tenantIDStr != "" {
		fmt.Sscanf(tenantIDStr, "%d", &tenantID)
	}
	if tenantID <= 0 {
		tenantID = 1
	}
	leaders, err := h.userService.GetLeaders(tenantID)
	if err != nil {
		fmt.Println("Error get names user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get name of user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"names": leaders})
}

func (h *UserHandler) GetTenants(c *gin.Context) {
	tenants, err := h.userService.GetTenants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tenants list"})
		return
	}
	c.JSON(http.StatusOK, tenants)
}

func (h *UserHandler) GetAllUserKaryawan(c *gin.Context) {
	list, err := h.userService.GetAllUserKaryawan(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user karyawan list"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *UserHandler) CreateUserKaryawan(c *gin.Context) {
	var req model.UserKaryawan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.userService.CreateUserKaryawan(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user karyawan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User karyawan created successfully"})
}

func (h *UserHandler) UpdateUserKaryawan(c *gin.Context) {
	idParam := c.Param("id")
	var id int
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	var req model.UserKaryawan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	req.ID = id

	err = h.userService.UpdateUserKaryawan(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user karyawan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User karyawan updated successfully"})
}

func (h *UserHandler) DeleteUserKaryawan(c *gin.Context) {
	idParam := c.Param("id")
	var id int
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	err = h.userService.DeleteUserKaryawan(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user karyawan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User karyawan deleted successfully"})
}

func (h *UserHandler) GetLeadersList(c *gin.Context) {
	leaders, err := h.userService.GetLeadersList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get leaders list"})
		return
	}
	c.JSON(http.StatusOK, leaders)
}

// Leader CRUD
func (h *UserHandler) GetAllLeaders(c *gin.Context) {
	list, err := h.userService.GetAllLeaders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get leaders list"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *UserHandler) CreateLeader(c *gin.Context) {
	var req model.KaryawanLeader
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	err := h.userService.CreateLeader(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create leader"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Leader created successfully"})
}

func (h *UserHandler) UpdateLeader(c *gin.Context) {
	idParam := c.Param("id")
	var id int
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}
	var req model.KaryawanLeader
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	req.ID = id
	err = h.userService.UpdateLeader(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leader"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Leader updated successfully"})
}

func (h *UserHandler) DeleteLeader(c *gin.Context) {
	idParam := c.Param("id")
	var id int
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}
	err = h.userService.DeleteLeader(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete leader"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Leader deleted successfully"})
}

// User Account CRUD
func (h *UserHandler) GetAllUserAccounts(c *gin.Context) {
	list, err := h.userService.GetAllUserAccounts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user accounts list"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *UserHandler) CreateUserAccount(c *gin.Context) {
	var req model.UserAccountCRUD
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	req.ID = uuid.New().String()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	err := h.userService.CreateUserAccount(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User account created successfully"})
}

func (h *UserHandler) UpdateUserAccount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}
	var req model.UserAccountCRUD
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	req.ID = id
	req.UpdatedAt = time.Now()

	err := h.userService.UpdateUserAccount(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User account updated successfully"})
}

func (h *UserHandler) DeleteUserAccount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}
	err := h.userService.DeleteUserAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User account deleted successfully"})
}

func (h *UserHandler) RegisterTenant(c *gin.Context) {
	var req model.RegisterTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input data tidak valid: " + err.Error()})
		return
	}

	err := h.userService.RegisterTenant(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant dan admin berhasil didaftarkan"})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req model.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input data tidak valid: " + err.Error()})
		return
	}

	err := h.userService.ForgotPasswordReset(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diperbarui"})
}

