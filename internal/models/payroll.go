package models

import (
	"time"

	"gorm.io/gorm"
)

type Payroll struct {
	gorm.Model
	UserID          uint      `gorm:"not null" json:"user_id"`
	User            User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Month           int       `gorm:"not null" json:"month"`
	Year            int       `gorm:"not null" json:"year"`
	BaseSalary      float64   `gorm:"type:numeric(14,2);not null" json:"base_salary"`
	ColdStorageAllowance float64 `gorm:"type:numeric(14,2);not null;default:0" json:"cold_storage_allowance"`
	OvertimePay     float64   `gorm:"type:numeric(14,2);not null;default:0" json:"overtime_pay"`
	GrossSalary     float64   `gorm:"type:numeric(14,2);not null;default:0" json:"gross_salary"`
	TotalAttendance int       `gorm:"not null" json:"total_attendance"`
	AbsentCount     int       `gorm:"not null" json:"absent_count"`
	LateCount       int       `gorm:"not null" json:"late_count"`
	TotalDeduction  float64   `gorm:"type:numeric(14,2);not null" json:"total_deduction"`
	NetSalary       float64   `gorm:"type:numeric(14,2);not null" json:"net_salary"`
	GeneratedAt     time.Time `gorm:"autoCreateTime" json:"generated_at"`
}
