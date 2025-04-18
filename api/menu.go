package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/internal/menu"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type MenuHandler struct {
	DB       *gorm.DB
	MenuRepo menu.Repository
}

func NewMenuHandler(db *gorm.DB, menuRepo menu.Repository) *MenuHandler {
	return &MenuHandler{
		DB:       db,
		MenuRepo: menuRepo,
	}
}

func (h *MenuHandler) GetMenuItems(c *fiber.Ctx) error {
	var menuItems []models.Menu
	if err := h.DB.Find(&menuItems).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch menu"})
	}
	return c.JSON(menuItems)
}

func (h *MenuHandler) GetMenuItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var menuItem models.Menu
	if err := h.DB.Where("item_id = ?", id).First(&menuItem).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Menu item not found"})
	}
	return c.JSON(menuItem)
}

func (h *MenuHandler) CreateMenuItem(c *fiber.Ctx) error {
	var menuItem models.Menu
	if err := c.BodyParser(&menuItem); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	menuItem.ItemID = uuid.New()
	if err := h.DB.Create(&menuItem).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create menu item"})
	}
	return c.Status(201).JSON(menuItem)
}

func (h *MenuHandler) UpdateMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var menuItem models.Menu
	if err := h.DB.Where("item_id = ?", id).First(&menuItem).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Menu item not found"})
	}

	if err := c.BodyParser(&menuItem); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	h.DB.Save(&menuItem)
	return c.JSON(menuItem)
}

func (h *MenuHandler) DeleteMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.DB.Where("item_id = ?", id).Delete(&models.Menu{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete menu item"})
	}
	return c.JSON(fiber.Map{"message": "Menu item deleted successfully"})
}
