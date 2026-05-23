package repositories

import (
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/config"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
)

type AssetRepository struct {
}

func NewAssetRepository() *AssetRepository {
	return &AssetRepository{}
}

func (r *AssetRepository) CreateAsset(asset *models.Asset) error {
	return config.DB.Create(asset).Error
}

func (r *AssetRepository) GetAssetByID(id uint) (*models.Asset, error) {
	var asset models.Asset
	if err := config.DB.Preload("User").First(&asset, id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepository) GetAllAssets() ([]models.Asset, error) {
	var assets []models.Asset
	if err := config.DB.Preload("User").Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *AssetRepository) UpdateAsset(asset *models.Asset) error {
	return config.DB.Save(asset).Error
}

func (r *AssetRepository) DeleteAsset(id uint) error {
	return config.DB.Delete(&models.Asset{}, id).Error
}

func (r *AssetRepository) GetAssetsByUser(userID uint) ([]models.Asset, error) {
	var assets []models.Asset
	if err := config.DB.Where("user_id = ?", userID).Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *AssetRepository) CountUnreturnedAssetsByUser(userID uint) (int64, error) {
	var count int64
	if err := config.DB.Model(&models.Asset{}).
		Where("user_id = ? AND status = ?", userID, "borrowed").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *AssetRepository) GetBorrowedAssetsByUser(userID uint) ([]models.Asset, error) {
	var assets []models.Asset
	err := config.DB.Where("user_id = ? AND status = ?", userID, "borrowed").Find(&assets).Error
	return assets, err
}
