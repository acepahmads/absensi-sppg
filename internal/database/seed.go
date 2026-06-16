package database

import (
	"absensi-sppg/internal/model"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&model.UserAccount{}).Count(&count)
	if count > 0 {
		log.Println("Data user sudah ada, skip seeding.")
		return
	}

	// 1. Seed Admin
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	adminUserID := uuid.New().String()
	adminUser := model.UserAccount{
		ID:        adminUserID,
		Name:      sql.NullString{String: "admin", Valid: true},
		Email:     "admin@example.com",
		Password:  string(adminPassword),
		Role:      "SuperAdmin",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Fatalf("Gagal seeding user admin: %v", err)
	} else {
		log.Println("User admin berhasil dibuat: admin@example.com / admin123")
	}

	// Seed UserInfo for Admin
	userInfo := model.UserInfo{
		ID:           uuid.New().String(),
		UserID:       adminUserID,
		FullName:     "Kang Ujang",
		Phone:        "081234567890",
		Email:        "",
		Gender:       "Laki-laki",
		BirthDate:    time.Now(),
		Address:      "Jl. Admin No. 1",
		Password:     "",
		PhotoFileURL: "default_photo.png",
		CreatedBy:    "system",
	}

	if err := db.Create(&userInfo).Error; err != nil {
		log.Fatalf("Gagal seeding user info: %v", err)
	}
}
