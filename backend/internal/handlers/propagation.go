package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nitefawkes/ham-radio-cloud/internal/services"
)

type PropagationHandler struct {
	service *services.PropagationService
}

func NewPropagationHandler(service *services.PropagationService) *PropagationHandler {
	return &PropagationHandler{service: service}
}

func (h *PropagationHandler) GetCurrent(c *fiber.Ctx) error {
	data, err := h.service.GetCurrent()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "fetch_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.JSON(data)
}

func (h *PropagationHandler) GetBandConditions(c *fiber.Ctx) error {
	// Get current propagation data
	data, err := h.service.GetCurrent()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "fetch_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Calculate band conditions
	conditions := h.service.GetBandConditions(data)

	return c.JSON(fiber.Map{
		"timestamp":  data.Timestamp,
		"conditions": conditions,
		"source":     data.Source,
	})
}

func (h *PropagationHandler) RefreshData(c *fiber.Ctx) error {
	if err := h.service.FetchAndStore(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "refresh_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Propagation data refreshed successfully",
	})
}
