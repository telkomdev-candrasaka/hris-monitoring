package repositories

import (
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/config"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
)

type AttendanceRepository struct {
}

func NewAttendanceRepository() *AttendanceRepository {
	return &AttendanceRepository{}
}

func (r *AttendanceRepository) CreateAttendance(attendance *models.Attendance) error {
	return config.DB.Create(attendance).Error
}

func (r *AttendanceRepository) UpdateAttendance(attendance *models.Attendance) error {
	return config.DB.Save(attendance).Error
}

func (r *AttendanceRepository) GetAttendanceByID(id uint) (*models.Attendance, error) {
	var attendance models.Attendance
	if err := config.DB.Preload("Location").Preload("User").First(&attendance, id).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *AttendanceRepository) GetAttendancesByUser(userID uint) ([]models.Attendance, error) {
	var attendances []models.Attendance
	if err := config.DB.Preload("Location").Where("user_id = ?", userID).Order("check_in desc").Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *AttendanceRepository) GetAttendancesByUserMonth(userID uint, month, year int) ([]models.Attendance, error) {
	var attendances []models.Attendance
	if err := config.DB.Preload("Location").Where("user_id = ? AND EXTRACT(MONTH FROM check_in) = ? AND EXTRACT(YEAR FROM check_in) = ?", userID, month, year).Order("check_in asc").Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *AttendanceRepository) CountAttendanceStatusesByUserMonth(userID uint, month, year int) (int, int, int, error) {
	var total int64
	var late int64
	var absent int64

	if err := config.DB.Model(&models.Attendance{}).
		Where("user_id = ? AND EXTRACT(MONTH FROM check_in) = ? AND EXTRACT(YEAR FROM check_in) = ?", userID, month, year).
		Count(&total).Error; err != nil {
		return 0, 0, 0, err
	}

	if err := config.DB.Model(&models.Attendance{}).
		Where("user_id = ? AND status = ? AND EXTRACT(MONTH FROM check_in) = ? AND EXTRACT(YEAR FROM check_in) = ?", userID, "terlambat", month, year).
		Count(&late).Error; err != nil {
		return 0, 0, 0, err
	}

	if err := config.DB.Model(&models.Attendance{}).
		Where("user_id = ? AND status = ? AND EXTRACT(MONTH FROM check_in) = ? AND EXTRACT(YEAR FROM check_in) = ?", userID, "mangkir", month, year).
		Count(&absent).Error; err != nil {
		return 0, 0, 0, err
	}

	return int(total), int(late), int(absent), nil
}

func (r *AttendanceRepository) GetOpenAttendance(userID uint) (*models.Attendance, error) {
	var attendance models.Attendance
	if err := config.DB.Preload("Location").Where("user_id = ? AND check_out IS NULL", userID).Order("check_in desc").First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}
