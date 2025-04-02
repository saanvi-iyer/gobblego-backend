package api

import (
	"encoding/json"

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
	CartID string `json:"cart_id"`
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	cartUUID, err := uuid.Parse(req.CartID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart ID format"})
	}

	cartIDsJSON, err := json.Marshal([]uuid.UUID{cartUUID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process cart IDs"})
	}

	order := models.Order{
		OrderID:     uuid.New(),
		CartIDs:     cartIDsJSON,
		OrderStatus: "pending",
	}

	err = h.Repo.CreateOrder(h.DB, &order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
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
