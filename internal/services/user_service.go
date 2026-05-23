package services

import (
	"errors"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        *repositories.UserRepository
	assetRepo   *repositories.AssetRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func NewUserServiceWithAsset(repo *repositories.UserRepository, assetRepo *repositories.AssetRepository) *UserService {
	return &UserService{repo: repo, assetRepo: assetRepo}
}

func (s *UserService) hashPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s *UserService) CreateUser(user *models.User) error {
	hashed, err := s.hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) UpdateUser(user *models.User) error {
	if user.Password != "" {
		hashed, err := s.hashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
	}
	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id uint) error {
	if s.assetRepo != nil {
		count, err := s.assetRepo.CountUnreturnedAssetsByUser(id)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("tidak bisa menghapus user: masih ada aset yang belum dikembalikan")
		}
	}
	return s.repo.DeleteUser(id)
}

func (s *UserService) CanResign(id uint) (bool, error) {
	if s.assetRepo != nil {
		count, err := s.assetRepo.CountUnreturnedAssetsByUser(id)
		if err != nil {
			return false, err
		}
		if count > 0 {
			return false, errors.New("tidak bisa resign: masih ada aset yang belum dikembalikan")
		}
	}
	return true, nil
}
