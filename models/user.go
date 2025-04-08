package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `gorm:"primaryKey;type:uuid" json:"user_id"`
	CartID    uuid.UUID `gorm:"type:uuid;index"      json:"cart_id"`
	IsLeader  bool      `gorm:"default:false"        json:"is_leader"`
	CreatedAt time.Time `                            json:"created_at"`
	UserName  string    `gorm:"not null"             json:"user_name"`

	Cart Cart `gorm:"foreignKey:CartID;references:CartID" json:"-"`
}
