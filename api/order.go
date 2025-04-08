package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/internal/cart"
	"github.com/saanvi-iyer/gobblego-backend/internal/menu"
	"github.com/saanvi-iyer/gobblego-backend/internal/order"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB        *gorm.DB
	OrderRepo order.Repository
	CartRepo  cart.Repository
	MenuRepo  menu.Repository
}

func NewOrderHandler(
	db *gorm.DB,
	orderRepo order.Repository,
	cartRepo cart.Repository,
	menuRepo menu.Repository,
) *OrderHandler {
	return &OrderHandler{
		DB:        db,
		OrderRepo: orderRepo,
		CartRepo:  cartRepo,
		MenuRepo:  menuRepo,
	}
}

func (h *OrderHandler) PlaceOrder(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if !user.IsLeader {
		return c.Status(403).JSON(fiber.Map{"error": "Only the table leader can place orders"})
	}

	cartItems, err := h.CartRepo.GetCartItems(h.DB, user.CartID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch cart items"})
	}

	if len(cartItems) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Cart is empty"})
	}

	var totalAmount float64
	orderItems := make([]models.OrderItem, 0, len(cartItems))

	for _, item := range cartItems {
		menuItem, err := h.MenuRepo.GetMenuByID(h.DB, item.ItemID.String())
		if err != nil {
			continue
		}

		itemTotal := float64(item.Quantity) * menuItem.Price
		totalAmount += itemTotal

		orderItem := models.OrderItem{
			OrderItemID: uuid.New(),
			ItemID:      item.ItemID,
			Quantity:    item.Quantity,
			Price:       menuItem.Price,
			Notes:       item.Notes,
		}

		orderItems = append(orderItems, orderItem)
	}

	order := models.Order{
		OrderID:     uuid.New(),
		CartID:      user.CartID,
		UserID:      user.UserID,
		Status:      "pending",
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.OrderRepo.CreateOrder(h.DB, &order, orderItems); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
	}

	if err := h.CartRepo.ClearCartItems(h.DB, user.CartID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to clear cart"})
	}

	cart, err := h.CartRepo.GetCartByID(h.DB, user.CartID)
	if err == nil {
		cart.BillAmount = 0
		h.CartRepo.UpdateCart(h.DB, cart)
	}

	orderDetails, err := h.OrderRepo.GetOrderByID(h.DB, order.OrderID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch order details"})
	}

	orderItemsDetails, err := h.OrderRepo.GetOrderItems(h.DB, order.OrderID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch order items"})
	}

	return c.Status(201).JSON(fiber.Map{
		"order": orderDetails,
		"items": orderItemsDetails,
	})
}

func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	orders, err := h.OrderRepo.GetOrdersByCartID(h.DB, user.CartID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}

	return c.JSON(orders)
}

func (h *OrderHandler) GetOrderDetails(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	orderID := c.Params("order_id")
	oid, err := uuid.Parse(orderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	order, err := h.OrderRepo.GetOrderByID(h.DB, oid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	if order.CartID != user.CartID {
		return c.Status(403).JSON(fiber.Map{"error": "Order does not belong to your cart"})
	}

	orderItems, err := h.OrderRepo.GetOrderItems(h.DB, oid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch order items"})
	}

	type EnrichedOrderItem struct {
		models.OrderItem
		Item models.Menu `json:"item"`
	}

	enrichedItems := make([]EnrichedOrderItem, 0, len(orderItems))
	for _, item := range orderItems {
		menuItem, err := h.MenuRepo.GetMenuByID(h.DB, item.ItemID.String())
		if err != nil {
			continue
		}

		enriched := EnrichedOrderItem{
			OrderItem: item,
			Item:      *menuItem,
		}
		enrichedItems = append(enrichedItems, enriched)
	}

	return c.JSON(fiber.Map{
		"order": order,
		"items": enrichedItems,
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {

	_, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	orderID := c.Params("order_id")
	oid, err := uuid.Parse(orderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	type StatusUpdateRequest struct {
		Status string `json:"status"`
	}

	var req StatusUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"preparing": true,
		"ready":     true,
		"delivered": true,
	}

	if !validStatuses[req.Status] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status"})
	}

	if err := h.OrderRepo.UpdateOrderStatus(h.DB, oid, req.Status); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update order status"})
	}

	order, err := h.OrderRepo.GetOrderByID(h.DB, oid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch updated order"})
	}

	return c.JSON(order)
}
