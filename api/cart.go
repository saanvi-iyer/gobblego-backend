package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/internal/cart"
	"github.com/saanvi-iyer/gobblego-backend/internal/menu"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type CartHandler struct {
	DB       *gorm.DB
	CartRepo cart.Repository
	MenuRepo menu.Repository
}

func NewCartHandler(db *gorm.DB, cartRepo cart.Repository, menuRepo menu.Repository) *CartHandler {
	return &CartHandler{
		DB:       db,
		CartRepo: cartRepo,
		MenuRepo: menuRepo,
	}
}

func (h *CartHandler) AddItemToCart(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type AddItemRequest struct {
		ItemID   string `json:"item_id"`
		Quantity int    `json:"quantity"`
		Notes    string `json:"notes"`
	}

	var req AddItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if req.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Quantity must be positive"})
	}

	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	menuItem, err := h.MenuRepo.GetMenuByID(h.DB, req.ItemID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Item not found"})
	}

	if !menuItem.IsAvailable {
		return c.Status(400).JSON(fiber.Map{"error": "Item is not available"})
	}

	cartItems, err := h.CartRepo.GetCartItems(h.DB, user.CartID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check cart items"})
	}

	for _, item := range cartItems {
		if item.ItemID == itemID {

			item.Quantity += req.Quantity
			item.Notes = req.Notes
			item.UserID = user.UserID

			if err := h.CartRepo.UpdateCartItem(h.DB, &item); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to update cart item"})
			}

			if err := h.CartRepo.UpdateCartTotal(h.DB, user.CartID); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to update cart total"})
			}

			return c.JSON(item)
		}
	}

	cartItem := models.CartItem{
		CartItemID: uuid.New(),
		CartID:     user.CartID,
		ItemID:     itemID,
		UserID:     user.UserID,
		Quantity:   req.Quantity,
		Notes:      req.Notes,
		ItemPrice:  menuItem.Price,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.CartRepo.AddCartItem(h.DB, &cartItem); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add item to cart"})
	}

	if err := h.CartRepo.UpdateCartTotal(h.DB, user.CartID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update cart total"})
	}

	return c.Status(201).JSON(cartItem)
}

func (h *CartHandler) GetCartItems(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	cartItems, err := h.CartRepo.GetCartItems(h.DB, user.CartID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch cart items"})
	}

	type EnrichedCartItem struct {
		models.CartItem
		Item models.Menu `json:"item"`
	}

	enrichedItems := make([]EnrichedCartItem, 0, len(cartItems))
	for _, item := range cartItems {
		menuItem, err := h.MenuRepo.GetMenuByID(h.DB, item.ItemID.String())
		if err != nil {
			continue
		}

		enriched := EnrichedCartItem{
			CartItem: item,
			Item:     *menuItem,
		}
		enrichedItems = append(enrichedItems, enriched)
	}

	return c.JSON(enrichedItems)
}

func (h *CartHandler) RemoveItemFromCart(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	cartItemID := c.Params("cart_item_id")
	cid, err := uuid.Parse(cartItemID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid cart item ID"})
	}

	cartItem, err := h.CartRepo.GetCartItemByID(h.DB, cid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Cart item not found"})
	}

	if cartItem.CartID != user.CartID {
		return c.Status(403).JSON(fiber.Map{"error": "Cart item does not belong to your cart"})
	}

	if err := h.CartRepo.DeleteCartItem(h.DB, cid); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to remove item from cart"})
	}

	if err := h.CartRepo.UpdateCartTotal(h.DB, user.CartID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update cart total"})
	}

	return c.JSON(fiber.Map{"message": "Item removed from cart"})
}

func (h *CartHandler) UpdateCartItemQuantity(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	cartItemID := c.Params("cart_item_id")
	cid, err := uuid.Parse(cartItemID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid cart item ID"})
	}

	type UpdateRequest struct {
		Quantity int    `json:"quantity"`
		Notes    string `json:"notes"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if req.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Quantity must be positive"})
	}

	cartItem, err := h.CartRepo.GetCartItemByID(h.DB, cid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Cart item not found"})
	}

	if cartItem.CartID != user.CartID {
		return c.Status(403).JSON(fiber.Map{"error": "Cart item does not belong to your cart"})
	}

	cartItem.Quantity = req.Quantity
	cartItem.Notes = req.Notes
	cartItem.UserID = user.UserID
	cartItem.UpdatedAt = time.Now()

	if err := h.CartRepo.UpdateCartItem(h.DB, cartItem); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update cart item"})
	}

	if err := h.CartRepo.UpdateCartTotal(h.DB, user.CartID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update cart total"})
	}

	menuItem, err := h.MenuRepo.GetMenuByID(h.DB, cartItem.ItemID.String())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch menu item details"})
	}

	enrichedItem := struct {
		models.CartItem
		Item models.Menu `json:"item"`
	}{
		CartItem: *cartItem,
		Item:     *menuItem,
	}

	return c.JSON(enrichedItem)
}
