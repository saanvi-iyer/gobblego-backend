package api

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/internal/cart"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type CartHandler struct {
	Repo cart.Repository
	DB   *gorm.DB
}

type CartResponse struct {
	ItemID   uuid.UUID `json:"item_id"`
	ItemName string    `json:"item_name"`
	Image    string    `json:"image"`
	Price    float64   `json:"price"`
	Quantity int       `json:"quantity"`
	UserName string    `json:"user_name"`
	UserID   uuid.UUID `json:"user_id"`
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{
		DB:   db,
		Repo: cart.NewCartRepo(),
	}
}

func (h *CartHandler) GetAllCarts(c *fiber.Ctx) error {
	var carts []models.Cart
	if err := h.DB.Find(&carts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch carts"})
	}

	for i, cart := range carts {
		var items []models.CartItem
		if err := json.Unmarshal(cart.Items, &items); err != nil {

			var item models.CartItem
			if err := json.Unmarshal(cart.Items, &item); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to parse cart items"})
			}
			items = append(items, item)
		}
		if newItems, err := json.Marshal(items); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to marshal cart items"})
		} else {
			carts[i].Items = newItems
		}
	}

	return c.JSON(carts)
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {

	var req struct {
		ItemID   uuid.UUID `json:"item_id"`
		Quantity int       `json:"quantity"`
		UserID   uuid.UUID `json:"user_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := h.DB.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user"})
	}

	if user.CartID == uuid.Nil {

		newCart := models.Cart{
			PaymentStatus: "pending",
			BillAmount:    0.0,
			Items:         []byte("[]"),
		}

		if err := h.DB.Create(&newCart).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create new cart"})
		}

		user.CartID = newCart.CartID
		if err := h.DB.Save(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update user with cart_id"})
		}
	}

	item := models.CartItem{
		ItemID:   req.ItemID,
		UserID:   req.UserID,
		Quantity: req.Quantity,
	}

	if err := h.Repo.AddToCart(h.DB, user.CartID, &item); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add item to cart: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Item added to cart"})
}
