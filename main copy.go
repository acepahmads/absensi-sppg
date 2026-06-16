package main

import (
	"log"
)

func main1() {
	log.Println("Hello, World!")
}

// // cmd/main.go
// package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"

// 	"absensi-sppg/internal/config"
// 	"absensi-sppg/internal/repository"
// 	"absensi-sppg/internal/service"
// 	"absensi-sppg/pkg/database"
// 	"absensi-sppg/pkg/middleware"
// 	"absensi-sppg/pkg/utils"

// 	"absensi-sppg/internal/handler"

// 	// "github.com/gorilla/mux"
// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	// Load DB config
// 	cfg := config.LoadConfig()
// 	// Connect to DB
// 	db, err := database.NewDB(&cfg)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	authRepo := repository.NewAuthRepository(db)
// 	authService := service.NewAuthService(authRepo)
// 	authHandler := handler.NewAuthHandler(authService)

// 	userRepo := repository.NewUserRepository(db)
// 	userService := service.NewUserService(userRepo)
// 	userHandler := handler.NewUserHandler(userService)

// 	complaindRepo := repository.NewComplainRepository(db)
// 	complaindService := service.NewComplainService(complaindRepo)
// 	complaindHandler := handler.NewComplainHandler(complaindService)

// 	absensiRepo := repository.NewAbsensiRepository(db)
// 	absensiService := service.NewAbsensiService(absensiRepo)
// 	absensiHandler := handler.NewAbsensiHandler(absensiService)

// 	wd, _ := os.Getwd()
// 	fmt.Println("WORKDIR:", wd)
// 	// Public routes
// 	r := gin.Default()
// 	r.Static("/static", "/home/wqms/app-cais/cais/cais/static")
// 	r.Use(DeviceDetector()) // ← INI WAJIB

// 	r.POST("/auth/login", authHandler.Login)
// 	r.POST("/auth/logout", authHandler.Logout)
// 	// r.POST("/auth/register", authHandler.Register)
// 	r.LoadHTMLGlob("templates/*.html")
// 	r.GET("/", func(c *gin.Context) {
// 		device := utils.DetectDeviceTypeAdvanced(c.Request)
// 		fmt.Println("Device:", device)
// 		c.Set("device_type", device)

// 		if device == utils.DeviceMobile {
// 			c.HTML(200, "index_mobile.html", nil)
// 			return
// 		}

// 		c.HTML(http.StatusOK, "index.html", nil)
// 	})

// 	r.GET("/user_registration", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "user_registration.html", nil)
// 	})
// 	//make api without auth
// 	r.POST("/api/user/register", userHandler.Register)
// 	r.GET("/api/user/registered", userHandler.Registered)
// 	r.GET("/api/user/leaders", userHandler.GetLeaders)

// 	auth := r.Group("/auth")
// 	auth.Use(middleware.JWTAuth()) // jika pakai middleware auth
// 	{
// 		// auth.POST("/change-password", authHandler.ChangePassword)
// 		auth.GET("/user/profile", userHandler.GetUserInfoByID)
// 		// auth.GET("/profile/photo/:filename", userHandler.ServeProfilePhoto)
// 		auth.POST("/v1/members/register", userHandler.Register)
// 		auth.GET("/tickets", complaindHandler.GetAll)
// 		auth.GET("/absensi", absensiHandler.GetAll)
// 		auth.POST("/absensi/:id/validate", absensiHandler.UpdateStatus)
// 		auth.POST("/absensi/:id/hide", absensiHandler.UpdateHide)
// 		auth.POST("/input/absensi", absensiHandler.InputAbsensi)
// 		auth.GET("/absensi/last", absensiHandler.GetLast)
// 		auth.POST("/absensi/konfirmasi", absensiHandler.KonfirmasiAbsensi)
// 	}

// 	// r.GET("/dashboard", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 	// 	claimsRaw, exists := c.Get("claims")
// 	// 	if !exists {
// 	// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
// 	// 		return
// 	// 	}
// 	// 	claims, ok := claimsRaw.(*utils.Claims)
// 	// 	if !ok {
// 	// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
// 	// 		return
// 	// 	}
// 	// 	c.HTML(http.StatusOK, "dashboard.html", gin.H{
// 	// 		"email": claims.Email,
// 	// 	})
// 	// })

// 	r.GET("/dashboard", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 		// Ambil device
// 		device, _ := c.Get("device_type")

// 		// Ambil claims (AMAN)
// 		claimsRaw, exists := c.Get("claims")
// 		if !exists {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "Token tidak ditemukan"})
// 			return
// 		}

// 		claims, ok := claimsRaw.(*utils.Claims)
// 		if !ok {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "Token tidak valid"})
// 			return
// 		}

// 		data := gin.H{
// 			"email":  claims.Email,
// 			"device": device,
// 		}

// 		fmt.Println("device", device)
// 		if device == utils.DeviceMobile {
// 			c.HTML(200, "dashboard_mobile.html", data)
// 			return
// 		}

// 		c.HTML(200, "dashboard.html", data)
// 	})

// 	r.GET("/complain_handling", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 		claimsRaw, exists := c.Get("claims")
// 		if !exists {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
// 			return
// 		}
// 		claims, ok := claimsRaw.(*utils.Claims)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
// 			return
// 		}
// 		c.HTML(http.StatusOK, "complain_handling.html", gin.H{
// 			"email": claims.Email,
// 		})
// 	})

