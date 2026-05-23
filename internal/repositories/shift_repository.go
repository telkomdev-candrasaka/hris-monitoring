package repositories

import (
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/config"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
)

type ShiftRepository struct{}

func NewShiftRepository() *ShiftRepository {
	return &ShiftRepository{}
}

func (r *ShiftRepository) CreateShift(shift *models.Shift) error {
	return config.DB.Create(shift).Error
}

func (r *ShiftRepository) GetShiftByID(id uint) (*models.Shift, error) {
	var shift models.Shift
	if err := config.DB.Preload("Location").First(&shift, id).Error; err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *ShiftRepository) GetShiftsByLocation(locationID uint) ([]models.Shift, error) {
	var shifts []models.Shift
	if err := config.DB.Where("location_id = ?", locationID).Find(&shifts).Error; err != nil {
		return nil, err
	}
	return shifts, nil
}
