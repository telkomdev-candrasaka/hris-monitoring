package services

import (
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/repositories"
)

type LocationService struct {
	repo *repositories.LocationRepository
}

func NewLocationService(repo *repositories.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) CreateLocation(location *models.Location) error {
	return s.repo.CreateLocation(location)
}

func (s *LocationService) GetLocationByID(id uint) (*models.Location, error) {
	return s.repo.GetLocationByID(id)
}

func (s *LocationService) GetAllLocations() ([]models.Location, error) {
	return s.repo.GetAllLocations()
}

func (s *LocationService) UpdateLocation(location *models.Location) error {
	return s.repo.UpdateLocation(location)
}

func (s *LocationService) DeleteLocation(id uint) error {
	return s.repo.DeleteLocation(id)
}
