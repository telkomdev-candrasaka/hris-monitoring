package repositories

import (
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/config"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
)

type LocationRepository struct {
}

func NewLocationRepository() *LocationRepository {
	return &LocationRepository{}
}

func (r *LocationRepository) CreateLocation(location *models.Location) error {
	return config.DB.Create(location).Error
}

func (r *LocationRepository) GetLocationByID(id uint) (*models.Location, error) {
	var location models.Location
	if err := config.DB.First(&location, id).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *LocationRepository) GetAllLocations() ([]models.Location, error) {
	var locations []models.Location
	if err := config.DB.Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *LocationRepository) UpdateLocation(location *models.Location) error {
	return config.DB.Save(location).Error
}

func (r *LocationRepository) DeleteLocation(id uint) error {
	return config.DB.Delete(&models.Location{}, id).Error
}
