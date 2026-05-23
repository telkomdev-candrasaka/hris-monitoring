package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Email      string    `gorm:"type:varchar(150);not null;uniqueIndex" json:"email"`
	Role       string    `gorm:"type:varchar(50);not null" json:"role"`
	LocationID uint      `gorm:"not null" json:"location_id"`
	Location   Location  `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	ShiftID    *uint     `json:"shift_id,omitempty"`
	Shift      *Shift    `gorm:"foreignKey:ShiftID" json:"shift,omitempty"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-"`
	BaseSalary float64   `gorm:"type:numeric(14,2);not null;default:0" json:"base_salary"`
	JoinedAt   time.Time `gorm:"autoCreateTime" json:"joined_at"`
	Payrolls   []Payroll `gorm:"foreignKey:UserID" json:"payrolls,omitempty"`
}
