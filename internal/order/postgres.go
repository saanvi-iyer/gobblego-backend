package order

import (
	"time"

	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type orderRepo struct{}

func NewOrderRepo() Repository {
	return &orderRepo{}
}

func (r *orderRepo) CreateOrder(db *gorm.DB, order *models.Order, items []models.OrderItem) error {
	tx := db.Begin()

	if order.OrderID == uuid.Nil {
		order.OrderID = uuid.New()
	}
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range items {
		items[i].OrderID = order.OrderID
		items[i].OrderItemID = uuid.New()
		items[i].CartID = order.CartID

		if err := tx.Create(&items[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *orderRepo) GetOrderByID(db *gorm.DB, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := db.Where("order_id = ?", id).First(&order).Error
	return &order, err
}

func (r *orderRepo) GetOrdersByCartID(db *gorm.DB, cartID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	err := db.Preload("OrderItems.Item").
		Where("cart_id = ?", cartID).
		Order("created_at desc").
		Find(&orders).
		Error
	return orders, err
}

func (r *orderRepo) UpdateOrderStatus(db *gorm.DB, id uuid.UUID, status string) error {
	return db.Model(&models.Order{}).Where("order_id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
}

func (r *orderRepo) GetOrderItems(db *gorm.DB, orderID uuid.UUID) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := db.Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}

func (r *orderRepo) GetTotalAmountForPendingOrders(db *gorm.DB, cartID uuid.UUID) (float64, error) {
	var totalAmount float64
	err := db.Model(&models.Order{}).
		Where("cart_id = ? AND status = ?", cartID, "pending").
		Select("SUM(total_amount)").
		Row().
		Scan(&totalAmount)
	return totalAmount, err
}

func (r *orderRepo) GetAllOrders(db *gorm.DB) ([]models.Order, error) {
	var orders []models.Order
	err := db.Preload("OrderItems.Item").Order("created_at desc").Find(&orders).Error
	return orders, err
}
