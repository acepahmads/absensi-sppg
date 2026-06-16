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
