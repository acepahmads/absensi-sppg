package main

import (
	"log"

	"absensi-sppg/internal/config"
	"absensi-sppg/internal/database"
	pkgdb "absensi-sppg/pkg/database"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.InitGorm(&cfg)
	if err != nil {
		log.Fatal("Installer gagal:", err)
	}

	database.SeedData(db)

	// Trigger raw schema setup and seeding
	sqlxDB, err := pkgdb.NewDB(&cfg)
	if err != nil {
		log.Fatal("Installer gagal inisialisasi SQLX database:", err)
	}
	defer sqlxDB.Close()

	log.Println("Migrasi selesai ✅")
}
