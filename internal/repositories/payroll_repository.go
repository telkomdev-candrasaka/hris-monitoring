package repositories

import (
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/config"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
)

type PayrollRepository struct {
}

func NewPayrollRepository() *PayrollRepository {
	return &PayrollRepository{}
}

func (r *PayrollRepository) GetPayrollByUserMonth(userID uint, month, year int) (*models.Payroll, error) {
	var payroll models.Payroll
	if err := config.DB.Where("user_id = ? AND month = ? AND year = ?", userID, month, year).First(&payroll).Error; err != nil {
		return nil, err
	}
	return &payroll, nil
}

func (r *PayrollRepository) CreatePayroll(payroll *models.Payroll) error {
	return config.DB.Create(payroll).Error
}

func (r *PayrollRepository) UpdatePayroll(payroll *models.Payroll) error {
	return config.DB.Save(payroll).Error
}

func (r *PayrollRepository) GetPayrollHistory(userID uint) ([]models.Payroll, error) {
	var payrolls []models.Payroll
	if err := config.DB.Where("user_id = ?", userID).Order("year desc, month desc").Find(&payrolls).Error; err != nil {
		return nil, err
	}
	return payrolls, nil
}
