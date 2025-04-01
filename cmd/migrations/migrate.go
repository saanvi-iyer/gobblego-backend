package main

import (
	"log"
	"github.com/saanvi-iyer/gobblego-backend/config"
	"github.com/saanvi-iyer/gobblego-backend/models"
)

func main() {
	db := config.InitDB()

	err_menu := db.AutoMigrate(&models.Menu{})
	if err_menu != nil {
		log.Fatal("Failed to migrate Menu table:", err_menu)
	}

	err_cart := db.AutoMigrate(&models.Cart{})
	if err_cart != nil {
		log.Fatal("Failed to migrate Cart table:", err_cart)
	}

	err_user := db.AutoMigrate(&models.User{})
	if err_user != nil {
		log.Fatal("Failed to migrate User table:", err_user)
	}

	log.Println("Database migrations completed successfully")
}