package repository

import (
	"absensi-sppg/internal/model"
	"absensi-sppg/pkg/utils"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user *model.RegisterAccount) (string, error)
	FindByEmail(email string) (string, error) // return user_id, error
	FindUserInfoByID(id string) (*model.UserAccount, error)
	Registered() ([]string, error)
	GetLeaders() ([]string, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(userp *model.RegisterAccount) (string, error) {
	role := "Operator"

	photoFilenames := []string{}
	for _, v := range userp.Photos {
		fmt.Println("Photo:", v)
		photoFilename, err := utils.SaveBase64Image(
			v,
			strings.ReplaceAll(userp.Email, " ", "_"),
		)
		if err != nil {
			return userp.Email, err
		}
		photoFilenames = append(photoFilenames, photoFilename)
		_, _, err = utils.RegisterFace(userp.Name, v)
		if err != nil {
			return userp.Email, err
		}
	}
	photosJSON, err := json.Marshal(photoFilenames)
	if err != nil {
		return userp.Email, err
	}

	UserID := uuid.New().String()

	// Mulai transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return userp.Email, err
	}

	// Kalau ada error, rollback otomatis
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT id FROM karyawan_leader WHERE nama = ?`

	var idLeader string
	err = tx.Get(&idLeader, query, userp.Leader)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Leader tidak ditemukan:", userp.Leader)
			return userp.Email, fmt.Errorf("Leader tidak ditemukan")
		}

		fmt.Println("Error getting idLeader:", err)
		return userp.Email, err
	}

	query = "INSERT INTO user_karyawan (nama_mesin_absen, status, id_leader) VALUES (?, ?, ?)"
	_, err = tx.Exec(query, userp.Name, 1, idLeader)
	if err != nil {
		fmt.Println("Error creating userKaryawan:", err)
		return userp.Email, err
	}

	query = `SELECT id FROM user_karyawan WHERE nama_mesin_absen = ?`

	var idKaryawan string
	err = tx.Get(&idKaryawan, query, userp.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Karyawan tidak ditemukan:", userp.Name)
			return userp.Email, fmt.Errorf("karyawan tidak ditemukan")
		}

		fmt.Println("Error getting idKaryawan:", err)
		return userp.Email, err
	}

	// 1️⃣ Insert ke user_accounts
	query = `INSERT INTO user_accounts (id, email, password, role, status, created_at, updated_at, photos, id_karyawan, name, id_leader) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(query, UserID, userp.Email, userp.Password, role, 1, time.Now(), time.Now(), photosJSON, idKaryawan, userp.Name, idLeader)
	if err != nil {
		fmt.Println("Error creating userAccount:", err)
		return userp.Email, err
	}

	// registrtaionCode := GenerateRegistrationCode(userp.Role)
	// // 2️⃣ Insert ke user_infos
	// query = `INSERT INTO user_infos
	//     (id, user_id, full_name, jabatan, phone, email, gender, birth_date, address, photo_file_url, status, created_by)
	//     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	// _, err = tx.Exec(query,
	// 	uuid.New().String(),
	// 	UserID,
	// 	userp.NamaLengkap,
	// 	userp.Jabatan,
	// 	userp.NomorHP,
	// 	userp.Email,
	// 	userp.JenisKelamin,
	// 	userp.TanggalLahir,
	// 	userp.Alamat,
	// 	filenamePhoto,
	// 	0,
	// 	"SuperAdmin",
	// )

	if err != nil {
		fmt.Println("Error creating userInfo:", err)
		return userp.Email, err
	}

	// 3️⃣ Commit jika semua sukses
	err = tx.Commit()
	if err != nil {
		return userp.Email, err
	}

	return userp.Email, nil
}

func (r *userRepository) FindByEmail(email string) (string, error) {
	query := `
		SELECT email
		FROM user_accounts
		WHERE email = ?
	`

	var userEmail string
	err := r.db.Get(&userEmail, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // ⬅️ email tidak ditemukan (NORMAL)
		}
		return "", err // ⬅️ error DB lain
	}

	return userEmail, nil
}

func (r *userRepository) FindUserInfoByID(id string) (*model.UserAccount, error) {
	query := `SELECT * FROM user_accounts WHERE id = ?`

	var UserAccount model.UserAccount
	err := r.db.Get(&UserAccount, query, id)
	if err != nil {
		fmt.Println("Error finding user_info by id:", id, "error:", err)
		return nil, err
	}

	return &UserAccount, nil
}

func GenerateRegistrationCode(role string) string {
	RegistrationCode := ""
	if role == "Superadmin" {
		RegistrationCode += "SA"
	} else if role == "Admin" {
		RegistrationCode += "AD"
	} else if role == "Operator" {
		RegistrationCode += "OP"
	} else if role == "Anggota" {
		RegistrationCode += "AG"
	}
	if RegistrationCode == "" {
		return ""
	} else {
		RegistrationCode += "-"
		RegistrationCode += fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return RegistrationCode
}

func uploadPicture(Data map[string]interface{}) (string, error) {
	base64Str := Data["data"].(string)
	filename := Data["name"].(string)

	if idx := strings.Index(base64Str, ","); idx != -1 {
		base64Str = base64Str[idx+1:]
	}

	// Decode
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return "", err
	}

	// Simpan ke file
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join("static", "uploads", "profile_photos", newFilename)

	err = os.WriteFile(savePath, data, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return "", err
	}
	return newFilename, nil
}

func (r *userRepository) Registered() ([]string, error) {
	query := `SELECT names FROM karyawan_terdaftar order by names asc`
	var names []string
	err := r.db.Select(&names, query)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (r *userRepository) GetLeaders() ([]string, error) {
	query := `SELECT nama FROM karyawan_leader order by id asc`
	var names []string
	err := r.db.Select(&names, query)
	if err != nil {
		return nil, err
	}
	return names, nil
}
