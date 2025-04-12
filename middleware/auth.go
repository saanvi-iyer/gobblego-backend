package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"github.com/saanvi-iyer/gobblego-backend/utils"
	"gorm.io/gorm"
)

func Authenticate(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		cookie := c.Cookies("jwt")
		if cookie == "" {

			authHeader := c.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				cookie = authHeader[7:]
			}
		}

		if cookie == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		claims, err := utils.VerifyToken(cookie)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		var user models.User
		if err := db.Where("user_id = ?", claims["user_id"]).First(&user).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "User not found"})
		}

		c.Locals("user", user)

		fmt.Print(user)
		fmt.Println("User authenticated")

		return c.Next()
	}
}
