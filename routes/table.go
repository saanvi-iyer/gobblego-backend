package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saanvi-iyer/gobblego-backend/api"
)

func MountTableRoutes(app *fiber.App, handler *api.TableHandler) {
	menu := app.Group("/api/v1/table")
	menu.Get("/", handler.GetTable)
	menu.Get("/:id", handler.GetTableByID)
	menu.Post("/", handler.CreateTable)
}
