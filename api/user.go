package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) JoinTable(c *fiber.Ctx) error {
	type JoinRequest struct {
		TableID  string `json:"table_id"`
		UserName string `json:"user_name"`
	}
	
	var req JoinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	tid, err := uuid.Parse(req.TableID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid table ID"})
	}

	var userCount int64
	h.DB.Model(&models.User{}).Where("table_id = ?", tid).Count(&userCount)

	user := models.User{
		UserID:   uuid.New(),
		TableID:  tid,
		UserName: req.UserName,
		IsLeader: userCount == 0,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to join table"})
	}

	return c.Status(201).JSON(user)
}

func (h *UserHandler) GetTableUsers(c *fiber.Ctx) error {
	tableID := c.Params("table_id")
	var users []models.User
	if err := h.DB.Where("table_id = ?", tableID).Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	return c.JSON(users)
}