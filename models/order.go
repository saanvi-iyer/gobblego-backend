package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     uuid.UUID `gorm:"primaryKey;type:uuid"       json:"order_id"`
	CartID      uuid.UUID `gorm:"type:uuid;not null"         json:"cart_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"         json:"user_id"`
	Status      string    `gorm:"not null;default:'pending'" json:"status"`
	TotalAmount float64   `gorm:"not null"                   json:"total_amount"`
	CreatedAt   time.Time `                                  json:"created_at"`
	UpdatedAt   time.Time `                                  json:"updated_at"`

	Cart Cart `gorm:"foreignKey:CartID;references:CartID" json:"-"`
	User User `gorm:"foreignKey:UserID;references:UserID" json:"-"`
}

type OrderItem struct {
	OrderItemID uuid.UUID `gorm:"primaryKey;type:uuid" json:"order_item_id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null"   json:"order_id"`
	CartID      uuid.UUID `gorm:"type:uuid;not null"   json:"cart_id"`
	ItemID      uuid.UUID `gorm:"type:uuid;not null"   json:"item_id"`
	Quantity    int       `gorm:"not null"             json:"quantity"`
	Price       float64   `gorm:"not null"             json:"price"`
	Notes       string    `                            json:"notes"`

	Order Order `gorm:"foreignKey:OrderID;references:OrderID" json:"-"`
	Item  Menu  `gorm:"foreignKey:ItemID;references:ItemID"   json:"-"`
	Cart  Cart  `gorm:"foreignKey:CartID;references:CartID"   json:"-"`
}
