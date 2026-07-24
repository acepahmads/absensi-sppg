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
	Registered(tenantID int) ([]string, error)
	GetLeaders(tenantID int) ([]string, error)
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
	GetTenantIDByDeviceSN(ctx context.Context, sn string) (int, error)

	// Tenants
	GetTenants(ctx context.Context) ([]model.Tenant, error)
	CreateTenant(req *model.RegisterTenantRequest) error
	ResetPassword(email, newPassword string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(userp *model.RegisterAccount) (string, error) {
	if userp.TenantID == 0 {
		userp.TenantID = 1
	}
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

	query := `SELECT id FROM karyawan_leader WHERE nama = ? AND tenant_id = ?`

	var idLeader string
	err = tx.Get(&idLeader, query, userp.Leader, userp.TenantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Leader tidak ditemukan:", userp.Leader)
			return userp.Email, fmt.Errorf("Leader tidak ditemukan")
		}

		fmt.Println("Error getting idLeader:", err)
		return userp.Email, err
	}

	query = "INSERT INTO user_karyawan (nama_mesin_absen, status, id_leader, tenant_id) VALUES (?, ?, ?, ?)"
	_, err = tx.Exec(query, userp.Name, 1, idLeader, userp.TenantID)
	if err != nil {
		fmt.Println("Error creating userKaryawan:", err)
		return userp.Email, err
	}

	query = `SELECT id FROM user_karyawan WHERE nama_mesin_absen = ? AND tenant_id = ?`

	var idKaryawan string
	err = tx.Get(&idKaryawan, query, userp.Name, userp.TenantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Karyawan tidak ditemukan:", userp.Name)
			return userp.Email, fmt.Errorf("karyawan tidak ditemukan")
		}

		fmt.Println("Error getting idKaryawan:", err)
		return userp.Email, err
	}

	// 1️⃣ Insert ke user_accounts
	query = `INSERT INTO user_accounts (id, email, password, role, status, created_at, updated_at, photos, id_karyawan, name, id_leader, tenant_id) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(query, UserID, userp.Email, userp.Password, role, 1, time.Now(), time.Now(), photosJSON, idKaryawan, userp.Name, idLeader, userp.TenantID)
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

func (r *userRepository) Registered(tenantID int) ([]string, error) {
	query := `SELECT names FROM karyawan_terdaftar WHERE tenant_id = ? ORDER BY names ASC`
	var names []string
	err := r.db.Select(&names, query, tenantID)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (r *userRepository) GetLeaders(tenantID int) ([]string, error) {
	query := `SELECT nama FROM karyawan_leader WHERE tenant_id = ? ORDER BY id ASC`
	var names []string
	err := r.db.Select(&names, query, tenantID)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (r *userRepository) GetAllUserKaryawan(ctx context.Context) ([]model.UserKaryawan, error) {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `
		SELECT uk.id, uk.nama_mesin_absen, uk.pin_mesin, uk.status, uk.id_leader, uk.uang_makan, uk.uang_harian, uk.jabatan, COALESCE(kl.nama, '') as leader_nama
		FROM user_karyawan uk
		LEFT JOIN karyawan_leader kl ON uk.id_leader = kl.id AND kl.tenant_id = uk.tenant_id
		WHERE uk.tenant_id = ?
		ORDER BY uk.id DESC
	`
	list := []model.UserKaryawan{}
	err := r.db.SelectContext(ctx, &list, query, tenantID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *userRepository) CreateUserKaryawan(ctx context.Context, uk *model.UserKaryawan) error {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `
		INSERT INTO user_karyawan (nama_mesin_absen, pin_mesin, status, id_leader, uang_makan, uang_harian, jabatan, tenant_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, uk.NamaMesinAbsen, uk.PinMesin, uk.Status, uk.IDLeader, uk.UangMakan, uk.UangHarian, uk.Jabatan, tenantID)
	return err
}

func (r *userRepository) UpdateUserKaryawan(ctx context.Context, uk *model.UserKaryawan) error {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `
		UPDATE user_karyawan
		SET nama_mesin_absen = ?, pin_mesin = ?, status = ?, id_leader = ?, uang_makan = ?, uang_harian = ?, jabatan = ?
		WHERE id = ? AND tenant_id = ?
	`
	_, err := r.db.ExecContext(ctx, query, uk.NamaMesinAbsen, uk.PinMesin, uk.Status, uk.IDLeader, uk.UangMakan, uk.UangHarian, uk.Jabatan, uk.ID, tenantID)
	return err
}

