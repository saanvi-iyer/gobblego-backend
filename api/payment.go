package api

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	razorpay "github.com/razorpay/razorpay-go"
	utils "github.com/razorpay/razorpay-go/utils"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	DB        *gorm.DB
	RzpClient *razorpay.Client
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	rzpKey := os.Getenv("RAZORPAY_KEY_ID")
	rzpSecret := os.Getenv("RAZORPAY_KEY_SECRET")

	client := razorpay.NewClient(rzpKey, rzpSecret)

	return &PaymentHandler{
		DB:        db,
		RzpClient: client,
	}
}

func (h *PaymentHandler) VerifyPayment(c *fiber.Ctx) error {
	type VerifyRequest struct {
		PaymentID         string `json:"payment_id"`
		RazorpayPaymentID string `json:"razorpay_payment_id"`
		RazorpayOrderID   string `json:"razorpay_order_id"`
		RazorpaySignature string `json:"razorpay_signature"`
	}

	var req VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payment ID"})
	}

	var payment models.Payment
	if err := h.DB.First(&payment, "payment_id = ?", paymentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Payment not found"})
	}

	params := map[string]interface{}{
		"razorpay_order_id":   req.RazorpayOrderID,
		"razorpay_payment_id": req.RazorpayPaymentID,
	}

	signature := req.RazorpaySignature
	secret := os.Getenv("RAZORPAY_KEY_SECRET")

	isValid := utils.VerifyPaymentSignature(params, signature, secret)

	if !isValid {

		h.DB.Model(&payment).Updates(map[string]interface{}{
			"status":              "failed",
			"razorpay_payment_id": req.RazorpayPaymentID,
			"updated_at":          time.Now(),
		})

		return c.Status(400).JSON(fiber.Map{"error": "Payment verification failed"})
	}

	if err := h.DB.Model(&payment).Updates(map[string]interface{}{
		"status":              "successful",
		"razorpay_payment_id": req.RazorpayPaymentID,
		"updated_at":          time.Now(),
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update payment status"})
	}

	if err := h.DB.Model(&models.Cart{}).Where("cart_id = ?", payment.CartID).Update("payment_status", "completed").Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update cart status"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Payment verified successfully",
		"payment": payment,
	})
}

func (h *PaymentHandler) GetPaymentDetails(c *fiber.Ctx) error {
	paymentID := c.Params("payment_id")

	var payment models.Payment
	if err := h.DB.First(&payment, "payment_id = ?", paymentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Payment not found"})
	}

	return c.JSON(payment)
}

func (h *PaymentHandler) GetCartPayments(c *fiber.Ctx) error {
	cartID := c.Params("cart_id")

	var payments []models.Payment
	if err := h.DB.Where("cart_id = ?", cartID).Find(&payments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch payments"})
	}

	return c.JSON(payments)
}
