package payment

import (
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreatePayment(db *gorm.DB, payment *models.Payment) error
	GetPaymentByID(db *gorm.DB, id string) (*models.Payment, error)
	GetPaymentsByCartID(db *gorm.DB, cartID string) ([]models.Payment, error)
	UpdatePaymentStatus(db *gorm.DB, paymentID string, status string) error
	UpdateRazorpayPaymentID(db *gorm.DB, paymentID, razorpayPaymentID string) error
}
