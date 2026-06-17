// cmd/main.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"absensi-sppg/internal/config"
	"absensi-sppg/internal/handler"
	"absensi-sppg/internal/repository"
	"absensi-sppg/internal/service"
	"absensi-sppg/pkg/database"
	"absensi-sppg/pkg/middleware"
	"absensi-sppg/pkg/utils"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	// ===============================
	// CONFIG & DATABASE
	// ===============================
	cfg := config.LoadConfig()

	db, err := database.NewDB(&cfg)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// ===============================
	// DEPENDENCY INJECTION
	// ===============================
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	absensiRepo := repository.NewAbsensiRepository(db)
	absensiService := service.NewAbsensiService(absensiRepo)
	absensiHandler := handler.NewAbsensiHandler(absensiService)

	// ===============================
	// GIN SETUP
	// ===============================
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // atau spesifik: http://localhost:5173
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(DeviceDetector())
	// r.Use(SecurityHeaders())

	if _, err := os.Stat("./static"); err == nil {
		r.Static("/static", "./static")
	} else {
		r.Static("/static", "/home/cais/apps/cais/cais/static")
	}
	r.Delims("{[{", "}]}")
	r.LoadHTMLGlob("templates/*.html")

	// ===============================
	// PUBLIC ROUTES
	// ===============================
	r.GET("/", func(c *gin.Context) {
		device, _ := c.Get("device_type")
		token, err := c.Cookie("authToken")
		fmt.Println("tok", token)
		if err == nil && token != "" {
			if device == utils.DeviceMobile {
				c.Redirect(http.StatusFound, "/input_absensi")
			} else {
				c.Redirect(http.StatusFound, "/dashboard")
			}
			return
		}

		if device == utils.DeviceMobile {
			c.HTML(200, "index_mobile.html", nil)
			return
		}
		c.HTML(200, "index.html", nil)
	})

	r.GET("/user_registration", middleware.JWTAuthDashboard(), func(c *gin.Context) {
		claims, ok := utils.GetClaims(c)
		if !ok {
			c.Redirect(http.StatusFound, "/")
			return
		}
		userAccount, err := userService.GetUserInfoByID(claims.UserID)
		if err != nil || (userAccount.Role != "SuperAdmin" && userAccount.Role != "Manager") {
			c.Redirect(http.StatusFound, "/dashboard")
			return
		}
		c.HTML(200, "user_registration.html", nil)
	})
	r.GET("/tenant_registration", func(c *gin.Context) {
		c.HTML(200, "tenant_registration.html", nil)
	})
	r.GET("/forgot_password", func(c *gin.Context) {
		c.HTML(200, "forgot_password.html", nil)
	})
	r.GET("/absen", func(c *gin.Context) {
		c.HTML(200, "absen.html", nil)
	})
	r.GET("/absen_test", func(c *gin.Context) {
		c.HTML(200, "absen.html", nil)
	})

	// ===============================
	// AUTH PUBLIC API
	// ===============================
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/logout", authHandler.Logout)

	// ===============================
	// PUBLIC API (NO AUTH)
	// ===============================
	api := r.Group("/api")
	{
		api.POST("/user/register", middleware.JWTAuthDashboard(), userHandler.Register)
		api.POST("/user/forgot-password", userHandler.ForgotPassword)
		api.POST("/tenant/register", userHandler.RegisterTenant)
		api.GET("/user/registered", userHandler.Registered)
		api.GET("/user/leaders", userHandler.GetLeaders)
		api.GET("/tenants", userHandler.GetTenants)
		api.GET("/indHolidays", absensiHandler.GetIndHolidays)
		api.POST("/absen", absensiHandler.InputAbsenMesin)
		api.POST("/verify-face", func(c *gin.Context) {
			type VerifyRequest struct {
				Image string `json:"image"`
			} // utils.FaceVerify(c, req.PhotoBase64)
			var req VerifyRequest

			// Parse JSON
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "invalid JSON"})
				return
			}

			// fmt.Printf("Image Data: %s\n", req.Image)
			imageData := req.Image
			if imageData == "" {
				c.JSON(400, gin.H{"error": "No image provided"})
				return
			}
			utils.FaceVerify(c, imageData, db)
		})
	}

	// ===============================
	// AUTH API (JWT)
	// ===============================
	auth := r.Group("/auth")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/user/profile", userHandler.GetUserInfoByID)
		auth.POST("/v1/members/register", userHandler.Register)

		// User Karyawan CRUD endpoints
		auth.GET("/user-karyawan", userHandler.GetAllUserKaryawan)
		auth.POST("/user-karyawan", userHandler.CreateUserKaryawan)
		auth.PUT("/user-karyawan/:id", userHandler.UpdateUserKaryawan)
		auth.DELETE("/user-karyawan/:id", userHandler.DeleteUserKaryawan)
		auth.GET("/leaders-list", userHandler.GetLeadersList)

		// Leader CRUD endpoints
		auth.GET("/setup/leaders", userHandler.GetAllLeaders)
		auth.POST("/setup/leaders", userHandler.CreateLeader)
		auth.PUT("/setup/leaders/:id", userHandler.UpdateLeader)
		auth.DELETE("/setup/leaders/:id", userHandler.DeleteLeader)

		// User Account CRUD endpoints
		auth.GET("/setup/users", userHandler.GetAllUserAccounts)
		auth.POST("/setup/users", userHandler.CreateUserAccount)
		auth.PUT("/setup/users/:id", userHandler.UpdateUserAccount)
		auth.DELETE("/setup/users/:id", userHandler.DeleteUserAccount)

		auth.GET("/absensi", absensiHandler.GetAll)
		auth.GET("/dashboard/stats", absensiHandler.GetDashboardStats)
		auth.GET("/absensi/statistik", absensiHandler.GetAttendanceStats)
		auth.GET("/absensi/statistik/individu", absensiHandler.GetIndividualStats)
		auth.GET("/absensi/last", absensiHandler.GetLast)
		auth.GET("/absensi/perhitungan", absensiHandler.GetAllPerhitungan)
		auth.GET("/absensi/karyawan", absensiHandler.GetAbsensiByKaryawan)
		auth.GET("/absensi/rekap-karyawan", absensiHandler.GetRekapAbsensiByKaryawan)
		auth.GET("/absensi/saya", absensiHandler.GetAbsensiSaya)
		auth.GET("/absensi/site-report", absensiHandler.GetSiteReports)
		auth.GET("/absensi/list-lembur", absensiHandler.GetLemburList)
		auth.GET("/absensi/lembur-detail", absensiHandler.GetLemburDetail)
		auth.GET("/absensi/daily-report", absensiHandler.GetDailyReports)
		auth.GET("/absensi/daily-report/:id", absensiHandler.GetDailyReportByID)
		auth.POST("/input/absensi", absensiHandler.InputAbsensi)
		auth.POST("/absensi/:id/validate", absensiHandler.UpdateStatus)
		auth.POST("/absensi/:id/hide", absensiHandler.UpdateHide)
		auth.POST("/absensi/konfirmasi", absensiHandler.KonfirmasiAbsensi)
		auth.POST("/absensi/input-lembur", absensiHandler.InputLembur)
		auth.POST("/absensi/inputleader", absensiHandler.InputAbsensiLeader)
		auth.POST("/absensi/site-report", absensiHandler.InputSiteReport)
		auth.POST("/absensi/lembur-approve", absensiHandler.ApproveLembur)
		auth.POST("/absensi/lembur-reject", absensiHandler.RejectLembur)
		auth.POST("/absensi/lembur-revise", absensiHandler.ReviseLembur)
		auth.POST("/absensi/daily-report", absensiHandler.InputDailyReport)
		auth.DELETE("/absensi/inputleader", absensiHandler.DeleteAbsensiLeader)
		auth.PUT("/absensi/inputleader", absensiHandler.UpdateAbsensiLeader)
	}

	// ===============================
	// DASHBOARD ROUTES (JWT + DEVICE)
	// ===============================
	r.GET("/dashboard", middleware.JWTAuthDashboard(), func(c *gin.Context) {
		renderDashboard(c, "dashboard.html", "dashboard_mobile.html")
	})

	r.GET("/konfirmasi_keterlambatan", middleware.JWTAuthDashboard(), func(c *gin.Context) {
		renderDashboard(c, "absensi_konfirmasi.html", "absensi_konfirmasi.html")
	})

	// ===============================
	// DASHBOARD STATIC PAGES
	// ===============================

	dashboardPages := map[string]string{
		"/rekapan_absensi":              "absensi_rekapan.html",
		"/input_absensi":                "absensi_input.html",
		"/absensi_perhitungan":          "absensi_perhitungan.html",
		"/absensi_karyawan":             "absensi_karyawan.html",
		"/input_lembur":                 "absensi_input_lembur.html",
		"/absensi_site_report":          "absensi_site_report.html",
		"/absensi_list_site_report":     "absensi_list_site_report.html",
		"/absensi_list_lembur":          "absensi_list_lembur.html",
		"/input_maintenance_da":         "maintenance_da_input.html",
		"/rekapan_maintenance_da":       "maintenance_da_rekapan.html",
		"/absensi_daily_report":         "absensi_daily_report.html",
		"/absensi_daily_report_rekapan": "absensi_daily_report_rekapan.html",
		"/setup/user_karyawan":          "setup_user_karyawan.html",
		"/setup/leaders":                "setup_leaders.html",
		"/setup/users":                  "setup_users.html",
		"/absensi_statistik":            "absensi_statistik.html",
		"/laporan_penggajian":           "laporan_penggajian.html",
	}

	for path, tpl := range dashboardPages {
		r.GET(path, middleware.JWTAuthDashboard(), dashboardPageHandler(tpl))
	}

	// ===============================
	// START SERVER
	// ===============================
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "2027"
	}

	log.Println("Server running at :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

//
// ===============================
// HELPERS
// ===============================
//

// Render dashboard by device
func renderDashboard(c *gin.Context, desktopTpl, mobileTpl string) {
	device, _ := c.Get("device_type")

	claims, ok := utils.GetClaims(c)
	if !ok {
		return
	}

	data := gin.H{
		"email":  claims.Email,
		"device": device,
	}

	if device == utils.DeviceMobile {
		c.HTML(200, mobileTpl, data)
		return
	}

	c.HTML(200, desktopTpl, data)
}

// Static dashboard handler
func dashboardPageHandler(tpl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := utils.GetClaims(c)
		if !ok {
			return
		}
		c.HTML(200, tpl, gin.H{"email": claims.Email})
	}
}

// Device detector middleware
func DeviceDetector() gin.HandlerFunc {
	return func(c *gin.Context) {
		device := utils.DetectDeviceTypeAdvanced(c.Request)
		c.Set("device_type", device)
		c.Next()
	}
}

// Optional proxy helper
func ProxyGet(c *gin.Context, url string) {
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch external API"})
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func HomeHandler(c *gin.Context) {
	token, err := c.Cookie("authToken")

	if err == nil && token != "" {
		// masih login → ke dashboard
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	// belum login → tampilkan halaman login
	c.HTML(http.StatusOK, "login.html", nil)
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("X-XSS-Protection", "1; mode=block")

		// CSP – sangat penting
		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data:; "+
				"connect-src 'self';")

		c.Next()
	}
}
