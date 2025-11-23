package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/Nitefawkes/Skypig/backend/internal/models"
	"github.com/Nitefawkes/Skypig/backend/internal/services"
)

// ADIFHandler handles ADIF import/export HTTP requests
type ADIFHandler struct {
	service *services.ADIFService
}

// NewADIFHandler creates a new ADIF handler
func NewADIFHandler(service *services.ADIFService) *ADIFHandler {
	return &ADIFHandler{service: service}
}

// ImportADIF handles POST /api/v1/qsos/import
func (h *ADIFHandler) ImportADIF(c *fiber.Ctx) error {
	// TODO: Get userID from JWT token
	userID := int64(1)

	// Get ADIF content from request body
	content := string(c.Body())
	if content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "EMPTY_CONTENT",
				"message": "Request body is empty",
			},
		})
	}

	// Check for strict mode query parameter
	strict := c.Query("strict", "false") == "true"

	// Import ADIF
	result, err := h.service.ImportADIF(c.Context(), userID, content, strict)
	if err != nil && strict {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "IMPORT_FAILED",
				"message": err.Error(),
			},
			"data": result,
		})
	}

	// Return result
	status := fiber.StatusOK
	if result.FailedRecords > 0 || result.SkippedRecords > 0 {
		status = fiber.StatusPartialContent
	}

	return c.Status(status).JSON(fiber.Map{
		"data": result,
		"meta": fiber.Map{
			"message": fmt.Sprintf("Imported %d of %d records", result.ImportedRecords, result.TotalRecords),
		},
	})
}

// ExportADIF handles GET /api/v1/qsos/export
func (h *ADIFHandler) ExportADIF(c *fiber.Ctx) error {
	// TODO: Get userID from JWT token
	userID := int64(1)

	// Parse filter parameters (same as ListQSOs)
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

	// Parse limit (default to all if not specified)
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	// Export ADIF
	adifContent, err := h.service.ExportADIF(c.Context(), userID, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "EXPORT_FAILED",
				"message": err.Error(),
			},
		})
	}

	// Set headers for file download
	filename := fmt.Sprintf("hamradio_cloud_export_%s.adi", time.Now().Format("20060102_150405"))
	c.Set("Content-Type", "text/plain; charset=utf-8")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return c.SendString(adifContent)
}

// ValidateADIF handles POST /api/v1/qsos/validate
func (h *ADIFHandler) ValidateADIF(c *fiber.Ctx) error {
	// Get ADIF content from request body
	content := string(c.Body())
	if content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "EMPTY_CONTENT",
				"message": "Request body is empty",
			},
		})
	}

	// Validate ADIF
	result, err := h.service.ValidateADIF(content)

	status := fiber.StatusOK
	if err != nil {
		status = fiber.StatusBadRequest
	}

	return c.Status(status).JSON(fiber.Map{
		"data": result,
		"meta": fiber.Map{
			"valid":   err == nil,
			"message": func() string {
				if err == nil {
					return fmt.Sprintf("Valid ADIF with %d records", result.TotalRecords)
				}
				return fmt.Sprintf("Invalid ADIF: %d errors found", result.FailedRecords)
			}(),
		},
	})
}
