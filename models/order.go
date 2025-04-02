package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     uuid.UUID       `json:"order_id" gorm:"type:uuid;primaryKey"`
	CartIDs     json.RawMessage `json:"cart_ids" gorm:"type:jsonb"`
	OrderStatus string          `json:"order_status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
