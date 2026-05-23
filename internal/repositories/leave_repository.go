package repositories

import (
	"time"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/config"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
)

type LeaveRepository struct {
}

func NewLeaveRepository() *LeaveRepository {
	return &LeaveRepository{}
}

func (r *LeaveRepository) CreateLeave(leave *models.Leave) error {
	return config.DB.Create(leave).Error
}

func (r *LeaveRepository) GetLeaveByID(id uint) (*models.Leave, error) {
	var leave models.Leave
	if err := config.DB.Preload("User").First(&leave, id).Error; err != nil {
		return nil, err
	}
	return &leave, nil
}

func (r *LeaveRepository) GetLeavesByUser(userID uint) ([]models.Leave, error) {
	var leaves []models.Leave
	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&leaves).Error; err != nil {
		return nil, err
	}
	return leaves, nil
}

func (r *LeaveRepository) UpdateLeave(leave *models.Leave) error {
	return config.DB.Save(leave).Error
}

func (r *LeaveRepository) GetPendingLeaves() ([]models.Leave, error) {
	var leaves []models.Leave
	if err := config.DB.Where("status = ?", "pending").Order("created_at asc").Find(&leaves).Error; err != nil {
		return nil, err
	}
	return leaves, nil
}

func (r *LeaveRepository) CountApprovedLeavesByLocationAndDateRange(locationID uint, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := config.DB.Table("leaves l").
		Joins("join users u on u.id = l.user_id").
		Where("u.location_id = ? AND l.status = ? AND l.start_date <= ? AND l.end_date >= ?", locationID, "approved", endDate, startDate).
		Count(&count).Error
	return count, err
}
