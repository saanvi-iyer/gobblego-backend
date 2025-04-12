package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saanvi-iyer/gobblego-backend/api"
)

func MountCartRoutes(app *fiber.App, handler *api.CartHandler, authMiddleware fiber.Handler) {
	cart := app.Group("/api/v1/cart")

	cart.Post("/", handler.CreateCart)

	cart.Use(authMiddleware)
	cart.Post("/items", handler.AddItemToCart)
	cart.Get("/items", handler.GetCartItems)
	cart.Delete("/items/:cart_item_id", handler.RemoveItemFromCart)
	cart.Put("/items/:cart_item_id", handler.UpdateCartItemQuantity)
}
