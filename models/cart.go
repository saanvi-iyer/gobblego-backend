package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	CartID        uuid.UUID `gorm:"primaryKey;type:uuid"       json:"cart_id"`
	PaymentStatus string    `gorm:"not null;default:'pending'" json:"payment_status"`
	BillAmount    float64   `gorm:"not null;default:0"         json:"bill_amount"`
	CreatedAt     time.Time `                                  json:"created_at"`
	UpdatedAt     time.Time `                                  json:"updated_at"`
}

type CartItem struct {
	CartItemID uuid.UUID `gorm:"primaryKey;type:uuid" json:"cart_item_id"`
	CartID     uuid.UUID `gorm:"type:uuid;not null"   json:"cart_id"`
	ItemID     uuid.UUID `gorm:"type:uuid;not null"   json:"item_id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"   json:"user_id"`
	Quantity   int       `gorm:"not null;default:1"   json:"quantity"`
	ItemPrice  float64   `gorm:"not null"             json:"item_price"`
	Notes      string    `                            json:"notes"`
	CreatedAt  time.Time `                            json:"created_at"`
	UpdatedAt  time.Time `                            json:"updated_at"`

	Cart Cart `gorm:"foreignKey:CartID;references:CartID" json:"-"`
	Item Menu `gorm:"foreignKey:ItemID;references:ItemID" json:"-"`
	User User `gorm:"foreignKey:UserID;references:UserID" json:"-"`
}
