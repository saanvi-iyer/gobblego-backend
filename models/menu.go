package models

import "github.com/google/uuid"

type Menu struct {
	ItemID      uuid.UUID    `gorm:"primaryKey" json:"item_id"`
	ItemName    string  `gorm:"not null" json:"item_name"`
	Price       float64 `gorm:"not null" json:"price"`
	IsAvailable bool    `gorm:"default:true" json:"is_available"`
	Category    string  `gorm:"not null" json:"category"`
	EstPrepTime int     `gorm:"not null" json:"est_prep_time"`
	Description string  `gorm:"type:text" json:"description"`
}
