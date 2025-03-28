package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `gorm:"primaryKey" json:"user_id"`
	TableID   uuid.UUID `gorm:"type:uuid;index" json:"table_id"`
	IsLeader  bool      `gorm:"default:false" json:"is_leader"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `gorm:"not null" json:"user_name"`
}