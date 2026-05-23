package models

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	Name           string  `gorm:"type:varchar(100);not null" json:"name"`
	Type           string  `gorm:"type:varchar(50);not null;default:'outlet'" json:"type"`
	Address        string  `gorm:"type:varchar(255);not null" json:"address"`
	City           string  `gorm:"type:varchar(100);not null" json:"city"`
	Province       string  `gorm:"type:varchar(100);not null" json:"province"`
	Latitude       float64 `gorm:"type:double precision;not null" json:"latitude"`
	Longitude      float64 `gorm:"type:double precision;not null" json:"longitude"`
	GeofenceRadius float64 `gorm:"not null" json:"geofence_radius"`
	MinimumStaffing int    `gorm:"not null;default:0" json:"minimum_staffing"`
}
