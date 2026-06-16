package service

import (
	"absensi-sppg/internal/repository"
	"absensi-sppg/pkg/utils"
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
