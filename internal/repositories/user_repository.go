package repositories

import (
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/config"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := config.DB.Preload("Location").Preload("Shift").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := config.DB.Preload("Location").Preload("Shift").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}

func (r *UserRepository) DeleteUser(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}

func (r *UserRepository) CountUsersByLocation(locationID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.User{}).Where("location_id = ?", locationID).Count(&count).Error
	return count, err
}
