package cart

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type cartRepo struct{}

func NewCartRepo() Repository {
	return &cartRepo{}
}

func (r *cartRepo) CreateCartItem(db *gorm.DB, cart *models.Cart, item *models.CartItem) error {
	var items []models.CartItem

	if len(cart.Items) > 0 {
		if err := json.Unmarshal(cart.Items, &items); err != nil {
			return err
		}
	}

	items = append(items, *item)

	updatedItems, err := json.Marshal(items)
	if err != nil {
		return err
	}

	cart.Items = updatedItems
	return db.Save(cart).Error
}

func (r *cartRepo) GetCartByID(db *gorm.DB, id string) (*models.Cart, error) {
	var cart models.Cart
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = db.Where("cart_id = ?", uid).First(&cart).Error
	return &cart, err
}

func (r *cartRepo) AddToCart(db *gorm.DB, cartID uuid.UUID, item *models.CartItem) error {
	var cart models.Cart

	if err := db.Where("cart_id = ?", cartID).First(&cart).Error; err != nil {
		return err
	}

	var items []models.CartItem
	if len(cart.Items) > 0 {
		if err := json.Unmarshal(cart.Items, &items); err != nil {
			return errors.New("failed to unmarshal existing items in cart")
		}
	}

	items = append(items, *item)

	updatedItems, err := json.Marshal(items)
	if err != nil {
		return err
	}

	cart.Items = updatedItems
	return db.Save(&cart).Error
}
func (r *cartRepo) ListCartItems(db *gorm.DB, cartID uuid.UUID) ([]models.CartItem, error) {
	var cart models.Cart
	err := db.Where("cart_id = ?", cartID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	var items []models.CartItem
	if len(cart.Items) > 0 {
		if err := json.Unmarshal(cart.Items, &items); err != nil {
			return nil, errors.New("failed to parse cart items")
		}
	}

	return items, nil
}

func (r *cartRepo) GetAllCarts(db *gorm.DB) ([]models.Cart, error) {
	var carts []models.Cart
	err := db.Find(&carts).Error
	if err != nil {
		return nil, err
	}

	for i := range carts {
		var items []models.CartItem
		if len(carts[i].Items) > 0 {
			if err := json.Unmarshal(carts[i].Items, &items); err != nil {
				return nil, errors.New("failed to parse cart items")
			}
		}
	}

	return carts, nil
}
