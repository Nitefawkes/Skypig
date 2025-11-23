package services

import (
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/nitefawkes/ham-radio-cloud/internal/models"
	"github.com/nitefawkes/ham-radio-cloud/internal/repositories"
	"github.com/nitefawkes/ham-radio-cloud/pkg/kiwisdr"
)

type SDRService struct {
	repo          *repositories.SDRRepository
	kiwiSDRClient *kiwisdr.Client
}

func NewSDRService(repo *repositories.SDRRepository) *SDRService {
	return &SDRService{
		repo:          repo,
		kiwiSDRClient: kiwisdr.NewClient(),
	}
}

// RefreshDirectory fetches the latest SDR directory from KiwiSDR and updates the database
func (s *SDRService) RefreshDirectory() error {
	log.Println("🔄 Refreshing SDR directory from KiwiSDR network...")

	entries, err := s.kiwiSDRClient.GetDirectory()
	if err != nil {
		return fmt.Errorf("failed to fetch KiwiSDR directory: %w", err)
	}

	log.Printf("📡 Found %d KiwiSDR receivers", len(entries))

	// Convert KiwiSDR entries to our SDR model
	sdrs := make([]models.SDRReceiver, 0, len(entries))
	for _, entry := range entries {
		sdr := s.convertKiwiSDREntry(entry)
		sdrs = append(sdrs, sdr)
	}

	// Bulk upsert to database
	if err := s.repo.BulkUpsert(sdrs); err != nil {
		return fmt.Errorf("failed to update SDR database: %w", err)
	}

	log.Printf("✅ SDR directory updated: %d receivers", len(sdrs))
	return nil
}

// convertKiwiSDREntry converts a KiwiSDR directory entry to our SDR model
func (s *SDRService) convertKiwiSDREntry(entry kiwisdr.KiwiSDREntry) models.SDRReceiver {
	now := time.Now()

	// Determine status
	status := "online"
	if entry.Offline {
		status = "offline"
	}

	// Normalize callsign
	callsign := kiwisdr.NormalizeCallsign(entry.Callsign)
	var callsignPtr *string
	if callsign != "" {
		callsignPtr = &callsign
	}

	// Normalize grid square
	gridSquare := kiwisdr.ParseGridSquare(entry.GridSquare)
	var gridPtr *string
	if gridSquare != "" {
		gridPtr = &gridSquare
	}

	// Location and country
	var locationPtr *string
	if entry.Location != "" {
		locationPtr = &entry.Location
	}

	// Parse country from location (simplified - would need geocoding for accuracy)
	var countryPtr *string
	if entry.Location != "" {
		country := parseCountryFromLocation(entry.Location)
		if country != "" {
			countryPtr = &country
		}
	}

	// Antenna info
	var antennaPtr *string
	if entry.Antenna != "" {
		antennaPtr = &entry.Antenna
	}

	// Convert frequency from kHz to MHz
	var freqMin, freqMax *float64
	if entry.FreqMin > 0 {
		fm := entry.FreqMin / 1000.0 // Convert kHz to MHz
		freqMin = &fm
	}
	if entry.FreqMax > 0 {
		fm := entry.FreqMax / 1000.0 // Convert kHz to MHz
		freqMax = &fm
	}

	// Users max
	var usersMaxPtr *int
	if entry.UsersMax > 0 {
		usersMaxPtr = &entry.UsersMax
	}

	// Latitude/Longitude
	var lat, lon *float64
	if entry.Latitude != 0 {
		lat = &entry.Latitude
	}
	if entry.Longitude != 0 {
		lon = &entry.Longitude
	}

	// Bands
	bands := pq.StringArray(entry.Bands)
	if len(bands) == 0 {
		bands = pq.StringArray{}
	}

	// Modes (KiwiSDR supports all modes, but default to common ones)
	modes := pq.StringArray{"AM", "SSB", "CW", "FM"}

	return models.SDRReceiver{
		Name:         entry.Name,
		Callsign:     callsignPtr,
		URL:          entry.URL,
		Type:         "kiwisdr",
		Location:     locationPtr,
		GridSquare:   gridPtr,
		Latitude:     lat,
		Longitude:    lon,
		Country:      countryPtr,
		Bands:        bands,
		Modes:        modes,
		AntennaInfo:  antennaPtr,
		FrequencyMin: freqMin,
		FrequencyMax: freqMax,
		UsersMax:     usersMaxPtr,
		Status:       status,
		LastSeen:     &now,
		IsPublic:     true,
	}
}

// parseCountryFromLocation attempts to extract country from location string
func parseCountryFromLocation(location string) string {
	// This is a simplified version - would need proper geocoding for production
	// For now, just return the last part after comma
	if location == "" {
		return ""
	}

	parts := []rune(location)
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == ',' {
			if i+2 < len(parts) {
				return string(parts[i+2:])
			}
		}
	}

	return ""
}

// List returns SDRs based on filter criteria
func (s *SDRService) List(filter *models.SDRFilter) ([]models.SDRReceiver, int, error) {
	// Set default limit if not specified
	if filter.Limit == 0 {
		filter.Limit = 50
	}

	sdrs, err := s.repo.List(filter)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(filter)
	if err != nil {
		return nil, 0, err
	}

	return sdrs, count, nil
}

// GetByID returns a single SDR by ID
func (s *SDRService) GetByID(id string) (*models.SDRReceiver, error) {
	// Parse UUID
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid SDR ID: %w", err)
	}

	return s.repo.GetByID(uuid)
}

// Search searches SDRs by name, location, or callsign
func (s *SDRService) Search(query string, limit int) ([]models.SDRReceiver, error) {
	if limit == 0 {
		limit = 20
	}

	filter := &models.SDRFilter{
		Search: query,
		Limit:  limit,
		Status: "online", // Only return online SDRs for search
	}

	sdrs, err := s.repo.List(filter)
	if err != nil {
		return nil, err
	}

	return sdrs, nil
}
