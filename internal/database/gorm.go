package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"absensi-sppg/internal/config"
	"absensi-sppg/internal/model"
)

var GormDB *gorm.DB

func InitGorm(cfg *config.Config) (*gorm.DB, error) {
	// Create database if not exists
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort,
	)
	dbRaw, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})
	if err != nil {
		log.Printf("Gagal koneksi ke MySQL untuk membuat database: %v", err)
	} else {
		createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.DBName)
		if err := dbRaw.Exec(createDBQuery).Error; err != nil {
			log.Printf("Gagal membuat database %s: %v", cfg.DBName, err)
		} else {
			log.Printf("Database '%s' siap atau berhasil dibuat.", cfg.DBName)
		}
		sqlDB, err := dbRaw.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	log.Println("Migrasi database dengan GORM...")
	fmt.Println(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database dengan GORM: %v", err)
		return nil, err
	}

	// AutoMigrate akan membuat atau update tabel berdasarkan model
	err = db.AutoMigrate(
		&model.Tenant{},
		&model.UserAccount{},
		&model.UserInfo{},
	)
	if err != nil {
		log.Fatalf("Gagal migrasi: %v", err)
		return nil, err
	}

	log.Println("Migrasi berhasil dengan GORM")
	GormDB = db
	return db, nil
}
