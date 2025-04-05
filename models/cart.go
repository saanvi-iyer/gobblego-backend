package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CartItem struct {
	ItemID   uuid.UUID   `json:"item_id"`
	UserIDs  []uuid.UUID `json:"user_ids"`
	Quantity int         `json:"quantity"`
}

type Cart struct {
	CartID        uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"cart_id"`
	Items         json.RawMessage `gorm:"type:jsonb"                                     json:"items"`
	PaymentStatus string          `gorm:"not null"                                       json:"payment_status"`
	BillAmount    float64         `gorm:"not null"                                       json:"bill_amount"`
}
