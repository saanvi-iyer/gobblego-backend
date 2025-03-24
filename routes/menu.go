package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saanvi-iyer/gobblego-backend/api"
)

func MountMenuRoutes(app *fiber.App, handler *api.MenuHandler) {
	menu := app.Group("/api/v1/menu")
	menu.Get("/", handler.GetMenuItems)
	menu.Get("/:id", handler.GetMenuItemByID)
	menu.Post("/", handler.CreateMenuItem)
	menu.Put("/:id", handler.UpdateMenuItem)
	menu.Delete("/:id", handler.DeleteMenuItem)
}
