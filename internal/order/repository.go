package order

import (
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateOrder(db *gorm.DB, order *models.Order, items []models.OrderItem) error
	GetOrderByID(db *gorm.DB, id uuid.UUID) (*models.Order, error)
	GetOrdersByCartID(db *gorm.DB, cartID uuid.UUID) ([]models.Order, error)
	UpdateOrderStatus(db *gorm.DB, id uuid.UUID, status string) error
	GetOrderItems(db *gorm.DB, orderID uuid.UUID) ([]models.OrderItem, error)
}
