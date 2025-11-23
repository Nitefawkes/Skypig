package handlers

import (
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nitefawkes/ham-radio-cloud/internal/models"
	"github.com/nitefawkes/ham-radio-cloud/internal/services"
	"github.com/nitefawkes/ham-radio-cloud/pkg/adif"
)

type ADIFHandler struct {
	qsoService *services.QSOService
}

func NewADIFHandler(qsoService *services.QSOService) *ADIFHandler {
	return &ADIFHandler{qsoService: qsoService}
}

func (h *ADIFHandler) Import(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		userID = "test-user-id"
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "missing_file",
			Message: "ADIF file is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Check file extension
	if !isADIFFile(file.Filename) {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_file",
			Message: "File must be .adi or .adif",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Open file
	fileReader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "file_error",
			Message: "Failed to open file",
			Code:    fiber.StatusInternalServerError,
		})
	}
	defer fileReader.Close()

	// Parse ADIF file
	parser := adif.NewParser(fileReader)
	qsos, err := parser.Parse()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "parse_error",
			Message: fmt.Sprintf("Failed to parse ADIF file: %v", err),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Set user ID for all QSOs
	for i := range qsos {
		qsos[i].UserID = userID.(string)
	}

	// Import QSOs (this will use bulk insert)
	imported := 0
	skipped := 0
	var importErrors []string

	for _, qso := range qsos {
		if err := h.qsoService.CreateQSO(&qso); err != nil {
			skipped++
			importErrors = append(importErrors, fmt.Sprintf("QSO with %s: %v", qso.Callsign, err))
		} else {
			imported++
		}
	}

	return c.JSON(fiber.Map{
		"message":  "ADIF import completed",
		"imported": imported,
		"skipped":  skipped,
		"total":    len(qsos),
		"errors":   importErrors,
	})
}

func (h *ADIFHandler) Export(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		userID = "test-user-id"
	}

	// Parse filters (same as List endpoint)
	filter := &models.QSOFilter{
		Callsign: c.Query("callsign"),
		Band:     c.Query("band"),
		Mode:     c.Query("mode"),
		Limit:    10000, // Export all matching records (up to 10k)
	}

	// Get QSOs
	qsos, err := h.qsoService.GetQSOs(userID.(string), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Generate ADIF file
	var buf bytes.Buffer
	exporter := adif.NewExporter(&buf)

	if err := exporter.Export(qsos, "Ham Radio Cloud", "1.0.0"); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "export_error",
			Message: fmt.Sprintf("Failed to generate ADIF file: %v", err),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Set headers for file download
	c.Set("Content-Type", "text/plain")
	c.Set("Content-Disposition", "attachment; filename=logbook.adi")

	return c.Send(buf.Bytes())
}

func isADIFFile(filename string) bool {
	return len(filename) > 4 && (
		filename[len(filename)-4:] == ".adi" ||
		filename[len(filename)-5:] == ".adif")
}
