package repository

import (
	"absensi-sppg/internal/model"
	"absensi-sppg/pkg/utils"
	"context"
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
	GetAllUserKaryawan(ctx context.Context) ([]model.UserKaryawan, error)
	CreateUserKaryawan(ctx context.Context, uk *model.UserKaryawan) error
	UpdateUserKaryawan(ctx context.Context, uk *model.UserKaryawan) error
	DeleteUserKaryawan(ctx context.Context, id int) error
	GetLeadersList(ctx context.Context) ([]model.KaryawanLeader, error)

	// Leader CRUD
	GetAllLeaders(ctx context.Context) ([]model.KaryawanLeader, error)
	CreateLeader(ctx context.Context, leader *model.KaryawanLeader) error
	UpdateLeader(ctx context.Context, leader *model.KaryawanLeader) error
	DeleteLeader(ctx context.Context, id int) error

	// User Account CRUD
	GetAllUserAccounts(ctx context.Context) ([]model.UserAccountCRUD, error)
	CreateUserAccount(ctx context.Context, ua *model.UserAccountCRUD) error
	UpdateUserAccount(ctx context.Context, ua *model.UserAccountCRUD) error
	DeleteUserAccount(ctx context.Context, id string) error
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

func (r *userRepository) GetAllUserKaryawan(ctx context.Context) ([]model.UserKaryawan, error) {
	query := `
		SELECT uk.id, uk.nama_mesin_absen, uk.status, uk.id_leader, uk.uang_makan, uk.uang_harian, uk.jabatan, COALESCE(kl.nama, '') as leader_nama
		FROM user_karyawan uk
		LEFT JOIN karyawan_leader kl ON uk.id_leader = kl.id
		ORDER BY uk.id DESC
	`
	var list []model.UserKaryawan
	err := r.db.SelectContext(ctx, &list, query)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *userRepository) CreateUserKaryawan(ctx context.Context, uk *model.UserKaryawan) error {
	query := `
		INSERT INTO user_karyawan (nama_mesin_absen, status, id_leader, uang_makan, uang_harian, jabatan)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, uk.NamaMesinAbsen, uk.Status, uk.IDLeader, uk.UangMakan, uk.UangHarian, uk.Jabatan)
	return err
}

func (r *userRepository) UpdateUserKaryawan(ctx context.Context, uk *model.UserKaryawan) error {
	query := `
		UPDATE user_karyawan
		SET nama_mesin_absen = ?, status = ?, id_leader = ?, uang_makan = ?, uang_harian = ?, jabatan = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, uk.NamaMesinAbsen, uk.Status, uk.IDLeader, uk.UangMakan, uk.UangHarian, uk.Jabatan, uk.ID)
	return err
}

func (r *userRepository) DeleteUserKaryawan(ctx context.Context, id int) error {
	query := `DELETE FROM user_karyawan WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) GetLeadersList(ctx context.Context) ([]model.KaryawanLeader, error) {
	query := `SELECT id, nama, divisi, status FROM karyawan_leader ORDER BY id ASC`
	var list []model.KaryawanLeader
	err := r.db.SelectContext(ctx, &list, query)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Leader CRUD
func (r *userRepository) GetAllLeaders(ctx context.Context) ([]model.KaryawanLeader, error) {
	query := `SELECT id, nama, divisi, status FROM karyawan_leader ORDER BY id DESC`
	var list []model.KaryawanLeader
	err := r.db.SelectContext(ctx, &list, query)
	return list, err
}

func (r *userRepository) CreateLeader(ctx context.Context, leader *model.KaryawanLeader) error {
	query := `INSERT INTO karyawan_leader (nama, divisi, status) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, leader.Nama, leader.Divisi, leader.Status)
	return err
}

func (r *userRepository) UpdateLeader(ctx context.Context, leader *model.KaryawanLeader) error {
	query := `UPDATE karyawan_leader SET nama = ?, divisi = ?, status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, leader.Nama, leader.Divisi, leader.Status, leader.ID)
	return err
}

func (r *userRepository) DeleteLeader(ctx context.Context, id int) error {
	query := `DELETE FROM karyawan_leader WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// User Account CRUD
func (r *userRepository) GetAllUserAccounts(ctx context.Context) ([]model.UserAccountCRUD, error) {
	query := `
		SELECT 
			ua.id, ua.name, ua.email, ua.role, ua.status, ua.id_karyawan, ua.id_leader, ua.created_at, ua.updated_at,
			COALESCE(uk.nama_mesin_absen, '') AS nama_karyawan,
			COALESCE(kl.nama, '') AS nama_leader
		FROM user_accounts ua
		LEFT JOIN user_karyawan uk ON ua.id_karyawan = uk.id
		LEFT JOIN karyawan_leader kl ON ua.id_leader = kl.id
		ORDER BY ua.created_at DESC
	`
	var list []model.UserAccountCRUD
	err := r.db.SelectContext(ctx, &list, query)
	return list, err
}

func (r *userRepository) CreateUserAccount(ctx context.Context, ua *model.UserAccountCRUD) error {
	query := `
		INSERT INTO user_accounts (id, name, email, password, role, status, id_karyawan, id_leader, created_at, updated_at, photos)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '[]')
	`
	_, err := r.db.ExecContext(ctx, query, ua.ID, ua.Name, ua.Email, ua.Password, ua.Role, ua.Status, ua.IDKaryawan, ua.IDLeader, ua.CreatedAt, ua.UpdatedAt)
	return err
}

func (r *userRepository) UpdateUserAccount(ctx context.Context, ua *model.UserAccountCRUD) error {
	var err error
	if ua.Password != "" {
		query := `
			UPDATE user_accounts
			SET name = ?, email = ?, password = ?, role = ?, status = ?, id_karyawan = ?, id_leader = ?, updated_at = ?
			WHERE id = ?
		`
		_, err = r.db.ExecContext(ctx, query, ua.Name, ua.Email, ua.Password, ua.Role, ua.Status, ua.IDKaryawan, ua.IDLeader, ua.UpdatedAt, ua.ID)
	} else {
		query := `
			UPDATE user_accounts
			SET name = ?, email = ?, role = ?, status = ?, id_karyawan = ?, id_leader = ?, updated_at = ?
			WHERE id = ?
		`
		_, err = r.db.ExecContext(ctx, query, ua.Name, ua.Email, ua.Role, ua.Status, ua.IDKaryawan, ua.IDLeader, ua.UpdatedAt, ua.ID)
	}
	return err
}

func (r *userRepository) DeleteUserAccount(ctx context.Context, id string) error {
	query := `DELETE FROM user_accounts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
