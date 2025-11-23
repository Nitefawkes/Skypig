package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router, qsoHandler *QSOHandler, adifHandler *ADIFHandler, propHandler *PropagationHandler) {
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

	// QSO routes (fully implemented)
	qso := v1.Group("/qso")
	qso.Get("/", qsoHandler.List)
	qso.Post("/", qsoHandler.Create)
	qso.Put("/:id", qsoHandler.Update)
	qso.Delete("/:id", qsoHandler.Delete)
	qso.Get("/stats", qsoHandler.GetStats)
	qso.Post("/import/adif", adifHandler.Import)
	qso.Get("/export/adif", adifHandler.Export)

	// Propagation routes (fully implemented)
	prop := v1.Group("/propagation")
	prop.Get("/current", propHandler.GetCurrent)
	prop.Get("/bands", propHandler.GetBandConditions)
	prop.Post("/refresh", propHandler.RefreshData)

	// User routes (placeholder)
	user := v1.Group("/user")
	user.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "User profile - to be implemented",
		})
	})
}
