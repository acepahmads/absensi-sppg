package service

import (
	"absensi-sppg/internal/repository"
	"absensi-sppg/pkg/utils"
	"context"
	"errors"
	"fmt"

	"absensi-sppg/internal/model"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return UserService{
		repo: repo,
	}
}

func (s *UserService) RegisterUser(userp model.RegisterAccount) (string, error) {
	// Check if email already exists
	existingUser, err := s.repo.FindByEmail(userp.Email)
	if err != nil {
		fmt.Println("Error finding user by email:", err)
		return "", err
	}
	if existingUser != "" {
		fmt.Println("Email already exists")
		return "", errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(userp.Password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return "", err
	}

	// Create the user
	userp.Password = hashedPassword
	return s.repo.Create(&userp)
}

func (s *UserService) GetUserInfoByID(id string) (*model.UserAccount, error) {
	return s.repo.FindUserInfoByID(id)
}

func (s *UserService) Registered() ([]string, error) {
	return s.repo.Registered()
}

func (s *UserService) GetLeaders() ([]string, error) {
	return s.repo.GetLeaders()
}

func (s *UserService) GetAllUserKaryawan(ctx context.Context) ([]model.UserKaryawan, error) {
	return s.repo.GetAllUserKaryawan(ctx)
}

func (s *UserService) CreateUserKaryawan(ctx context.Context, uk *model.UserKaryawan) error {
	return s.repo.CreateUserKaryawan(ctx, uk)
}

func (s *UserService) UpdateUserKaryawan(ctx context.Context, uk *model.UserKaryawan) error {
	return s.repo.UpdateUserKaryawan(ctx, uk)
}

func (s *UserService) DeleteUserKaryawan(ctx context.Context, id int) error {
	return s.repo.DeleteUserKaryawan(ctx, id)
}

func (s *UserService) GetLeadersList(ctx context.Context) ([]model.KaryawanLeader, error) {
	return s.repo.GetLeadersList(ctx)
}

// Leader CRUD
func (s *UserService) GetAllLeaders(ctx context.Context) ([]model.KaryawanLeader, error) {
	return s.repo.GetAllLeaders(ctx)
}

func (s *UserService) CreateLeader(ctx context.Context, leader *model.KaryawanLeader) error {
	return s.repo.CreateLeader(ctx, leader)
}

func (s *UserService) UpdateLeader(ctx context.Context, leader *model.KaryawanLeader) error {
	return s.repo.UpdateLeader(ctx, leader)
}

func (s *UserService) DeleteLeader(ctx context.Context, id int) error {
	return s.repo.DeleteLeader(ctx, id)
}

// User Account CRUD
func (s *UserService) GetAllUserAccounts(ctx context.Context) ([]model.UserAccountCRUD, error) {
	return s.repo.GetAllUserAccounts(ctx)
}

func (s *UserService) CreateUserAccount(ctx context.Context, ua *model.UserAccountCRUD) error {
	if ua.Password != "" {
		hashed, err := utils.HashPassword(ua.Password)
		if err != nil {
			return err
		}
		ua.Password = hashed
	}
	return s.repo.CreateUserAccount(ctx, ua)
}

func (s *UserService) UpdateUserAccount(ctx context.Context, ua *model.UserAccountCRUD) error {
	if ua.Password != "" {
		hashed, err := utils.HashPassword(ua.Password)
		if err != nil {
			return err
		}
		ua.Password = hashed
	}
	return s.repo.UpdateUserAccount(ctx, ua)
}

func (s *UserService) DeleteUserAccount(ctx context.Context, id string) error {
	return s.repo.DeleteUserAccount(ctx, id)
}
