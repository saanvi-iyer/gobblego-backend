package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	PaymentID         uuid.UUID `gorm:"primaryKey;type:uuid" json:"payment_id"`
	CartID            uuid.UUID `gorm:"type:uuid;index"      json:"cart_id"`
	RazorpayOrderID   string    `gorm:"size:255"             json:"razorpay_order_id"`
	RazorpayPaymentID string    `gorm:"size:255"             json:"razorpay_payment_id"`
	Amount            float64   `                            json:"amount"`
	Currency          string    `gorm:"size:3;default:INR"   json:"currency"`
	Status            string    `gorm:"size:50"              json:"status"`
	CreatedAt         time.Time `                            json:"created_at"`
	UpdatedAt         time.Time `                            json:"updated_at"`

	Cart Cart `gorm:"foreignKey:CartID;references:CartID" json:"-"`
}