func (r *userRepository) DeleteUserKaryawan(ctx context.Context, id int) error {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `DELETE FROM user_karyawan WHERE id = ? AND tenant_id = ?`
	_, err := r.db.ExecContext(ctx, query, id, tenantID)
	return err
}

func (r *userRepository) GetLeadersList(ctx context.Context) ([]model.KaryawanLeader, error) {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `SELECT id, nama, divisi, status FROM karyawan_leader WHERE tenant_id = ? ORDER BY id ASC`
	list := []model.KaryawanLeader{}
	err := r.db.SelectContext(ctx, &list, query, tenantID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Leader CRUD
func (r *userRepository) GetAllLeaders(ctx context.Context) ([]model.KaryawanLeader, error) {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `SELECT id, nama, divisi, status FROM karyawan_leader WHERE tenant_id = ? ORDER BY id DESC`
	list := []model.KaryawanLeader{}
	err := r.db.SelectContext(ctx, &list, query, tenantID)
	return list, err
}

func (r *userRepository) CreateLeader(ctx context.Context, leader *model.KaryawanLeader) error {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `INSERT INTO karyawan_leader (nama, divisi, status, tenant_id) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, leader.Nama, leader.Divisi, leader.Status, tenantID)
	return err
}

func (r *userRepository) UpdateLeader(ctx context.Context, leader *model.KaryawanLeader) error {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `UPDATE karyawan_leader SET nama = ?, divisi = ?, status = ? WHERE id = ? AND tenant_id = ?`
	_, err := r.db.ExecContext(ctx, query, leader.Nama, leader.Divisi, leader.Status, leader.ID, tenantID)
	return err
}

func (r *userRepository) DeleteLeader(ctx context.Context, id int) error {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `DELETE FROM karyawan_leader WHERE id = ? AND tenant_id = ?`
	_, err := r.db.ExecContext(ctx, query, id, tenantID)
	return err
}

// User Account CRUD
func (r *userRepository) GetAllUserAccounts(ctx context.Context) ([]model.UserAccountCRUD, error) {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `
		SELECT 
			ua.id, ua.name, ua.email, ua.role, ua.status, ua.id_karyawan, ua.id_leader, ua.created_at, ua.updated_at, ua.tenant_id,
			COALESCE(ua.sn_mesin, '') AS sn_mesin,
			COALESCE(uk.nama_mesin_absen, '') AS nama_karyawan,
			COALESCE(kl.nama, '') AS nama_leader
		FROM user_accounts ua
		LEFT JOIN user_karyawan uk ON ua.id_karyawan = uk.id AND uk.tenant_id = ua.tenant_id
		LEFT JOIN karyawan_leader kl ON ua.id_leader = kl.id AND kl.tenant_id = ua.tenant_id
		WHERE ua.tenant_id = ?
		ORDER BY ua.created_at DESC
	`
	list := []model.UserAccountCRUD{}
	err := r.db.SelectContext(ctx, &list, query, tenantID)
	return list, err
}

func (r *userRepository) CreateUserAccount(ctx context.Context, ua *model.UserAccountCRUD) error {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `
		INSERT INTO user_accounts (id, name, email, password, role, status, id_karyawan, id_leader, created_at, updated_at, photos, tenant_id, sn_mesin)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '[]', ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, ua.ID, ua.Name, ua.Email, ua.Password, ua.Role, ua.Status, ua.IDKaryawan, ua.IDLeader, ua.CreatedAt, ua.UpdatedAt, tenantID, ua.SNMesin)
	return err
}

func (r *userRepository) UpdateUserAccount(ctx context.Context, ua *model.UserAccountCRUD) error {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	var err error
	if ua.Password != "" {
		query := `
			UPDATE user_accounts
			SET name = ?, email = ?, password = ?, role = ?, status = ?, id_karyawan = ?, id_leader = ?, updated_at = ?, sn_mesin = ?
			WHERE id = ? AND tenant_id = ?
		`
		_, err = r.db.ExecContext(ctx, query, ua.Name, ua.Email, ua.Password, ua.Role, ua.Status, ua.IDKaryawan, ua.IDLeader, ua.UpdatedAt, ua.SNMesin, ua.ID, tenantID)
	} else {
		query := `
			UPDATE user_accounts
			SET name = ?, email = ?, role = ?, status = ?, id_karyawan = ?, id_leader = ?, updated_at = ?, sn_mesin = ?
			WHERE id = ? AND tenant_id = ?
		`
		_, err = r.db.ExecContext(ctx, query, ua.Name, ua.Email, ua.Role, ua.Status, ua.IDKaryawan, ua.IDLeader, ua.UpdatedAt, ua.SNMesin, ua.ID, tenantID)
	}
	return err
}

func (r *userRepository) GetTenantIDByDeviceSN(ctx context.Context, sn string) (int, error) {
	if sn == "" {
		return 0, errors.New("empty SN")
	}
	var tenantID int
	query := `SELECT tenant_id FROM user_accounts WHERE FIND_IN_SET(?, REPLACE(sn_mesin, ' ', '')) OR sn_mesin LIKE ? LIMIT 1`
	pattern := "%" + sn + "%"
	err := r.db.GetContext(ctx, &tenantID, query, sn, pattern)
	if err != nil {
		return 0, err
	}
	return tenantID, nil
}

func (r *userRepository) DeleteUserAccount(ctx context.Context, id string) error {
	tenantID, _ := ctx.Value("tenantID").(int)
	if tenantID == 0 {
		tenantID = 1
	}
	query := `DELETE FROM user_accounts WHERE id = ? AND tenant_id = ?`
	_, err := r.db.ExecContext(ctx, query, id, tenantID)
	return err
}

func (r *userRepository) GetTenants(ctx context.Context) ([]model.Tenant, error) {
	query := `SELECT id, name, code, status, created_at, updated_at FROM tenants WHERE status = 1 ORDER BY name ASC`
	var list []model.Tenant
	err := r.db.SelectContext(ctx, &list, query)
	return list, err
}

func (r *userRepository) CreateTenant(req *model.RegisterTenantRequest) error {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM tenants WHERE code = ?", req.TenantCode)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("kode tenant sudah terdaftar")
	}

	err = r.db.Get(&count, "SELECT COUNT(*) FROM user_accounts WHERE email = ?", req.AdminEmail)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email administrator sudah terdaftar")
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	res, err := tx.Exec("INSERT INTO tenants (name, code, status) VALUES (?, ?, 1)", req.TenantName, req.TenantCode)
	if err != nil {
		return err
	}
	tenantID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	resLeader, err := tx.Exec("INSERT INTO karyawan_leader (nama, divisi, status, tenant_id) VALUES (?, 'Management', 1, ?)", req.AdminName, tenantID)
	if err != nil {
		return err
	}
	leaderID, err := resLeader.LastInsertId()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO karyawan_terdaftar (names, tenant_id) VALUES (?, ?)", req.AdminName, tenantID)
	if err != nil {
		return err
	}

	resKaryawan, err := tx.Exec("INSERT INTO user_karyawan (nama_mesin_absen, status, id_leader, tenant_id, jabatan) VALUES (?, 1, ?, ?, 'Manager')", req.AdminName, leaderID, tenantID)
	if err != nil {
		return err
	}
	karyawanID, err := resKaryawan.LastInsertId()
	if err != nil {
		return err
	}

	userID := uuid.New().String()
	query := `INSERT INTO user_accounts (id, email, password, role, status, created_at, updated_at, photos, id_karyawan, name, id_leader, tenant_id) 
              VALUES (?, ?, ?, 'SuperAdmin', 1, NOW(), NOW(), '[]', ?, ?, ?, ?)`
	_, err = tx.Exec(query, userID, req.AdminEmail, req.AdminPassword, karyawanID, req.AdminName, leaderID, tenantID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *userRepository) ResetPassword(email, newPassword string) error {
	var user model.UserAccount
	query := `SELECT * FROM user_accounts WHERE email = ? LIMIT 1`
	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data verifikasi tidak cocok dengan akun manapun")
		}
		return err
	}

	_, err = r.db.Exec("UPDATE user_accounts SET password = ?, updated_at = NOW() WHERE id = ?", newPassword, user.ID)
	return err
}