// 	r.GET("/rekapan_absensi", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 		claimsRaw, exists := c.Get("claims")
// 		if !exists {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
// 			return
// 		}
// 		claims, ok := claimsRaw.(*utils.Claims)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
// 			return
// 		}
// 		c.HTML(http.StatusOK, "rekapan_absensi.html", gin.H{
// 			"email": claims.Email,
// 		})
// 	})

// 	r.GET("/input_absensi", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 		claimsRaw, exists := c.Get("claims")
// 		if !exists {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
// 			return
// 		}
// 		claims, ok := claimsRaw.(*utils.Claims)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
// 			return
// 		}
// 		c.HTML(http.StatusOK, "input_absensi1.html", gin.H{
// 			"email": claims.Email,
// 		})
// 	})

// 	r.GET("/input_maintenance_da", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 		claimsRaw, exists := c.Get("claims")
// 		if !exists {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
// 			return
// 		}
// 		claims, ok := claimsRaw.(*utils.Claims)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
// 			return
// 		}
// 		c.HTML(http.StatusOK, "input_maintenance_da.html", gin.H{
// 			"email": claims.Email,
// 		})
// 	})

// 	r.GET("/rekapan_maintenance_da", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 		claimsRaw, exists := c.Get("claims")
// 		if !exists {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
// 			return
// 		}
// 		claims, ok := claimsRaw.(*utils.Claims)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
// 			return
// 		}
// 		c.HTML(http.StatusOK, "rekapan_maintenance_da.html", gin.H{
// 			"email": claims.Email,
// 		})
// 	})

// 	r.GET("/inventory", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 		claimsRaw, exists := c.Get("claims")
// 		if !exists {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
// 			return
// 		}
// 		claims, ok := claimsRaw.(*utils.Claims)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
// 			return
// 		}
// 		c.HTML(http.StatusOK, "inventory.html", gin.H{
// 			"email": claims.Email,
// 		})
// 	})

// 	r.GET("/inventory_management", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 		claimsRaw, exists := c.Get("claims")
// 		if !exists {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
// 			return
// 		}
// 		claims, ok := claimsRaw.(*utils.Claims)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
// 			return
// 		}
// 		c.HTML(http.StatusOK, "inventory_management.html", gin.H{
// 			"email": claims.Email,
// 		})
// 	})

// 	r.GET("/inventory_peminjaman", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 		claimsRaw, exists := c.Get("claims")
// 		if !exists {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
// 			return
// 		}
// 		claims, ok := claimsRaw.(*utils.Claims)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
// 			return
// 		}
// 		c.HTML(http.StatusOK, "inventory_peminjaman.html", gin.H{
// 			"email": claims.Email,
// 		})
// 	})

// 	r.GET("/konfirmasi_keterlambatan", middleware.JWTAuthDashboard(), func(c *gin.Context) {
// 		// Ambil device
// 		device, _ := c.Get("device_type")

// 		// Ambil claims (AMAN)
// 		claimsRaw, exists := c.Get("claims")
// 		if !exists {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "Token tidak ditemukan"})
// 			return
// 		}

// 		claims, ok := claimsRaw.(*utils.Claims)
// 		if !ok {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "Token tidak valid"})
// 			return
// 		}

// 		data := gin.H{
// 			"email":  claims.Email,
// 			"device": device,
// 		}

// 		fmt.Println("device", device)
// 		if device == utils.DeviceMobile {
// 			c.HTML(200, "absensi_konfirmasi.html", data)
// 			return
// 		}

// 		c.HTML(200, "dashboard_mobile.html", data)
// 	})

// 	//change without middleware.JWTAuthDashboard
// 	r.GET("/aicp_scada", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "aicp_scada.html", gin.H{
// 			"email": "ok",
// 		})
// 	})

// 	r.GET("/aicp_chemical_performance", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "aicp_chemical_performance.html", gin.H{
// 			"email": "ok",
// 		})
// 	})

// 	r.GET("/aicp_security", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "aicp_security.html", gin.H{
// 			"email": "ok",
// 		})
// 	})

// 	r.GET("/aicp_filosofi_dashboard", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "aicp_filosofi_dashboard.html", gin.H{
// 			"email": "ok",
// 		})
// 	})

// 	// Start server
// 	port := os.Getenv("APP_PORT")
// 	if port == "" {
// 		port = "2027"
// 	}
// 	log.Println("Server running at :" + port)
// 	log.Fatal(http.ListenAndServe(":"+port, r))
// }

// func ProxyGet(c *gin.Context, url string) {
// 	resp, err := http.Get(url)
// 	fmt.Println("resp", resp)
// 	if err != nil {
// 		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch external API"})
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Set Content-Type sesuai response eksternal
// 	c.Header("Content-Type", resp.Header.Get("Content-Type"))
// 	c.Status(resp.StatusCode)
// 	io.Copy(c.Writer, resp.Body)
// }

// func DeviceDetector() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		device := utils.DetectDeviceTypeAdvanced(c.Request)
// 		c.Set("device_type", device)
// 		c.Next()
// 	}
// }
