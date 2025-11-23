package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nitefawkes/ham-radio-cloud/internal/models"
	"github.com/nitefawkes/ham-radio-cloud/internal/services"
)

type QSOHandler struct {
	service *services.QSOService
}

func NewQSOHandler(service *services.QSOService) *QSOHandler {
	return &QSOHandler{service: service}
}

type CreateQSORequest struct {
	Callsign       string    `json:"callsign"`
	Frequency      float64   `json:"frequency"`
	Band           string    `json:"band"`
	Mode           string    `json:"mode"`
	RST_Sent       string    `json:"rst_sent"`
	RST_Received   string    `json:"rst_received"`
	QSODate        string    `json:"qso_date"`
	TimeOn         string    `json:"time_on"`
	TimeOff        string    `json:"time_off,omitempty"`
	GridSquare     string    `json:"grid_square,omitempty"`
	Country        string    `json:"country,omitempty"`
	State          string    `json:"state,omitempty"`
	County         string    `json:"county,omitempty"`
	Comment        string    `json:"comment,omitempty"`
	TXPower        float64   `json:"tx_power,omitempty"`
}

func (h *QSOHandler) List(c *fiber.Ctx) error {
	// TODO: Get user ID from JWT token in middleware
	// For now, use a test user ID
	userID := c.Locals("userID")
	if userID == nil {
		userID = "test-user-id" // Development fallback
	}

	// Parse filters
	filter := &models.QSOFilter{
		Callsign:  c.Query("callsign"),
		Band:      c.Query("band"),
		Mode:      c.Query("mode"),
		Limit:     c.QueryInt("limit", 100),
		Offset:    c.QueryInt("offset", 0),
	}

	// Parse date filters
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}

	qsos, err := h.service.GetQSOs(userID.(string), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"data":   qsos,
		"total":  len(qsos),
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *QSOHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		userID = "test-user-id" // Development fallback
	}

	var req CreateQSORequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Parse timestamps
	var timeOn time.Time
	var timeOff *time.Time
	var qsoDate time.Time

	if req.TimeOn != "" {
		t, err := time.Parse(time.RFC3339, req.TimeOn)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "invalid_time",
				Message: "Invalid time_on format. Use RFC3339",
				Code:    fiber.StatusBadRequest,
			})
		}
		timeOn = t
	} else {
		timeOn = time.Now()
	}

	if req.TimeOff != "" {
		t, err := time.Parse(time.RFC3339, req.TimeOff)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "invalid_time",
				Message: "Invalid time_off format. Use RFC3339",
				Code:    fiber.StatusBadRequest,
			})
		}
		timeOff = &t
	}

	if req.QSODate != "" {
		t, err := time.Parse("2006-01-02", req.QSODate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "invalid_date",
				Message: "Invalid qso_date format. Use YYYY-MM-DD",
				Code:    fiber.StatusBadRequest,
			})
		}
		qsoDate = t
	} else {
		qsoDate = timeOn
	}

	qso := &models.QSO{
		UserID:       userID.(string),
		Callsign:     req.Callsign,
		Frequency:    req.Frequency,
		Band:         req.Band,
		Mode:         req.Mode,
		RST_Sent:     req.RST_Sent,
		RST_Received: req.RST_Received,
		QSODate:      qsoDate,
		TimeOn:       timeOn,
		TimeOff:      timeOff,
		GridSquare:   req.GridSquare,
		Country:      req.Country,
		State:        req.State,
		County:       req.County,
		Comment:      req.Comment,
		TXPower:      req.TXPower,
	}

	if err := h.service.CreateQSO(qso); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "QSO created successfully",
		"id":      qso.ID,
		"qso":     qso,
	})
}

func (h *QSOHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		userID = "test-user-id"
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "missing_id",
			Message: "QSO ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	var req CreateQSORequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Parse timestamps
	timeOn, _ := time.Parse(time.RFC3339, req.TimeOn)
	var timeOff *time.Time
	if req.TimeOff != "" {
		t, _ := time.Parse(time.RFC3339, req.TimeOff)
		timeOff = &t
	}
	qsoDate, _ := time.Parse("2006-01-02", req.QSODate)

	qso := &models.QSO{
		ID:           id,
		UserID:       userID.(string),
		Callsign:     req.Callsign,
		Frequency:    req.Frequency,
		Band:         req.Band,
		Mode:         req.Mode,
		RST_Sent:     req.RST_Sent,
		RST_Received: req.RST_Received,
		QSODate:      qsoDate,
		TimeOn:       timeOn,
		TimeOff:      timeOff,
		GridSquare:   req.GridSquare,
		Country:      req.Country,
		State:        req.State,
		County:       req.County,
		Comment:      req.Comment,
		TXPower:      req.TXPower,
	}

	if err := h.service.UpdateQSO(qso); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "update_error",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	return c.JSON(fiber.Map{
		"message": "QSO updated successfully",
		"qso":     qso,
	})
}

func (h *QSOHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		userID = "test-user-id"
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "missing_id",
			Message: "QSO ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	if err := h.service.DeleteQSO(id, userID.(string)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
			Code:    fiber.StatusNotFound,
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *QSOHandler) GetStats(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		userID = "test-user-id"
	}

	count, err := h.service.GetQSOCount(userID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"total_qsos": count,
		"user_id":    userID,
	})
}
