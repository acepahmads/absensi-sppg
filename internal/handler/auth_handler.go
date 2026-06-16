package handler

import (
	"fmt"
	"net/http"

	"absensi-sppg/internal/service"
	"absensi-sppg/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type DeviceInfo struct {
	UserAgent        string `json:"user_agent"`
	ScreenResolution string `json:"screen_resolution"`
	Timezone         string `json:"timezone"`
	Language         string `json:"language"`
}

type RegisterRequest struct {
	ProfilePhoto         string     `json:"profile_photo" binding:"required"`
	ProfilePhotoName     string     `json:"profile_photo_name" binding:"required"`
	NamaLengkap          string     `json:"nama_lengkap" binding:"required"`
	Email                string     `json:"email" binding:"required,email"`
	NoTelepon            string     `json:"no_telepon" binding:"required"`
	Jabatan              string     `json:"jabatan"`
	Alamat               string     `json:"alamat_lengkap"`
	Password             string     `json:"password" binding:"required,min=6"`
	PasswordConfirmation string     `json:"password_confirmation" binding:"required,eqfield=Password"`
	Role                 string     `json:"role" binding:"required,oneof=pengurus warga"`
	DeviceInfo           DeviceInfo `json:"device_info"`
	CreatedAt            string     `json:"created_at"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("BIND ERROR:", err.Error()) // <== Tambahkan ini!
		utils.RespondError(c, http.StatusBadRequest, "Input tidak valid")
		return
	}
	// fmt.Println("Email:", req.Email, "Password:", req.Password)
	id, token, role, nama_karyawan, jabatan, id_karyawan, id_leader, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.SetCookie("authToken", token, 3600*24*365*10, "/", "", false, false)                     // maxAge in seconds
	c.SetCookie("nama_karyawan", nama_karyawan, 3600*24*365*10, "/", "", false, false)         // maxAge in seconds
	c.SetCookie("jabatan", jabatan, 3600*24*365*10, "/", "", false, false)                     // maxAge in seconds
	c.SetCookie("id", fmt.Sprint(id), 3600*24*365*10, "/", "", false, false)                   // maxAge in seconds
	c.SetCookie("id_karyawan", fmt.Sprint(id_karyawan), 3600*24*365*10, "/", "", false, false) // maxAge in seconds
	c.SetCookie("id_leader", fmt.Sprint(id_leader), 3600*24*365*10, "/", "", false, false)     // maxAge in seconds
	c.SetCookie("role", role, 3600*24*365*10, "/", "", false, false)                           // maxAge in seconds
	// fmt.Println("Token set in cookie:", token)

	utils.RespondJSON(c, http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": token,
		"user_data": gin.H{
			"user": req.Email,
			"role": role,
		},
		"token_expires_in": 72, // 72 hours
		"message":          "Login berhasil",
		"nama_karyawan":    nama_karyawan,
		"jabatan":          jabatan,
		"id":               id,
		"id_karyawan":      id_karyawan,
		"id_leader":        id_leader,
	})

}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Jika pakai cookie
	c.SetCookie("authToken", "", -1, "/", "", false, false)
	c.SetCookie("nama_karyawan", "", -1, "/", "", false, false)
	c.SetCookie("jabatan", "", -1, "/", "", false, false)
	c.SetCookie("id", "", -1, "/", "", false, false)
	c.SetCookie("id_karyawan", "", -1, "/", "", false, false)
	c.SetCookie("id_leader", "", -1, "/", "", false, false)
	c.SetCookie("date_konfirmasi_absen", "", -1, "/", "", false, false)
	c.SetCookie("role", "", -1, "/", "", false, false)

	// Jika hanya frontend yang simpan token (misal di localStorage), tidak perlu apa-apa
	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}
