package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Nitefawkes/Skypig/backend/internal/services"
)

// PropagationHandler handles propagation-related HTTP requests
type PropagationHandler struct {
	service *services.PropagationService
}

// NewPropagationHandler creates a new propagation handler
func NewPropagationHandler(service *services.PropagationService) *PropagationHandler {
	return &PropagationHandler{service: service}
}

// GetCurrent handles GET /api/v1/propagation
func (h *PropagationHandler) GetCurrent(c *fiber.Ctx) error {
	data, err := h.service.GetCurrent(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "FETCH_FAILED",
				"message": "Failed to fetch propagation data",
				"details": err.Error(),
			},
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"cache_ttl": 900, // 15 minutes in seconds
		},
	})
}

// GetForecast handles GET /api/v1/propagation/forecast
func (h *PropagationHandler) GetForecast(c *fiber.Ctx) error {
	forecast, err := h.service.GetForecast(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "FORECAST_FAILED",
				"message": "Failed to generate forecast",
				"details": err.Error(),
			},
		})
	}

	return c.JSON(fiber.Map{
		"data": forecast,
	})
}
