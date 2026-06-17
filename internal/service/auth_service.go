package service

import (
	"absensi-sppg/internal/repository"
	"absensi-sppg/pkg/utils"
	"errors"
)

type AuthService interface {
	Login(email, password string) (string, string, string, string, string, int, int, error)
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{authRepo}
}

func (s *authService) Login(email, password string) (string, string, string, string, string, int, int, error) {
	user, _, err := s.authRepo.FindUserByEmail(email)
	if err != nil {
		return "", "", "", "", "", 0, 0, errors.New("Email/Akun tidak ditemukan")
	}
	_, err = s.authRepo.CheckStatus(user.Email)
	if err != nil {
		return "", "", "", "", "", 0, 0, errors.New("Email/Akun anda belum aktif")
	}
	if !utils.CheckPassword(user.Password, password) {
		return "", "", "", "", "", 0, 0, errors.New("Password salah")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.TenantID)
	if err != nil {
		return "", "", "", "", "", 0, 0, err
	}

	namaKaryawan := ""
	if user.NamaKaryawan.Valid {
		namaKaryawan = user.NamaKaryawan.String
	}

	jabatan := ""
	if user.Jabatan.Valid {
		jabatan = user.Jabatan.String
	}

	return user.ID, token, user.Role, namaKaryawan, jabatan, user.IDKaryawan, user.IDLeader, nil
}
