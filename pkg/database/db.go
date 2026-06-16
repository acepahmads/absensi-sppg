// internal/database/database.go

package database

import (
	"fmt"
	"log"
	"time"

	"absensi-sppg/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDB(cfg *config.Config) (*sqlx.DB, error) {
	// Create database if not exists
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort,
	)
	dbRaw, err := sqlx.Connect("mysql", dsnWithoutDB)
	if err == nil {
		_, errCreate := dbRaw.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.DBName))
		if errCreate == nil {
			log.Printf("Database '%s' siap atau berhasil dibuat.", cfg.DBName)
		}
		dbRaw.Close()
	}

	var db *sqlx.DB
	for {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
		)

		db, err = sqlx.Connect("mysql", dsn)
		if err != nil {
			log.Printf("Gagal koneksi ke database: %v", err)
			sleepDuration := 1 // Detik
			log.Printf("Mencoba kembali dalam %d detik...", sleepDuration)
			time.Sleep(time.Duration(sleepDuration) * time.Second)
			continue
		} else {
			log.Println("Berhasil terkoneksi ke database")
			break
		}
	}

	err = EnsureSchema(db)
	if err != nil {
		log.Printf("Gagal inisialisasi schema database: %v", err)
	}

	return db, nil
}
