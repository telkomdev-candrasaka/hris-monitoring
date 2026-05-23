package models

import "gorm.io/gorm"

type MandatoryEquipment struct {
	gorm.Model
	LocationID        uint     `gorm:"not null" json:"location_id"`
	Location          Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Role              string   `gorm:"type:varchar(50);not null" json:"role"`
	EquipmentCategory string   `gorm:"type:varchar(100);not null" json:"equipment_category"`
	RequiredCondition string   `gorm:"type:varchar(50);not null;default:'Layak'" json:"required_condition"`
	IsActive          bool     `gorm:"not null;default:true" json:"is_active"`
}
