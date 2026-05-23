package services

import (
	"errors"
	"time"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/repositories"
)

type AssetService struct {
	repo *repositories.AssetRepository
}

func NewAssetService(repo *repositories.AssetRepository) *AssetService {
	return &AssetService{repo: repo}
}

func (s *AssetService) CreateAsset(name, category, serialNumber, condition, notes string) (*models.Asset, error) {
	if name == "" || category == "" || serialNumber == "" {
		return nil, errors.New("nama, kategori, dan serial number wajib diisi")
	}

	asset := &models.Asset{
		Name:         name,
		Category:     category,
		SerialNumber: serialNumber,
		Status:       "available",
		Condition:    condition,
		Notes:        notes,
	}

	if err := s.repo.CreateAsset(asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *AssetService) GetAssetByID(id uint) (*models.Asset, error) {
	return s.repo.GetAssetByID(id)
}

func (s *AssetService) GetAllAssets() ([]models.Asset, error) {
	return s.repo.GetAllAssets()
}

func (s *AssetService) BorrowAsset(assetID, userID uint) (*models.Asset, error) {
	asset, err := s.repo.GetAssetByID(assetID)
	if err != nil {
		return nil, err
	}

	if asset.Status != "available" {
		return nil, errors.New("aset tidak tersedia untuk dipinjam")
	}

	now := time.Now()
	asset.Status = "borrowed"
	asset.UserID = &userID
	asset.BorrowedAt = &now
	asset.ReturnedAt = nil

	if err := s.repo.UpdateAsset(asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *AssetService) ReturnAsset(assetID uint, condition string) (*models.Asset, error) {
	asset, err := s.repo.GetAssetByID(assetID)
	if err != nil {
		return nil, err
	}

	if asset.Status != "borrowed" {
		return nil, errors.New("aset tidak dalam status dipinjam")
	}

	now := time.Now()
	asset.Status = "available"
	asset.UserID = nil
	asset.ReturnedAt = &now
	asset.Condition = condition

	if err := s.repo.UpdateAsset(asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *AssetService) UpdateAsset(asset *models.Asset) error {
	return s.repo.UpdateAsset(asset)
}

func (s *AssetService) DeleteAsset(id uint) error {
	return s.repo.DeleteAsset(id)
}

func (s *AssetService) CanUserResign(userID uint) (bool, error) {
	count, err := s.repo.CountUnreturnedAssetsByUser(userID)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (s *AssetService) GetAssetsByUser(userID uint) ([]models.Asset, error) {
	return s.repo.GetAssetsByUser(userID)
}
