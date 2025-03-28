package models

import "github.com/google/uuid"

type Table struct {
	TableID       uuid.UUID `gorm:"primaryKey" json:"table_id"`
	PaymentStatus string    `gorm:"not null" json:"payment_status"`
	BillAmount    float64   `gorm:"not null" json:"bill_amount"`
}
