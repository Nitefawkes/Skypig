package services

import (
	"context"
	"fmt"

	"github.com/Nitefawkes/Skypig/backend/internal/database"
	"github.com/Nitefawkes/Skypig/backend/internal/models"
	"github.com/Nitefawkes/Skypig/backend/pkg/adif"
)

// ADIFService handles ADIF import/export operations
type ADIFService struct {
	qsoRepo  *database.QSORepository
	userRepo *database.UserRepository
	qsoSvc   *QSOService
}

// NewADIFService creates a new ADIF service
func NewADIFService(qsoRepo *database.QSORepository, userRepo *database.UserRepository, qsoSvc *QSOService) *ADIFService {
	return &ADIFService{
		qsoRepo:  qsoRepo,
		userRepo: userRepo,
		qsoSvc:   qsoSvc,
	}
}

// ImportResult contains import statistics
type ImportResult struct {
	TotalRecords    int      `json:"total_records"`
	ImportedRecords int      `json:"imported_records"`
	FailedRecords   int      `json:"failed_records"`
	Errors          []string `json:"errors,omitempty"`
	SkippedRecords  int      `json:"skipped_records"`
}

// ImportADIF imports QSOs from ADIF content
func (s *ADIFService) ImportADIF(ctx context.Context, userID int64, content string, strict bool) (*ImportResult, error) {
	// Get user to check limits
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse ADIF content
	parser := adif.NewParser(false) // Always use non-strict parsing
	qsos, err := parser.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ADIF: %w", err)
	}

	result := &ImportResult{
		TotalRecords: len(qsos),
		Errors:       []string{},
	}

	// Import each QSO
	for i, qso := range qsos {
		// Check if user can create more QSOs
		if !user.CanCreateQSO() {
			result.SkippedRecords++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Record %d: QSO limit reached (%d/%d)", i+1, user.QSOCount, user.QSOLimit))

			if strict {
				return result, fmt.Errorf("QSO limit reached after importing %d records", result.ImportedRecords)
			}
			continue
		}

		// Set user ID
		qso.UserID = userID

		// Validate and create QSO
		_, err := s.qsoSvc.CreateQSO(ctx, userID, qso)
		if err != nil {
			result.FailedRecords++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Record %d (%s): %v", i+1, qso.Callsign, err))

			if strict {
				return result, fmt.Errorf("import failed at record %d: %w", i+1, err)
			}
			continue
		}

		result.ImportedRecords++
		user.QSOCount++ // Increment for limit checking
	}

	return result, nil
}

// ExportADIF exports QSOs to ADIF format
func (s *ADIFService) ExportADIF(ctx context.Context, userID int64, filter models.QSOFilter) (string, error) {
	// Get QSOs with filter
	qsos, _, err := s.qsoRepo.List(ctx, userID, filter)
	if err != nil {
		return "", fmt.Errorf("failed to list QSOs: %w", err)
	}

	// Generate ADIF
	generator := adif.NewGenerator()
	adifContent := generator.Generate(qsos)

	return adifContent, nil
}

// ValidateADIF validates ADIF content without importing
func (s *ADIFService) ValidateADIF(content string) (*ImportResult, error) {
	parser := adif.NewParser(true) // Strict mode for validation
	qsos, err := parser.Parse(content)

	result := &ImportResult{
		TotalRecords:    len(qsos),
		ImportedRecords: 0,
		Errors:          []string{},
	}

	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		return result, err
	}

	// Validate each QSO (without saving)
	for i, qso := range qsos {
		if err := s.validateQSO(qso); err != nil {
			result.FailedRecords++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Record %d (%s): %v", i+1, qso.Callsign, err))
		}
	}

	if result.FailedRecords > 0 {
		return result, fmt.Errorf("validation failed: %d records have errors", result.FailedRecords)
	}

	result.ImportedRecords = result.TotalRecords
	return result, nil
}

// validateQSO validates a QSO without creating it
func (s *ADIFService) validateQSO(qso *models.QSO) error {
	// Reuse validation from QSOService
	// This is a basic implementation - in production, you'd extract validation
	// to a shared function
	if qso.Callsign == "" {
		return fmt.Errorf("callsign is required")
	}
	if qso.TimeOn.IsZero() {
		return fmt.Errorf("time_on is required")
	}
	return nil
}
