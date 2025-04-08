package cart

import (
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateCart(db *gorm.DB, cart *models.Cart) error
	GetCartByID(db *gorm.DB, id uuid.UUID) (*models.Cart, error)
	UpdateCart(db *gorm.DB, cart *models.Cart) error

	AddCartItem(db *gorm.DB, item *models.CartItem) error
	GetCartItems(db *gorm.DB, cartID uuid.UUID) ([]models.CartItem, error)
	GetCartItemByID(db *gorm.DB, id uuid.UUID) (*models.CartItem, error)
	UpdateCartItem(db *gorm.DB, item *models.CartItem) error
	DeleteCartItem(db *gorm.DB, id uuid.UUID) error
	ClearCartItems(db *gorm.DB, cartID uuid.UUID) error
	UpdateCartTotal(db *gorm.DB, cartID uuid.UUID) error
}
