package payment

import (
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type paymentRepo struct{}

func NewPaymentRepo() Repository {
	return &paymentRepo{}
}

func (r *paymentRepo) CreatePayment(db *gorm.DB, payment *models.Payment) error {
	if payment.PaymentID == uuid.Nil {
		payment.PaymentID = uuid.New()
	}
	return db.Create(payment).Error
}

func (r *paymentRepo) GetPaymentByID(db *gorm.DB, id string) (*models.Payment, error) {
	var payment models.Payment
	pid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = db.Where("payment_id = ?", pid).First(&payment).Error
	return &payment, err
}

func (r *paymentRepo) GetPaymentsByCartID(db *gorm.DB, cartID string) ([]models.Payment, error) {
	var payments []models.Payment
	cid, err := uuid.Parse(cartID)
	if err != nil {
		return nil, err
	}
	err = db.Where("cart_id = ?", cid).Find(&payments).Error
	return payments, err
}

func (r *paymentRepo) UpdatePaymentStatus(db *gorm.DB, paymentID string, status string) error {
	pid, err := uuid.Parse(paymentID)
	if err != nil {
		return err
	}
	return db.Model(&models.Payment{}).Where("payment_id = ?", pid).Update("status", status).Error
}

func (r *paymentRepo) UpdateRazorpayPaymentID(
	db *gorm.DB,
	paymentID, razorpayPaymentID string,
) error {
	pid, err := uuid.Parse(paymentID)
	if err != nil {
		return err
	}
	return db.Model(&models.Payment{}).
		Where("payment_id = ?", pid).
		Update("razorpay_payment_id", razorpayPaymentID).
		Error
}
