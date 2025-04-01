package cart

import (
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type Repository interface {
	AddToCart(db *gorm.DB, cartID uuid.UUID, item *models.CartItem) error
	GetCartByID(db *gorm.DB, id string) (*models.Cart, error)
	ListCartItems(db *gorm.DB, cartID uuid.UUID) ([]models.CartItem, error)
	GetAllCarts(db *gorm.DB) ([]models.Cart, error)
}
