// main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/saanvi-iyer/gobblego-backend/api"
	"github.com/saanvi-iyer/gobblego-backend/config"
	"github.com/saanvi-iyer/gobblego-backend/internal/cart"
	"github.com/saanvi-iyer/gobblego-backend/internal/menu"
	"github.com/saanvi-iyer/gobblego-backend/internal/order"
	"github.com/saanvi-iyer/gobblego-backend/middleware"
	"github.com/saanvi-iyer/gobblego-backend/routes"
)

func main() {
	db := config.InitDB()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	// Initialize repositories
	menuRepo := menu.NewMenuRepo()
	cartRepo := cart.NewCartRepo()
	orderRepo := order.NewOrderRepo()

	// Initialize handlers
	menuHandler := api.NewMenuHandler(db.DB, menuRepo)
	cartHandler := api.NewCartHandler(db.DB, cartRepo, menuRepo)
	userHandler := api.NewUserHandler(db.DB, cartHandler)
	orderHandler := api.NewOrderHandler(db.DB, orderRepo, cartRepo, menuRepo)

	// Create auth middleware
	authMiddleware := middleware.Authenticate(db.DB)

	// Mount routes
	routes.MountMenuRoutes(app, menuHandler)
	routes.MountCartRoutes(app, cartHandler, authMiddleware)
	routes.MountUserRoutes(app, userHandler)
	routes.MountOrderRoutes(app, orderHandler, authMiddleware)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Pong", "status": 200})
	})

	log.Println("Starting server on :8080...")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
