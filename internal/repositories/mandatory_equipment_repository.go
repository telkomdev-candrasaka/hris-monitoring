package repositories

import (
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/config"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
)

type MandatoryEquipmentRepository struct{}

func NewMandatoryEquipmentRepository() *MandatoryEquipmentRepository {
	return &MandatoryEquipmentRepository{}
}

func (r *MandatoryEquipmentRepository) CreateMandatoryEquipment(item *models.MandatoryEquipment) error {
	return config.DB.Create(item).Error
}

func (r *MandatoryEquipmentRepository) GetByLocationAndRole(locationID uint, role string) ([]models.MandatoryEquipment, error) {
	var items []models.MandatoryEquipment
	err := config.DB.Where("location_id = ? AND role = ? AND is_active = ?", locationID, role, true).Find(&items).Error
	return items, err
}
