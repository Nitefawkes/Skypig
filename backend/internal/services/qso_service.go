package services

import (
	"fmt"
	"time"

	"github.com/nitefawkes/ham-radio-cloud/internal/models"
	"github.com/nitefawkes/ham-radio-cloud/internal/repositories"
)

type QSOService struct {
	repo *repositories.QSORepository
}

func NewQSOService(repo *repositories.QSORepository) *QSOService {
	return &QSOService{repo: repo}
}

func (s *QSOService) GetQSOs(userID string, filter *models.QSOFilter) ([]models.QSO, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	qsos, err := s.repo.GetByUserID(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get QSOs: %w", err)
	}

	return qsos, nil
}

func (s *QSOService) CreateQSO(qso *models.QSO) error {
	if err := s.validateQSO(qso); err != nil {
		return err
	}

	// Set defaults
	if qso.QSODate.IsZero() {
		qso.QSODate = time.Now()
	}
	if qso.TimeOn.IsZero() {
		qso.TimeOn = time.Now()
	}

	if err := s.repo.Create(qso); err != nil {
		return fmt.Errorf("failed to create QSO: %w", err)
	}

	return nil
}

func (s *QSOService) UpdateQSO(qso *models.QSO) error {
	if err := s.validateQSO(qso); err != nil {
		return err
	}

	if err := s.repo.Update(qso); err != nil {
		return fmt.Errorf("failed to update QSO: %w", err)
	}

	return nil
}

func (s *QSOService) DeleteQSO(id, userID string) error {
	if id == "" || userID == "" {
		return fmt.Errorf("QSO ID and user ID are required")
	}

	if err := s.repo.Delete(id, userID); err != nil {
		return fmt.Errorf("failed to delete QSO: %w", err)
	}

	return nil
}

func (s *QSOService) GetQSOCount(userID string) (int, error) {
	count, err := s.repo.GetCount(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get QSO count: %w", err)
	}
	return count, nil
}

func (s *QSOService) validateQSO(qso *models.QSO) error {
	if qso.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if qso.Callsign == "" {
		return fmt.Errorf("callsign is required")
	}
	if qso.Mode == "" {
		return fmt.Errorf("mode is required")
	}
	if qso.Band == "" && qso.Frequency == 0 {
		return fmt.Errorf("either band or frequency is required")
	}

	// Auto-determine band from frequency if not provided
	if qso.Band == "" && qso.Frequency > 0 {
		qso.Band = s.determineBand(qso.Frequency)
	}

	return nil
}

func (s *QSOService) determineBand(frequency float64) string {
	switch {
	case frequency >= 1.8 && frequency < 2.0:
		return "160m"
	case frequency >= 3.5 && frequency < 4.0:
		return "80m"
	case frequency >= 7.0 && frequency < 7.3:
		return "40m"
	case frequency >= 10.1 && frequency < 10.15:
		return "30m"
	case frequency >= 14.0 && frequency < 14.35:
		return "20m"
	case frequency >= 18.068 && frequency < 18.168:
		return "17m"
	case frequency >= 21.0 && frequency < 21.45:
		return "15m"
	case frequency >= 24.89 && frequency < 24.99:
		return "12m"
	case frequency >= 28.0 && frequency < 29.7:
		return "10m"
	case frequency >= 50.0 && frequency < 54.0:
		return "6m"
	case frequency >= 144.0 && frequency < 148.0:
		return "2m"
	case frequency >= 420.0 && frequency < 450.0:
		return "70cm"
	default:
		return "unknown"
	}
}
