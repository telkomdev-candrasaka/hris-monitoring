package services

import (
	"errors"
	"math"
	"time"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/repositories"
	"gorm.io/gorm"
)

type PayrollService struct {
	payrollRepo    *repositories.PayrollRepository
	attendanceRepo *repositories.AttendanceRepository
	userRepo       *repositories.UserRepository
}

func NewPayrollService(payrollRepo *repositories.PayrollRepository, attendanceRepo *repositories.AttendanceRepository, userRepo *repositories.UserRepository) *PayrollService {
	return &PayrollService{payrollRepo: payrollRepo, attendanceRepo: attendanceRepo, userRepo: userRepo}
}

func (s *PayrollService) CalculatePayroll(userID uint, month, year int) (*models.Payroll, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if month < 1 || month > 12 {
		return nil, errors.New("bulan tidak valid")
	}

	if year < 2000 {
		return nil, errors.New("tahun tidak valid")
	}

	payroll, err := s.payrollRepo.GetPayrollByUserMonth(userID, month, year)
	if err == nil {
		return payroll, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	totalAttendance, lateCount, absentCount, err := s.attendanceRepo.CountAttendanceStatusesByUserMonth(userID, month, year)
	if err != nil {
		return nil, err
	}

	baseSalary := user.BaseSalary
	absentPenalty := float64(absentCount) * (baseSalary / 22)
	latePenalty := float64(lateCount) * (baseSalary * 0.01)
	totalDeduction := absentPenalty + latePenalty
	coldStorageAllowance := 0.0
	if user.Location.Type == "warehouse" {
		coldStorageAllowance = baseSalary * 0.10
	}

	overtimePay := 0.0
	attendances, err := s.attendanceRepo.GetAttendancesByUserMonth(userID, month, year)
	if err != nil {
		return nil, err
	}
	for _, attendance := range attendances {
		if attendance.CheckOut == nil || user.Shift == nil {
			continue
		}
		shiftEnd, parseErr := time.Parse("15:04", user.Shift.EndTime)
		if parseErr != nil {
			continue
		}
		shiftEndTime := time.Date(attendance.CheckIn.Year(), attendance.CheckIn.Month(), attendance.CheckIn.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, attendance.CheckIn.Location())
		if user.Shift.CrossMidnight && shiftEndTime.Before(attendance.CheckIn) {
			shiftEndTime = shiftEndTime.Add(24 * time.Hour)
		}
		if attendance.CheckOut.After(shiftEndTime) {
			overtimeHours := attendance.CheckOut.Sub(shiftEndTime).Hours()
			if overtimeHours > 0 {
				hourlyRate := baseSalary / (22 * 8)
				overtimePay += math.Round(overtimeHours*hourlyRate*1.5*100) / 100
			}
		}
	}

	grossSalary := baseSalary + coldStorageAllowance + overtimePay
	netSalary := grossSalary - totalDeduction
	if netSalary < 0 {
		netSalary = 0
	}

	payroll = &models.Payroll{
		UserID:          userID,
		Month:           month,
		Year:            year,
		BaseSalary:      baseSalary,
		ColdStorageAllowance: coldStorageAllowance,
		OvertimePay:     overtimePay,
		GrossSalary:     grossSalary,
		TotalAttendance: totalAttendance,
		AbsentCount:     absentCount,
		LateCount:       lateCount,
		TotalDeduction:  totalDeduction,
		NetSalary:       netSalary,
		GeneratedAt:     time.Now(),
	}

	if err := s.payrollRepo.CreatePayroll(payroll); err != nil {
		return nil, err
	}

	return payroll, nil
}

func (s *PayrollService) GetPayrollHistory(userID uint) ([]models.Payroll, error) {
	return s.payrollRepo.GetPayrollHistory(userID)
}

func (s *PayrollService) GetPayrollByUserMonth(userID uint, month, year int) (*models.Payroll, error) {
	return s.payrollRepo.GetPayrollByUserMonth(userID, month, year)
}
