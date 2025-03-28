package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type TableHandler struct {
	DB *gorm.DB
}

func NewTableHandler(db *gorm.DB) *TableHandler {
	return &TableHandler{DB: db}
}

func (h *TableHandler) GetTable(c *fiber.Ctx) error {
	var table []models.Table
	if err := h.DB.Find(&table).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch table"})
	}
	return c.JSON(table)
}

func (h *TableHandler) GetTableByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var table models.Menu
	if err := h.DB.First(&table, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Table not found"})
	}
	return c.JSON(table)
}

func (h *TableHandler) CreateTable(c *fiber.Ctx) error {
	var table models.Table
	if err := c.BodyParser(&table); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	table.TableID=uuid.New()
	if err := h.DB.Create(&table).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create Table"})
	}
	return c.Status(201).JSON(table)
}