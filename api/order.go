package api

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/internal/order"
	"github.com/saanvi-iyer/gobblego-backend/models"

	"gorm.io/gorm"
)

type OrderHandler struct {
	DB   *gorm.DB
	Repo order.Repository
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{
		DB:   db,
		Repo: order.NewOrderRepo(),
	}
}

type CreateOrderRequest struct {
	Items []models.CartItem `json:"items"`
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req struct {
		CartID string `json:"cart_id"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate cart ID format
	cartUUID, err := uuid.Parse(req.CartID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart ID format"})
	}

	// Fetch cart from DB
	var cart models.Cart
	if err := h.DB.First(&cart, "cart_id = ?", cartUUID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart not found"})
	}

	// Create order with cart items
	order := models.Order{
		OrderID:     uuid.New(),
		CartID:      cartUUID,
		Items:       cart.Items, // Copy cart items into order
		OrderStatus: "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save order to DB
	if err := h.Repo.CreateOrder(h.DB, &order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
	}

	// Clear cart items but keep cart_id, payment_status, and bill_amount intact
	emptyJSON, _ := json.Marshal([]models.CartItem{}) // Empty cart items
	cart.Items = emptyJSON

	if err := h.DB.Save(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to clear cart items"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"order_id": order.OrderID})
}

type UpdateOrderStatusRequest struct {
	OrderStatus string `json:"order_status"`
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	orderID := c.Params("order_id")

	var req UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.Repo.UpdateOrderStatus(h.DB, orderID, req.OrderStatus)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update order status"})
	}

	return c.JSON(fiber.Map{"message": "Order status updated"})
}

func (h *OrderHandler) ListAllOrders(c *fiber.Ctx) error {
	orders, err := h.Repo.ListAllOrders(h.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}

	return c.JSON(fiber.Map{"orders": orders})
}
