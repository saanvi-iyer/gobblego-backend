package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/saanvi-iyer/gobblego-backend/api"
	"github.com/saanvi-iyer/gobblego-backend/config"
	"github.com/saanvi-iyer/gobblego-backend/routes"
)

func main() {
	db := config.InitDB()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	MenuHandler := api.NewMenuHandler(db.DB)
	routes.MountMenuRoutes(app, MenuHandler)

	CartHandler := api.NewCartHandler(db.DB)
	routes.MountCartRoutes(app, CartHandler)

	UserHandler := api.NewUserHandler(db.DB, CartHandler)
	routes.MountUserRoutes(app, UserHandler)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Pong", "status": 200})
	})

	log.Println("Starting server on :8080...")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}