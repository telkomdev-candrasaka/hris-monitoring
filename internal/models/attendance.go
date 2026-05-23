package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	UserID         uint       `gorm:"not null" json:"user_id"`
	User           User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LocationID     uint       `gorm:"not null" json:"location_id"`
	Location       Location   `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Status         string     `gorm:"type:varchar(50);not null" json:"status"`
	CheckIn        time.Time  `gorm:"not null" json:"check_in"`
	CheckOut       *time.Time `json:"check_out,omitempty"`
	DeviceLatitude  float64    `gorm:"type:double precision;not null" json:"device_latitude"`
	DeviceLongitude float64    `gorm:"type:double precision;not null" json:"device_longitude"`
	SelfiePath     string     `gorm:"type:varchar(255)" json:"selfie_path,omitempty"`
}
