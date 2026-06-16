package database

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

func EnsureSchema(db *sqlx.DB) error {
	// 1. Create tenants table if not exists
	_, errTenant := db.Exec(`CREATE TABLE IF NOT EXISTS tenants (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		code VARCHAR(50) NOT NULL UNIQUE,
		status TINYINT DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if errTenant != nil {
		log.Printf("Warning: failed to create tenants table: %v", errTenant)
	} else {
		// Insert default tenant if none exists
		_, _ = db.Exec("INSERT IGNORE INTO tenants (id, name, code, status) VALUES (1, 'PT. Cakrawala Bima Instrument', 'cbi', 1)")
		_, _ = db.Exec("INSERT IGNORE INTO tenants (id, name, code, status) VALUES (2, 'PT. Bintang Baru', 'bintang', 1)")
	}

	// 2. Add tenant_id column to existing tables if they exist
	alterQueries := []string{
		"ALTER TABLE user_accounts ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE user_infos ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE user_karyawan ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE karyawan_leader ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE karyawan_terdaftar ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE karyawan_absensi ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE karyawan_lembur ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE karyawan_absensi_site_report ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE karyawan_daily_report ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE inventory ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE inventory_masuk ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE inventory_keluar ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
		"ALTER TABLE karyawan_holidays ADD COLUMN tenant_id INT NOT NULL DEFAULT 1",
	}
	for _, q := range alterQueries {
		_, err := db.Exec(q)
		if err != nil {
			errStr := strings.ToLower(err.Error())
			// Ignore "Duplicate column name" error (1060 in MySQL)
			if !strings.Contains(err.Error(), "1060") && !strings.Contains(errStr, "duplicate column") && !strings.Contains(errStr, "already exists") {
				log.Printf("Warning: alter query failed (%s): %v", q, err)
			}
		}
	}

	// Modify user_accounts role column to VARCHAR(50) to support all roles dynamically if table exists
	var tableExists int
	_ = db.Get(&tableExists, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'user_accounts'")
	if tableExists > 0 {
		if _, err := db.Exec("ALTER TABLE user_accounts MODIFY COLUMN role VARCHAR(50) NOT NULL"); err != nil {
			log.Printf("Warning: failed to alter user_accounts role column to VARCHAR: %v", err)
		}
	}

	// Check if user_karyawan table exists
	var exists int
	err := db.Get(&exists, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'user_karyawan'")
	if err != nil {
		return fmt.Errorf("failed to check table existence: %v", err)
	}

	if exists > 0 {
		log.Println("Tables already exist. Skipping schema initialization.")
		var countKaryawan int
		var countLeader int
		_ = db.Get(&countKaryawan, "SELECT COUNT(*) FROM user_karyawan")
		_ = db.Get(&countLeader, "SELECT COUNT(*) FROM karyawan_leader")
		if countKaryawan == 0 || countLeader == 0 {
			log.Println("user_karyawan or karyawan_leader is empty. Running database seeding...")
			if errSeed := runSeeds(db); errSeed != nil {
				log.Printf("Warning: database seeding failed: %v", errSeed)
			}
		}
		// Seed PT. Bintang Baru data for verification
		_, _ = db.Exec("INSERT IGNORE INTO tenants (id, name, code, status) VALUES (2, 'PT. Bintang Baru', 'bintang', 1)")
		_, _ = db.Exec("INSERT INTO karyawan_leader (id, nama, status, tenant_id) VALUES (100, 'Budi Santoso', 1, 2) ON DUPLICATE KEY UPDATE nama=VALUES(nama)")
		_, _ = db.Exec("INSERT INTO karyawan_terdaftar (names, tenant_id) SELECT 'Ahmad Subarjo', 2 WHERE NOT EXISTS (SELECT 1 FROM karyawan_terdaftar WHERE names = 'Ahmad Subarjo' AND tenant_id = 2)")
		return nil
	}

	log.Println("Initializing database schema...")

	// 1. Create Tables
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS karyawan_leader (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nama VARCHAR(100) NOT NULL,
			divisi VARCHAR(100) DEFAULT 'Operations',
			status TINYINT(1) DEFAULT 1,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS user_karyawan (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nama_mesin_absen VARCHAR(100) NOT NULL,
			status TINYINT(1) DEFAULT 1,
			id_leader INT NOT NULL,
			uang_makan DOUBLE DEFAULT 0,
			uang_harian DOUBLE DEFAULT 0,
			jabatan VARCHAR(100) DEFAULT NULL,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS karyawan_terdaftar (
			id INT AUTO_INCREMENT PRIMARY KEY,
			names VARCHAR(100) NOT NULL,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS karyawan_absensi (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nama VARCHAR(100) NOT NULL,
			attendance_type VARCHAR(50) DEFAULT NULL,
			jam_masuk DATETIME DEFAULT NULL,
			jam_pulang DATETIME DEFAULT NULL,
			lembur_masuk DATETIME DEFAULT NULL,
			lembur_pulang DATETIME DEFAULT NULL,
			gps_latitude DOUBLE DEFAULT 0,
			gps_longitude DOUBLE DEFAULT 0,
			location_name VARCHAR(255) DEFAULT NULL,
			start_date VARCHAR(50) DEFAULT NULL,
			end_date VARCHAR(50) DEFAULT NULL,
			duration_days DOUBLE DEFAULT 0,
			bukti_photo1 VARCHAR(255) DEFAULT NULL,
			bukti_photo2 VARCHAR(255) DEFAULT NULL,
			photo_masuk VARCHAR(255) DEFAULT NULL,
			photo_pulang VARCHAR(255) DEFAULT NULL,
			status VARCHAR(50) DEFAULT NULL,
			created_at DATETIME DEFAULT NULL,
			jumlah_potongan DOUBLE DEFAULT 0,
			keterangan TEXT DEFAULT NULL,
			keterlambatan VARCHAR(100) DEFAULT NULL,
			id_user_karyawan INT DEFAULT 0,
			keterangan_ybs TEXT DEFAULT NULL,
			validasi_atasan INT DEFAULT 0,
			hide VARCHAR(10) DEFAULT NULL,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS karyawan_lembur (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nama VARCHAR(100) NOT NULL,
			tanggal_lembur VARCHAR(50) NOT NULL,
			lembur_weekday_1 DOUBLE DEFAULT 0,
			lembur_weekday_2 DOUBLE DEFAULT 0,
			lembur_weekend_1 DOUBLE DEFAULT 0,
			lembur_weekend_2 DOUBLE DEFAULT 0,
			lembur_weekend_3 DOUBLE DEFAULT 0,
			daftar_pekerjaan TEXT DEFAULT NULL,
			bukti_persetujuan_atasan VARCHAR(255) DEFAULT NULL,
			bukti_pekerjaan TEXT DEFAULT NULL,
			jumlah_bayar DOUBLE DEFAULT 0,
			approval VARCHAR(10) DEFAULT '0',
			keterangan TEXT DEFAULT NULL,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS karyawan_absensi_site_report (
			id INT AUTO_INCREMENT PRIMARY KEY,
			id_karyawan VARCHAR(50) NOT NULL,
			nama VARCHAR(100) NOT NULL,
			jenis_report VARCHAR(50) DEFAULT NULL,
			nama_system VARCHAR(100) DEFAULT NULL,
			site VARCHAR(100) DEFAULT NULL,
			maintenance_day INT DEFAULT 0,
			jam_masuk VARCHAR(50) DEFAULT NULL,
			jam_pulang VARCHAR(50) DEFAULT NULL,
			pekerjaan_hari_ini TEXT DEFAULT NULL,
			yang_dikerjakan_esok TEXT DEFAULT NULL,
			hasil_pekerjaan TEXT DEFAULT NULL,
			kendala TEXT DEFAULT NULL,
			bukti_foto_1 VARCHAR(255) DEFAULT NULL,
			bukti_foto_2 VARCHAR(255) DEFAULT NULL,
			bukti_foto_3 VARCHAR(255) DEFAULT NULL,
			ada_penggantian_sparepart VARCHAR(50) DEFAULT NULL,
			foto_sparepart_sebelum VARCHAR(255) DEFAULT NULL,
			foto_sparepart_sesudah VARCHAR(255) DEFAULT NULL,
			calibration_attachment VARCHAR(255) DEFAULT NULL,
			submitted_at DATETIME DEFAULT NULL,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS karyawan_daily_report (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nama VARCHAR(100) NOT NULL,
			lokasi_kerja VARCHAR(255) DEFAULT NULL,
			tanggal VARCHAR(50) DEFAULT NULL,
			jam_mulai VARCHAR(50) DEFAULT NULL,
			jam_selesai VARCHAR(50) DEFAULT NULL,
			pekerjaan_list TEXT DEFAULT NULL,
			rencana_besok TEXT DEFAULT NULL,
			created_at DATETIME DEFAULT NULL,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS inventory (
			id INT AUTO_INCREMENT PRIMARY KEY,
			qr_code VARCHAR(100) NOT NULL,
			nama_barang VARCHAR(255) NOT NULL,
			kategori_barang VARCHAR(100) DEFAULT NULL,
			jenis_barang VARCHAR(100) DEFAULT NULL,
			satuan VARCHAR(50) DEFAULT NULL,
			stok_awal INT DEFAULT 0,
			stok_akhir INT DEFAULT 0,
			stok_masuk INT DEFAULT 0,
			posisi VARCHAR(255) DEFAULT NULL,
			gambar VARCHAR(255) DEFAULT NULL,
			barang_masuk INT DEFAULT 0,
			barang_keluar INT DEFAULT 0,
			harga_beli INT DEFAULT 0,
			harga_jual INT DEFAULT 0,
			keterangan TEXT DEFAULT NULL,
			created_at VARCHAR(50) DEFAULT NULL,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS inventory_masuk (
			id INT AUTO_INCREMENT PRIMARY KEY,
			qr_code VARCHAR(100) NOT NULL,
			tanggal_jam VARCHAR(50) DEFAULT NULL,
			nama_barang VARCHAR(255) DEFAULT NULL,
			kategori VARCHAR(100) DEFAULT NULL,
			jenis_barang VARCHAR(100) DEFAULT NULL,
			jumlah INT DEFAULT 0,
			keterangan TEXT DEFAULT NULL,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS inventory_keluar (
			id INT AUTO_INCREMENT PRIMARY KEY,
			qr_code VARCHAR(100) NOT NULL,
			tanggal_jam VARCHAR(50) DEFAULT NULL,
			nama_barang VARCHAR(255) DEFAULT NULL,
			kategori VARCHAR(100) DEFAULT NULL,
			jenis_barang VARCHAR(100) DEFAULT NULL,
			jumlah INT DEFAULT 0,
			keterangan TEXT DEFAULT NULL,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS karyawan_holidays (
			id INT AUTO_INCREMENT PRIMARY KEY,
			date DATE NOT NULL,
			name VARCHAR(255) NOT NULL,
			tenant_id INT NOT NULL DEFAULT 1
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	}

	for _, schema := range schemas {
		_, err := db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to execute schema: %v", err)
		}
	}

	return runSeeds(db)
}

func runSeeds(db *sqlx.DB) error {
	var err error
	// 2. Import cais.sql if it exists
	caisSqlPaths := []string{
		"D:/cbi-project-src/CBI Automation & Integrated System/cais.sql",
		"D:/cbi-project-src/CBI Automation & Integrated System/cais-update-12082025.sql",
	}

	for _, path := range caisSqlPaths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("Importing data from %s...", path)
			err = importSqlFile(db, path)
			if err != nil {
				log.Printf("Warning: failed to import SQL file %s: %v", path, err)
			}
		}
	}
	// 3. Seed Leaders
	leaders := []string{
		"Erwin Widianto",
		"DJUHARTONO",
		"DANI GUMILAR",
		"ACEP  AHMAD  SUPRIAT",
		"Sri Winardono",
		"Alif",
	}

	for i, l := range leaders {
		_, err := db.Exec("INSERT INTO karyawan_leader (id, nama, status, tenant_id) VALUES (?, ?, 1, 1) ON DUPLICATE KEY UPDATE nama=VALUES(nama)", i+1, l)
		if err != nil {
			return fmt.Errorf("failed to seed leader %s: %v", l, err)
		}
	}
	// Seed a leader for tenant 2 (PT. Bintang Baru)
	_, _ = db.Exec("INSERT INTO karyawan_leader (id, nama, status, tenant_id) VALUES (100, 'Budi Santoso', 1, 2) ON DUPLICATE KEY UPDATE nama=VALUES(nama)")

	// 4. Seed Registered Employees
	registeredEmployees := []string{
		"ABED  NEGO  SURINAND",
		"ABED  NEGO  SURINANDA",
		"ACEP  AHMAD  SUPRIAT",
		"AKHMAD  HIDAYAT",
		"ANISA  AYU  LESTARI",
		"Agung Setiawan",
		"Alfauzan",
		"Andra",
		"Angga",
		"Anna Anggriana",
		"BILLY  SUKMO  PRAKOS",
		"BILLY  SUKMO  PRAKOSO",
		"Bardan Salam",
		"Bardan salam",
		"DANI GUMILAR",
		"DEVI    ARYA  PUTRI",
		"DJUHARTONO",
		"Diana",
		"Diana eka Putri",
		"Dudi_Fadilah",
		"Erwin Widianto",
		"FIQIH  ARAFAT",
		"Fadhillah Salman",
		"FadhillahSalmanAlfarisi",
		"Fahmi Algipari",
		"Fawwa Muhammad Daffa",
		"Fawwaz Muhammad Daffa",
		"Fitra",
		"GALIH  RIYANDI    PR",
		"GALIH  RIYANDI    PRA",
		"HAIKAL",
		"HERVIAN  DICKY  S",
		"HERVIAN  DICKY  SAPUTR",
		"Haikal",
		"Heriyanto",
		"Ingka Fitra Oemardi",
		"Irfan",
		"Julia Imalatul Huda",
		"Krisna",
		"Luthfi",
		"M Abdul Fajar",
		"M Aziz Alfaris",
		"M Rizky Raditya",
		"M Shadiq Adi Nugraha",
		"M.Shadiq",
		"MEGA  MERIANA  PERMA",
		"MEGA  MERIANA  PERMATA",
		"MOCHAMAD KURNIAWAN R",
		"MUHAMAD BAGUS CAHYO",
		"MUHAMAD BAGUS CAHYO F",
		"MUHAMMAD BAGUS MAULA",
		"MUHAMMAD BAGUS MAULAN",
		"MUHAMMAD YUSUF FAHRI",
		"MUHAMMAD YUSUF FAHRI S",
		"Muhamad Nabil Fadila",
		"Muhammad Rizki",
		"Najwa Destania Azzah",
		"Novi Susanti",
		"Ramzi",
		"Ramzi Nugraha",
		"Rangga Noviansyah",
		"Rasya M Andika",
		"Regita",
		"Richard",
		"Royan",
		"Royyan Firdaus M",
		"SEPRYAN ISMAIL CHAND",
		"SEPRYAN ISMAIL CHANDR",
		"SRIYONO",
		"SUWARDI",
		"Syafri",
		"TRIA SILVIA DAMAYANT",
		"Taufik",
		"Tifanny Putri M",
		"ULFA MAESYAROH",
	}

	for _, emp := range registeredEmployees {
		_, err := db.Exec("INSERT INTO karyawan_terdaftar (names, tenant_id) SELECT ?, 1 WHERE NOT EXISTS (SELECT 1 FROM karyawan_terdaftar WHERE names = ? AND tenant_id = 1)", emp, emp)
		if err != nil {
			return fmt.Errorf("failed to seed employee %s: %v", emp, err)
		}
	}
	// Seed registered employee for tenant 2 (PT. Bintang Baru)
	_, _ = db.Exec("INSERT INTO karyawan_terdaftar (names, tenant_id) SELECT 'Ahmad Subarjo', 2 WHERE NOT EXISTS (SELECT 1 FROM karyawan_terdaftar WHERE names = 'Ahmad Subarjo' AND tenant_id = 2)")

	// 5. Seed user_karyawan and link user_accounts
	type UserMapping struct {
		Email        string
		KaryawanName string
		LeaderId     int
		UangHarian   float64
		UangMakan    float64
		Jabatan      string
	}

	mappings := []UserMapping{
		{"erwin@cbinstrument.com", "Erwin Widianto", 4, 150000, 25000, "Supervisor"},
		{"rasya@cbinstrument.com", "Rasya M Andika", 1, 100000, 25000, "Operator"},
		{"rizky@cbinstrument.com", "M Rizky Raditya", 1, 120000, 25000, "Operator"},
		{"agung@cbinstrument.com", "Agung Setiawan", 1, 100000, 25000, "Operator"},
		{"anisa@cbinstrument.com", "ANISA  AYU  LESTARI", 1, 150000, 25000, "Operator"},
		{"billy@cbinstrument.com", "BILLY  SUKMO  PRAKOSO", 1, 150000, 25000, "Operator"},
		{"diana@cbinstrument.com", "Diana eka Putri", 1, 120000, 25000, "Operator"},
		{"fawwaz@cbinstrument.com", "Fawwaz Muhammad Daffa", 1, 120000, 25000, "Operator"},
		{"ramzi@cbinstrument.com", "Ramzi Nugraha", 1, 120000, 25000, "Operator"},
		{"devi@cbinstrument.com", "DEVI    ARYA  PUTRI", 1, 150000, 25000, "Operator"},
		{"djuhartono@cbinstrument.com", "DJUHARTONO", 4, 150000, 25000, "Supervisor"},
		{"rangga@cbinstrument.com", "Rangga Noviansyah", 1, 150000, 25000, "Operator"},
		{"fiqih@cbinstrument.com", "FIQIH  ARAFAT", 1, 150000, 25000, "Operator"},
		{"galih@cbinstrument.com", "GALIH  RIYANDI    PRA", 1, 150000, 25000, "Operator"},
		{"hervian@cbinstrument.com", "HERVIAN  DICKY  SAPUTR", 1, 150000, 25000, "Operator"},
		{"novi@cbinstrument.com", "Novi Susanti", 1, 120000, 25000, "Operator"},
		{"dudi@cbinstrument.com", "Dudi_Fadilah", 1, 120000, 25000, "Operator"},
		{"rizki@cbinstrument.com", "Muhammad Rizki", 1, 100000, 25000, "Operator"},
		{"baguscahyo@cbinstrument.com", "MUHAMAD BAGUS CAHYO F", 1, 150000, 25000, "Operator"},
		{"bagusmaulana@cbinstrument.com", "MUHAMMAD BAGUS MAULAN", 1, 150000, 25000, "Operator"},
		{"yusuf@cbinstrument.com", "MUHAMMAD YUSUF FAHRI S", 1, 150000, 25000, "Operator"},
		{"akhmad@cbinstrument.com", "AKHMAD  HIDAYAT", 1, 150000, 25000, "Operator"},
		{"fitra@cbinstrument.com", "Fitra", 1, 100000, 25000, "Operator"},
		{"royan@cbinstrument.com", "Royan", 1, 100000, 25000, "Operator"},
		{"sepryan@cbinstrument.com", "SEPRYAN ISMAIL CHANDR", 1, 150000, 25000, "Operator"},
		{"sriyono@cbinstrument.com", "SRIYONO", 1, 150000, 25000, "Operator"},
		{"suwardi@cbinstrument.com", "SUWARDI", 1, 150000, 25000, "Operator"},
		{"luthfi@cbinstrument.com", "Luthfi", 1, 120000, 25000, "Operator"},
		{"regita@cbinstrument.com", "Regita", 1, 150000, 25000, "Operator"},
		{"haikal@cbinstrument.com", "Haikal", 1, 150000, 25000, "Operator"},
		{"mega@cbinstrument.com", "MEGA  MERIANA  PERMATA", 1, 150000, 25000, "Operator"},
		{"richard@cbinstrument.com", "Richard", 1, 150000, 25000, "Operator"},
		{"shadiq@cbinstrument.com", "M.Shadiq", 1, 120000, 25000, "Operator"},
		{"fadhillahsalmanalfarisi@cbinstrument.com", "FadhillahSalmanAlfarisi", 1, 120000, 25000, "Operator"},
		{"tifanny@cbinstrument.com", "Tifanny Putri M", 1, 150000, 25000, "Operator"},
		{"taufik@cbinstrument.com", "Taufik", 1, 120000, 25000, "Operator"},
		{"abed@cbinstrument.com", "ABED  NEGO  SURINANDA", 1, 150000, 25000, "Operator"},
		{"bardan@cbinstrument.com", "Bardan salam", 1, 120000, 25000, "Operator"},
		{"alfauzan@cbinstrument.com", "Alfauzan", 1, 100000, 25000, "Operator"},
		{"acep@cbinstrument.com", "ACEP  AHMAD  SUPRIAT", 4, 150000, 25000, "SuperAdmin"},
		{"dani@cbinstrument.com", "DANI GUMILAR", 4, 150000, 25000, "Supervisor"},
		{"anna@cbinstrument.com", "Anna Anggriana", 1, 150000, 25000, "Operator"},
	}

	for _, m := range mappings {
		// Insert into user_karyawan
		var kid int64
		err = db.Get(&kid, "SELECT id FROM user_karyawan WHERE nama_mesin_absen = ?", m.KaryawanName)
		if err == sql.ErrNoRows {
			res, err := db.Exec("INSERT INTO user_karyawan (nama_mesin_absen, status, id_leader, uang_harian, uang_makan, jabatan) VALUES (?, 1, ?, ?, ?, ?)", m.KaryawanName, m.LeaderId, m.UangHarian, m.UangMakan, m.Jabatan)
			if err != nil {
				return fmt.Errorf("failed to insert user_karyawan for %s: %v", m.KaryawanName, err)
			}
			kid, _ = res.LastInsertId()
		}

		// Update user_accounts if the account exists
		_, err = db.Exec("UPDATE user_accounts SET id_karyawan = ?, id_leader = ?, name = ? WHERE email = ?", kid, m.LeaderId, m.KaryawanName, m.Email)
		if err != nil {
			return fmt.Errorf("failed to update user_account for email %s: %v", m.Email, err)
		}
	}

	log.Println("Database schema initialization and seeding complete.")
	return nil
}

func importSqlFile(db *sqlx.DB, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	content := string(contentBytes)
	lines := strings.Split(content, "\n")

	var queries []string
	var currentQuery []string
	skipTableBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for target tables to skip dropping and creating
		if strings.HasPrefix(trimmed, "-- Table structure for user_accounts") ||
			strings.HasPrefix(trimmed, "-- Table structure for user_infos") ||
			strings.HasPrefix(trimmed, "-- Records of user_accounts") ||
			strings.HasPrefix(trimmed, "-- Records of user_infos") {
			skipTableBlock = true
			continue
		}

		// Reset skipping if it's another table structure comment
		if strings.HasPrefix(trimmed, "-- Table structure for") &&
			!strings.HasPrefix(trimmed, "-- Table structure for user_accounts") &&
			!strings.HasPrefix(trimmed, "-- Table structure for user_infos") {
			skipTableBlock = false
		}

		if skipTableBlock {
			continue
		}

		if trimmed == "" || strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "/*") {
			continue
		}

		currentQuery = append(currentQuery, line)

		if strings.HasSuffix(trimmed, ";") {
			queries = append(queries, strings.Join(currentQuery, "\n"))
			currentQuery = nil
		}
	}

	// Run all queries
	for _, query := range queries {
		trimmedQuery := strings.TrimSpace(query)
		if trimmedQuery == "" {
			continue
		}
		trimmedQuery = strings.ReplaceAll(trimmedQuery, "utf8mb4_0900_ai_ci", "utf8mb4_unicode_ci")
		trimmedQuery = strings.ReplaceAll(trimmedQuery, "utf8_0900_ai_ci", "utf8_unicode_ci")
		_, err := db.Exec(trimmedQuery)
		if err != nil {
			// Ignore errors like table already exists since we use CREATE TABLE
			if !strings.Contains(err.Error(), "already exists") {
				log.Printf("SQL import error in query: %s\nError: %v", trimmedQuery[:valMin(len(trimmedQuery), 100)], err)
			}
		}
	}

	return nil
}

func valMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
