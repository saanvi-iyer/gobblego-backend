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
		if len(cart.Items) > 0 {
			if err := json.Unmarshal(cart.Items, &items); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to parse cart items"})
			}
		}

		var enrichedItems []fiber.Map
		for _, item := range items {
			var menuItem models.Menu
			if err := h.DB.Where("item_id = ?", item.ItemID).First(&menuItem).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch item details"})
			}

			var user models.User
			if err := h.DB.Where("user_id = ?", item.UserID).First(&user).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user details"})
			}

			enrichedItems = append(enrichedItems, fiber.Map{
				"item_id":   item.ItemID,
				"item_name": menuItem.ItemName,
				"image":     menuItem.Images,
				"price":     menuItem.Price,
				"quantity":  item.Quantity,
				"user_id":   item.UserID,
				"user_name": user.UserName,
			})
		}

		if newItems, err := json.Marshal(enrichedItems); err != nil {
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

	var cart models.Cart
	if err := h.DB.Where("cart_id = ?", user.CartID).First(&cart).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch cart"})
	}

	var items []models.CartItem
	if len(cart.Items) > 0 {
		if err := json.Unmarshal(cart.Items, &items); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse cart items"})
		}
	}

	itemExists := false
	for i := range items {
		if items[i].ItemID == req.ItemID && items[i].UserID == req.UserID {

			items[i].Quantity += req.Quantity
			itemExists = true
			break
		}
	}

	if !itemExists {
		newItem := models.CartItem{
			ItemID:   req.ItemID,
			UserID:   req.UserID,
			Quantity: req.Quantity,
		}
		items = append(items, newItem)
	}

	updatedItems, err := json.Marshal(items)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update cart items"})
	}

	cart.Items = updatedItems
	if err := h.DB.Save(&cart).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save updated cart"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Item added/updated in cart"})
}

func (h *CartHandler) UpdateCartItem(c *fiber.Ctx) error {

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

	var cart models.Cart
	if err := h.DB.Where("cart_id = ?", user.CartID).First(&cart).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch cart"})
	}

	var items []models.CartItem
	if len(cart.Items) > 0 {
		if err := json.Unmarshal(cart.Items, &items); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse cart items"})
		}
	}

	itemFound := false
	for i := range items {
		if items[i].ItemID == req.ItemID && items[i].UserID == req.UserID {
			itemFound = true
			if req.Quantity > 0 {

				items[i].Quantity = req.Quantity
			} else {

				items = append(items[:i], items[i+1:]...)
			}
			break
		}
	}

	if !itemFound {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found in cart"})
	}

	updatedItems, err := json.Marshal(items)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update cart items"})
	}

	cart.Items = updatedItems
	if err := h.DB.Save(&cart).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save updated cart"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Cart item updated"})
}
