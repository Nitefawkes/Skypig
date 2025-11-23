package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router) {
	// Auth routes (placeholder for OAuth integration)
	auth := v1.Group("/auth")
	auth.Get("/login/qrz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "QRZ OAuth login endpoint - to be implemented",
		})
	})
	auth.Get("/callback/qrz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "QRZ OAuth callback endpoint - to be implemented",
		})
	})

	// QSO routes (placeholder)
	qso := v1.Group("/qso")
	qso.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "List QSOs - to be implemented",
			"data":    []interface{}{},
		})
	})
	qso.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Create QSO - to be implemented",
		})
	})

	// Propagation routes (placeholder)
	prop := v1.Group("/propagation")
	prop.Get("/current", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Current propagation data - to be implemented",
		})
	})

	// User routes (placeholder)
	user := v1.Group("/user")
	user.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "User profile - to be implemented",
		})
	})
}
