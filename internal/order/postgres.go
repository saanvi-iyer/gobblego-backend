package order

import (
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type orderRepo struct{}

func NewOrderRepo() Repository {
	return &orderRepo{}
}

func (r *orderRepo) CreateOrder(db *gorm.DB, order *models.Order) error {
	return db.Create(order).Error
}

func (r *orderRepo) GetOrderByID(db *gorm.DB, id string) (*models.Order, error) {
	var order models.Order
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = db.Where("order_id = ?", uid).First(&order).Error
	return &order, err
}

func (r *orderRepo) ListAllOrders(db *gorm.DB) ([]models.Order, error) {
	var orders []models.Order
	err := db.Find(&orders).Error
	return orders, err
}

func (r *orderRepo) UpdateOrderStatus(db *gorm.DB, orderID string, status string) error {
	uid, err := uuid.Parse(orderID)
	if err != nil {
		return err
	}
	return db.Model(&models.Order{}).Where("order_id = ?", uid).Update("order_status", status).Error
}
