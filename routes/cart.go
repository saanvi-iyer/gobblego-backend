package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saanvi-iyer/gobblego-backend/api"
)

func MountCartRoutes(app *fiber.App, handler *api.CartHandler) {
	cart := app.Group("/api/v1/cart")
	cart.Get("/", handler.GetAllCarts)
	cart.Post("/", handler.AddToCart)
	cart.Patch("/", handler.UpdateCartItem)
}
