package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     uuid.UUID       `json:"order_id" gorm:"type:uuid;primaryKey"`
	CartID      uuid.UUID       `json:"cart_id" gorm:"type:uuid"`
	Items       json.RawMessage `json:"items" gorm:"type:jsonb"`
	OrderStatus string          `json:"order_status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
