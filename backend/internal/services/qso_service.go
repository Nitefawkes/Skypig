package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Nitefawkes/Skypig/backend/internal/database"
	"github.com/Nitefawkes/Skypig/backend/internal/models"
)

// QSOService handles QSO business logic
type QSOService struct {
	qsoRepo  *database.QSORepository
	userRepo *database.UserRepository
}

// NewQSOService creates a new QSO service
func NewQSOService(qsoRepo *database.QSORepository, userRepo *database.UserRepository) *QSOService {
	return &QSOService{
		qsoRepo:  qsoRepo,
		userRepo: userRepo,
	}
}

// CreateQSO creates a new QSO after validation
func (s *QSOService) CreateQSO(ctx context.Context, userID int64, qso *models.QSO) (*models.QSO, error) {
	// Get user to check QSO limits
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user can create more QSOs
	if !user.CanCreateQSO() {
		return nil, fmt.Errorf("QSO limit reached for %s tier (%d/%d)", user.Tier, user.QSOCount, user.QSOLimit)
	}

	// Validate QSO
	if err := s.validateQSO(qso); err != nil {
		return nil, err
	}

	// Set user ID
	qso.UserID = userID

	// Set defaults
	if qso.QSODate.IsZero() {
		qso.QSODate = qso.TimeOn
	}

	// Create QSO
	if err := s.qsoRepo.Create(ctx, qso); err != nil {
		return nil, err
	}

	return qso, nil
}

// GetQSO retrieves a QSO by ID
func (s *QSOService) GetQSO(ctx context.Context, id, userID int64) (*models.QSO, error) {
	return s.qsoRepo.GetByID(ctx, id, userID)
}

// ListQSOs retrieves QSOs with filtering
func (s *QSOService) ListQSOs(ctx context.Context, userID int64, filter models.QSOFilter) ([]*models.QSO, int, error) {
	return s.qsoRepo.List(ctx, userID, filter)
}

// UpdateQSO updates an existing QSO
func (s *QSOService) UpdateQSO(ctx context.Context, userID int64, qso *models.QSO) (*models.QSO, error) {
	// Verify QSO exists and belongs to user
	existing, err := s.qsoRepo.GetByID(ctx, qso.ID, userID)
	if err != nil {
		return nil, err
	}

	// Validate updated QSO
	if err := s.validateQSO(qso); err != nil {
		return nil, err
	}

	// Ensure user ID hasn't changed
	qso.UserID = existing.UserID

	// Update QSO
	if err := s.qsoRepo.Update(ctx, qso); err != nil {
		return nil, err
	}

	return qso, nil
}

// DeleteQSO removes a QSO
func (s *QSOService) DeleteQSO(ctx context.Context, id, userID int64) error {
	return s.qsoRepo.Delete(ctx, id, userID)
}

// GetQSOStats returns statistics for a user's QSOs
func (s *QSOService) GetQSOStats(ctx context.Context, userID int64) (*models.QSOStats, error) {
	// For now, return basic stats - can be expanded later
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats := &models.QSOStats{
		TotalQSOs:     user.QSOCount,
		QSOLimit:      user.QSOLimit,
		RemainingQSOs: user.QSOLimit - user.QSOCount,
	}

	if user.Tier == "contester" {
		stats.RemainingQSOs = -1 // Unlimited
	}

	return stats, nil
}

// validateQSO performs validation on QSO data
func (s *QSOService) validateQSO(qso *models.QSO) error {
	// Required fields
	if qso.Callsign == "" {
		return fmt.Errorf("callsign is required")
	}

	if qso.TimeOn.IsZero() {
		return fmt.Errorf("time_on is required")
	}

	// Validate callsign format (basic)
	qso.Callsign = strings.ToUpper(strings.TrimSpace(qso.Callsign))
	if len(qso.Callsign) < 3 || len(qso.Callsign) > 20 {
		return fmt.Errorf("callsign must be 3-20 characters")
	}

	// Validate band if provided
	if qso.Band != "" {
		if !isValidBand(qso.Band) {
			return fmt.Errorf("invalid band: %s", qso.Band)
		}
	}

	// Validate mode if provided
	if qso.Mode != "" {
		qso.Mode = strings.ToUpper(qso.Mode)
	}

	// Validate grid square if provided
	if qso.GridSquare != "" {
		qso.GridSquare = strings.ToUpper(qso.GridSquare)
		if len(qso.GridSquare) < 4 || len(qso.GridSquare) > 8 {
			return fmt.Errorf("invalid grid square format")
		}
	}

	// Validate time_off is after time_on
	if !qso.TimeOff.IsZero() && qso.TimeOff.Before(qso.TimeOn) {
		return fmt.Errorf("time_off must be after time_on")
	}

	// Validate frequency ranges
	if qso.Freq < 0 || qso.Freq > 300000 { // 0-300 GHz in MHz
		return fmt.Errorf("invalid frequency")
	}

	// Validate power
	if qso.TXPower < 0 || qso.TXPower > 10000 { // Max 10kW
		return fmt.Errorf("invalid transmit power")
	}

	return nil
}

// isValidBand checks if the band is valid
func isValidBand(band string) bool {
	validBands := map[string]bool{
		"2190m": true, "630m": true, "560m": true, "160m": true,
		"80m": true, "60m": true, "40m": true, "30m": true,
		"20m": true, "17m": true, "15m": true, "12m": true,
		"10m": true, "6m": true, "4m": true, "2m": true,
		"1.25m": true, "70cm": true, "33cm": true, "23cm": true,
		"13cm": true, "9cm": true, "6cm": true, "3cm": true,
		"1.25cm": true, "6mm": true, "4mm": true, "2.5mm": true,
		"2mm": true, "1mm": true,
	}
	return validBands[strings.ToLower(band)]
}
