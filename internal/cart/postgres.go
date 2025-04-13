package cart

import (
	"time"

	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type cartRepo struct{}

func NewCartRepo() Repository {
	return &cartRepo{}
}

func (r *cartRepo) CreateCart(db *gorm.DB, cart *models.Cart) error {
	if cart.CartID == uuid.Nil {
		cart.CartID = uuid.New()
	}
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()
	return db.Create(cart).Error
}

func (r *cartRepo) GetCartByID(db *gorm.DB, id uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := db.Preload("User").Where("cart_id = ?", id).First(&cart).Error
	return &cart, err
}

func (r *cartRepo) UpdateCart(db *gorm.DB, cart *models.Cart) error {
	cart.UpdatedAt = time.Now()
	return db.Save(cart).Error
}

func (r *cartRepo) AddCartItem(db *gorm.DB, item *models.CartItem) error {
	if item.CartItemID == uuid.Nil {
		item.CartItemID = uuid.New()
	}
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	return db.Create(item).Error
}

func (r *cartRepo) GetCartItems(db *gorm.DB, cartID uuid.UUID) ([]models.CartItem, error) {
	var items []models.CartItem
	err := db.Preload("User").Where("cart_id = ?", cartID).Find(&items).Error
	return items, err
}

func (r *cartRepo) GetCartItemByID(db *gorm.DB, id uuid.UUID) (*models.CartItem, error) {
	var item models.CartItem
	err := db.Preload("User").Where("cart_item_id = ?", id).First(&item).Error
	return &item, err
}

func (r *cartRepo) UpdateCartItem(db *gorm.DB, item *models.CartItem) error {
	item.UpdatedAt = time.Now()
	return db.Save(item).Error
}

func (r *cartRepo) DeleteCartItem(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&models.CartItem{}, "cart_item_id = ?", id).Error
}

func (r *cartRepo) ClearCartItems(db *gorm.DB, cartID uuid.UUID) error {
	return db.Delete(&models.CartItem{}, "cart_id = ?", cartID).Error
}

func (r *cartRepo) UpdateCartTotal(db *gorm.DB, cartID uuid.UUID) error {
	var total float64

	rows, err := db.Raw(`
	SELECT SUM(ci.item_price * ci.quantity) 
	FROM cart_items ci 
	WHERE ci.cart_id = ?::uuid
	`, cartID.String()).Rows()

	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&total)
	}

	return db.Model(&models.Cart{}).Where("cart_id = ?", cartID).Updates(map[string]interface{}{
		"bill_amount": total,
		"updated_at":  time.Now(),
	}).Error
}
