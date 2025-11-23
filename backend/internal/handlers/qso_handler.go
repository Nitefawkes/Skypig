package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/Nitefawkes/Skypig/backend/internal/models"
	"github.com/Nitefawkes/Skypig/backend/internal/services"
)

// QSOHandler handles QSO-related HTTP requests
type QSOHandler struct {
	service *services.QSOService
}

// NewQSOHandler creates a new QSO handler
func NewQSOHandler(service *services.QSOService) *QSOHandler {
	return &QSOHandler{service: service}
}

// CreateQSO handles POST /api/v1/qsos
func (h *QSOHandler) CreateQSO(c *fiber.Ctx) error {
	// TODO: Get userID from JWT token
	// For now, using a hardcoded test user ID
	userID := int64(1)

	var req models.QSO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
				"details": err.Error(),
			},
		})
	}

	qso, err := h.service.CreateQSO(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "CREATE_FAILED",
				"message": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": qso,
	})
}

// GetQSO handles GET /api/v1/qsos/:id
func (h *QSOHandler) GetQSO(c *fiber.Ctx) error {
	// TODO: Get userID from JWT token
	userID := int64(1)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_ID",
				"message": "Invalid QSO ID",
			},
		})
	}

	qso, err := h.service.GetQSO(c.Context(), id, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": err.Error(),
			},
		})
	}

	return c.JSON(fiber.Map{
		"data": qso,
	})
}

// ListQSOs handles GET /api/v1/qsos
func (h *QSOHandler) ListQSOs(c *fiber.Ctx) error {
	// TODO: Get userID from JWT token
	userID := int64(1)

	// Parse query parameters
	filter := models.QSOFilter{
		Callsign: c.Query("callsign"),
		Band:     c.Query("band"),
		Mode:     c.Query("mode"),
		Country:  c.Query("country"),
	}

	// Parse date range
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filter.DateFrom = t
		}
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			filter.DateTo = t
		}
	}

	// Parse pagination
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	qsos, total, err := h.service.ListQSOs(c.Context(), userID, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "LIST_FAILED",
				"message": err.Error(),
			},
		})
	}

	return c.JSON(fiber.Map{
		"data": qsos,
		"meta": fiber.Map{
			"total":    total,
			"limit":    filter.Limit,
			"offset":   filter.Offset,
			"returned": len(qsos),
		},
	})
}

// UpdateQSO handles PUT /api/v1/qsos/:id
func (h *QSOHandler) UpdateQSO(c *fiber.Ctx) error {
	// TODO: Get userID from JWT token
	userID := int64(1)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_ID",
				"message": "Invalid QSO ID",
			},
		})
	}

	var req models.QSO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
				"details": err.Error(),
			},
		})
	}

	req.ID = id
	qso, err := h.service.UpdateQSO(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "UPDATE_FAILED",
				"message": err.Error(),
			},
		})
	}

	return c.JSON(fiber.Map{
		"data": qso,
	})
}

// DeleteQSO handles DELETE /api/v1/qsos/:id
func (h *QSOHandler) DeleteQSO(c *fiber.Ctx) error {
	// TODO: Get userID from JWT token
	userID := int64(1)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "INVALID_ID",
				"message": "Invalid QSO ID",
			},
		})
	}

	if err := h.service.DeleteQSO(c.Context(), id, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "DELETE_FAILED",
				"message": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// GetStats handles GET /api/v1/qsos/stats
func (h *QSOHandler) GetStats(c *fiber.Ctx) error {
	// TODO: Get userID from JWT token
	userID := int64(1)

	stats, err := h.service.GetQSOStats(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "STATS_FAILED",
				"message": err.Error(),
			},
		})
	}

	return c.JSON(fiber.Map{
		"data": stats,
	})
}
