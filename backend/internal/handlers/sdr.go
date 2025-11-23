package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nitefawkes/ham-radio-cloud/internal/models"
	"github.com/nitefawkes/ham-radio-cloud/internal/services"
)

type SDRHandler struct {
	service *services.SDRService
}

func NewSDRHandler(service *services.SDRService) *SDRHandler {
	return &SDRHandler{service: service}
}

// List returns a paginated list of SDR receivers
func (h *SDRHandler) List(c *fiber.Ctx) error {
	filter := &models.SDRFilter{
		Type:    c.Query("type"),
		Country: c.Query("country"),
		Status:  c.Query("status", "online"), // Default to online only
		Search:  c.Query("search"),
	}

	// Parse bands filter
	if bands := c.Query("bands"); bands != "" {
		filter.Bands = []string{bands}
	}

	// Parse pagination
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	} else {
		filter.Limit = 50 // Default limit
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}

	sdrs, total, err := h.service.List(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "fetch_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"sdrs":   sdrs,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetByID returns a single SDR by ID
func (h *SDRHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "SDR ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	sdr, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "not_found",
			Message: "SDR not found",
			Code:    fiber.StatusNotFound,
		})
	}

	return c.JSON(sdr)
}

// Search searches for SDRs by name, location, or callsign
func (h *SDRHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_query",
			Message: "Search query is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	sdrs, err := h.service.Search(query, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "search_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"sdrs":  sdrs,
		"query": query,
		"count": len(sdrs),
	})
}

// RefreshDirectory manually triggers a refresh of the SDR directory
func (h *SDRHandler) RefreshDirectory(c *fiber.Ctx) error {
	if err := h.service.RefreshDirectory(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "refresh_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"message": "SDR directory refreshed successfully",
	})
}

// GetStats returns statistics about the SDR directory
func (h *SDRHandler) GetStats(c *fiber.Ctx) error {
	// Get total count
	total, err := h.service.List(&models.SDRFilter{Limit: 1})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "stats_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Get online count
	onlineSDRs, onlineCount, err := h.service.List(&models.SDRFilter{
		Status: "online",
		Limit:  1,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "stats_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Get count by type
	kiwiSDRs, kiwiCount, _ := h.service.List(&models.SDRFilter{Type: "kiwisdr", Limit: 1})
	webSDRs, webCount, _ := h.service.List(&models.SDRFilter{Type: "websdr", Limit: 1})
	openWebRXs, openWebRXCount, _ := h.service.List(&models.SDRFilter{Type: "openwebrx", Limit: 1})

	return c.JSON(fiber.Map{
		"total":  len(total),
		"online": onlineCount,
		"by_type": fiber.Map{
			"kiwisdr":   kiwiCount,
			"websdr":    webCount,
			"openwebrx": openWebRXCount,
		},
		"sample_sdrs": onlineSDRs,
	})
}
