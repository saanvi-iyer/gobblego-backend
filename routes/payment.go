package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saanvi-iyer/gobblego-backend/api"
)

func MountPaymentRoutes(app *fiber.App, handler *api.PaymentHandler, authMiddleware fiber.Handler) {

	payment := app.Group("/api/v1/payments")
	payment.Use(authMiddleware)

	payment.Post("/verify", handler.VerifyPayment)
	payment.Get("/:payment_id", handler.GetPaymentDetails)
	payment.Get("/cart/:cart_id", handler.GetCartPayments)
}
