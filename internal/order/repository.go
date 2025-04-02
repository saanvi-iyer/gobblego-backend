package order

import (
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateOrder(db *gorm.DB, order *models.Order) error
	GetOrderByID(db *gorm.DB, id string) (*models.Order, error)
	ListAllOrders(db *gorm.DB) ([]models.Order, error)
	UpdateOrderStatus(db *gorm.DB, orderID string, status string) error
}
