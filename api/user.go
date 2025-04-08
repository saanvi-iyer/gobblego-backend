package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"github.com/saanvi-iyer/gobblego-backend/utils"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB          *gorm.DB
	CartHandler *CartHandler
}

func NewUserHandler(db *gorm.DB, cartHandler *CartHandler) *UserHandler {
	return &UserHandler{
		DB:          db,
		CartHandler: cartHandler,
	}
}

func (h *UserHandler) JoinCart(c *fiber.Ctx) error {
	type JoinRequest struct {
		CartID   string `json:"cart_id"`
		UserName string `json:"user_name"`
	}

	var req JoinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	cid, err := uuid.Parse(req.CartID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid cart ID"})
	}

	var cart models.Cart
	if err := h.DB.First(&cart, "cart_id = ?", cid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			newCart := models.Cart{
				CartID:        cid,
				PaymentStatus: "pending",
				BillAmount:    0.0,
			}
			if err := h.DB.Create(&newCart).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to create cart"})
			}
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to check cart existence"})
		}
	}

	var userCount int64
	h.DB.Model(&models.User{}).Where("cart_id = ?", cid).Count(&userCount)

	user := models.User{
		UserID:   uuid.New(),
		CartID:   cid,
		UserName: req.UserName,
		IsLeader: userCount == 0,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to join Cart"})
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		SameSite: "lax",
		MaxAge:   3600 * 24, // 1 day
	}
	c.Cookie(&cookie)

	response := fiber.Map{
		"user_id":   user.UserID,
		"cart_id":   user.CartID,
		"user_name": user.UserName,
		"is_leader": user.IsLeader,
		"token":     token,
	}

	return c.Status(201).JSON(response)
}

func (h *UserHandler) GetCartUsers(c *fiber.Ctx) error {
	cartID := c.Params("cart_id")
	var users []models.User
	if err := h.DB.Where("cart_id = ?", cartID).Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	return c.JSON(users)
}
