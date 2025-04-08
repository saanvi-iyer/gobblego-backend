package main

import (
	"log"

	"github.com/saanvi-iyer/gobblego-backend/config"
	"github.com/saanvi-iyer/gobblego-backend/models"
)

func main() {
	db := config.InitDB()

	if err := db.AutoMigrate(&models.Menu{}); err != nil {
		log.Fatal("Failed to migrate Menu table:", err)
	}
	log.Println("Menu table migrated successfully")

	if err := db.AutoMigrate(&models.Cart{}); err != nil {
		log.Fatal("Failed to migrate Cart table:", err)
	}
	log.Println("Cart table migrated successfully")

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate User table:", err)
	}
	log.Println("User table migrated successfully")

	if err := db.AutoMigrate(&models.CartItem{}); err != nil {
		log.Fatal("Failed to migrate CartItem table:", err)
	}
	log.Println("CartItem table migrated successfully")

	if err := db.AutoMigrate(&models.Order{}); err != nil {
		log.Fatal("Failed to migrate Order table:", err)
	}
	log.Println("Order table migrated successfully")

	if err := db.AutoMigrate(&models.OrderItem{}); err != nil {
		log.Fatal("Failed to migrate OrderItem table:", err)
	}
	log.Println("OrderItem table migrated successfully")

	log.Println("All database migrations completed successfully")
}
