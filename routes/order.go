package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saanvi-iyer/gobblego-backend/api"
)

func MountOrderRoutes(app *fiber.App, handler *api.OrderHandler, authMiddleware fiber.Handler) {
	orders := app.Group("/api/v1/orders")
	orders.Get("/all", handler.GetAllOrders)
	orders.Put("/:order_id/status", handler.UpdateOrderStatus)

	orders.Use(authMiddleware)

	orders.Post("/", handler.PlaceOrder)
	orders.Get("/", handler.GetOrders)
	orders.Get("/:order_id", handler.GetOrderDetails)
	orders.Post("/checkout", handler.Checkout)
}
