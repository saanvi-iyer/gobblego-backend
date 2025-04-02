package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saanvi-iyer/gobblego-backend/api"
)

func MountOrderRoutes(app *fiber.App, handler *api.OrderHandler) {
	orderGroup := app.Group("/api/v1/order")
	orderGroup.Post("/", handler.CreateOrder)
	orderGroup.Patch("/:order_id", handler.UpdateOrderStatus)
	orderGroup.Get("/", handler.ListAllOrders)

}
