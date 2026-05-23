package models

import (
	"time"

	"gorm.io/gorm"
)

type Leave struct {
	gorm.Model
	UserID       uint       `gorm:"not null" json:"user_id"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	StartDate    time.Time  `gorm:"not null" json:"start_date"`
	EndDate      time.Time  `gorm:"not null" json:"end_date"`
	LeaveType    string     `gorm:"type:varchar(100);not null" json:"leave_type"`
	Reason       string     `gorm:"type:text" json:"reason,omitempty"`
	Status       string     `gorm:"type:varchar(50);not null" json:"status"`
	DocumentPath string     `gorm:"type:varchar(255)" json:"document_path,omitempty"`
	ApprovedBy   uint       `json:"approved_by,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
}
