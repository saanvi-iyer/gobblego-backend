package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saanvi-iyer/gobblego-backend/api"
)

func MountUserRoutes(app *fiber.App, handler *api.UserHandler) {
	user := app.Group("/api/v1/users")
	user.Post("/", handler.JoinTable)
	user.Get("/:table_id", handler.GetTableUsers)
}