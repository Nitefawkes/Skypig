package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "ham-radio-cloud-api",
		Version:   "1.0.0",
	})
}
