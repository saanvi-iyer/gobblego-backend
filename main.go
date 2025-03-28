package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/saanvi-iyer/gobblego-backend/api"
	"github.com/saanvi-iyer/gobblego-backend/config"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"github.com/saanvi-iyer/gobblego-backend/routes"
	"log"
)

func main() {
	db := config.InitDB()
	err := db.AutoMigrate(&models.Menu{})
	if err != nil {
		log.Fatal("Failed to migrate Menu table:", err)
	}
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	MenuHandler := api.NewMenuHandler(db.DB)
	routes.MountMenuRoutes(app, MenuHandler)

	TableHandler := api.NewTableHandler(db.DB)
	routes.MountTableRoutes(app, TableHandler)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Pong", "status": 200})
	})

	log.Println("Starting server on :8080...")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
