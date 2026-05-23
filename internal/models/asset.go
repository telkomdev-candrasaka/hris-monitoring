package models

import (
	"time"

	"gorm.io/gorm"
)

type Asset struct {
	gorm.Model
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	Category      string    `gorm:"type:varchar(100);not null" json:"category"`
	SerialNumber  string    `gorm:"type:varchar(100);uniqueIndex" json:"serial_number"`
	UserID        *uint     `json:"user_id,omitempty"`
	User          *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status        string    `gorm:"type:varchar(50);not null" json:"status"`
	BorrowedAt    *time.Time `json:"borrowed_at,omitempty"`
	ReturnedAt    *time.Time `json:"returned_at,omitempty"`
	Condition     string    `gorm:"type:varchar(50);not null;default:'Layak'" json:"condition,omitempty"`
	Notes         string    `gorm:"type:text" json:"notes,omitempty"`
}
