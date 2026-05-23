package models

import "gorm.io/gorm"

type Shift struct {
	gorm.Model
	LocationID    uint     `gorm:"not null" json:"location_id"`
	Location      Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Name          string   `gorm:"type:varchar(100);not null" json:"name"`
	StartTime     string   `gorm:"type:varchar(5);not null" json:"start_time"`
	EndTime       string   `gorm:"type:varchar(5);not null" json:"end_time"`
	CrossMidnight bool     `gorm:"not null;default:false" json:"cross_midnight"`
	GraceMinutes  int      `gorm:"not null;default:15" json:"grace_minutes"`
}
