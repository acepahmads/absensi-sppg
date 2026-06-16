package model

import (
	"database/sql"
	"time"
)

// GORM models (dari ERD)
type UserAccount struct {
	ID           string         `gorm:"primaryKey;type:char(36)"`
	Name         sql.NullString `gorm:"size:100;not null"`
	Email        string         `gorm:"uniqueIndex;size:255;not null"`
	Password     string         `gorm:"size:255;not null"`
	Role         string         `gorm:"type:enum('SuperAdmin','Manager','Supervisor','Operator','Guest');not null"`
	Status       bool           `gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
	IDKaryawan   int            `json:"id_karyawan" db:"id_karyawan"`
	NamaKaryawan string         `json:"nama_karyawan" db:"nama_karyawan"`
	Jabatan      sql.NullString `json:"jabatan" db:"jabatan"`
	Photos       string         `json:"photos" db:"photos"`
	IDLeader     int            `json:"id_leader" db:"id_leader"`
}

type UserInfo struct {
	ID           string    `gorm:"primaryKey;type:char(36)" db:"id"`
	UserID       string    `gorm:"type:char(36);not null" json:"user_id" db:"user_id"`
	FullName     string    `gorm:"size:100;not null" json:"full_name" db:"full_name"`
	Jabatan      string    `gorm:"size:100;not null" json:"jabatan" db:"jabatan"`
	Phone        string    `gorm:"size:20;not null" json:"phone" db:"phone"`
	Email        string    `gorm:"uniqueIndex;size:255;not null" json:"email" db:"email"`
	Password     string    `gorm:"size:255;not null" json:"password" db:"password"`
	Gender       string    `gorm:"size:10;not null" json:"gender" db:"gender"`
	BirthDate    time.Time `gorm:"not null" json:"birth_date" db:"birth_date"`
	Address      string    `gorm:"type:text" json:"address" db:"address"`
	PhotoFileURL string    `json:"photo_file_url" db:"photo_file_url"`
	CreatedBy    string    `json:"created_by" db:"created_by"`
}

type UserPost struct {
	NamaLengkap     string                 `json:"namaLengkap"`
	Nik             string                 `json:"nik"`
	Email           string                 `json:"email"`
	NomorHP         string                 `json:"nomorHP"`
	JenisKelamin    string                 `json:"jenisKelamin"`
	TempatLahir     string                 `json:"tempatLahir"`
	TanggalLahir    string                 `json:"tanggalLahir"`
	Alamat          string                 `json:"alamat"`
	Role            string                 `json:"role"`
	Jabatan         string                 `json:"jabatan"`
	Password        string                 `json:"password"`
	ConfirmPassword string                 `json:"confirmPassword"`
	UploadFoto      map[string]interface{} `json:"uploadFoto"`
}

type RegisterAccount struct {
	Name     string   `json:"name" binding:"required"`
	Leader   string   `json:"leader" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required"`
	Photos   []string `json:"photos" binding:"required"`
}

type UserInfoAccount struct {
	Name     string `json:"name" binding:"required"`
	Position string `json:"position" binding:"required"`
	Photo    string `json:"photo" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UserKaryawan struct {
	ID             int     `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	NamaMesinAbsen string  `json:"nama_mesin_absen" db:"nama_mesin_absen" gorm:"size:100;not null"`
	Status         int     `json:"status" db:"status" gorm:"default:1"`
	IDLeader       int     `json:"id_leader" db:"id_leader" gorm:"not null"`
	UangMakan      float64 `json:"uang_makan" db:"uang_makan" gorm:"default:0"`
	UangHarian     float64 `json:"uang_harian" db:"uang_harian" gorm:"default:0"`
	Jabatan        string  `json:"jabatan" db:"jabatan" gorm:"size:100;default:null"`
	LeaderNama     string  `json:"leader_nama" db:"leader_nama" gorm:"-"`
}

type KaryawanLeader struct {
	ID   int    `json:"id" db:"id"`
	Nama string `json:"nama" db:"nama"`
}
