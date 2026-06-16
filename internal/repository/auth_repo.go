package repository

import (
	"absensi-sppg/internal/model"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindUserByEmail(email string) (*model.UserAccount, string, error)
	CheckStatus(email string) (*model.UserAccount, error)
}

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByEmail(email string) (*model.UserAccount, string, error) {
	var user model.UserAccount

	query := `
		SELECT 
			ua.*,
			COALESCE(uk.nama_mesin_absen, '') AS nama_karyawan, 
			uk.jabatan
		FROM user_accounts ua
		LEFT JOIN user_karyawan uk ON uk.id = ua.id_karyawan AND uk.tenant_id = ua.tenant_id
		WHERE ua.email = ?
		LIMIT 1
	`
	err := r.db.Get(&user, query, email)
	if err != nil {
		fmt.Println("Error finding user by email:", err)
		return nil, "", err
	}
	return &user, "", nil
}

func (r *authRepository) CheckStatus(email string) (*model.UserAccount, error) {
	var user model.UserAccount
	err := r.db.Get(&user, "SELECT * FROM user_accounts WHERE email = ? AND status = 1", email)
	if err != nil {
		// fmt.Println("Error finding user by email:", err)
		return nil, err
	}
	return &user, nil
}
