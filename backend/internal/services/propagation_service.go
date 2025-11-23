package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Nitefawkes/Skypig/backend/internal/models"
)

// PropagationService handles propagation data fetching and analysis
type PropagationService struct {
	cache      *models.PropagationData
	cacheMutex sync.RWMutex
	cacheExpiry time.Time
	httpClient *http.Client
	useMock    bool
}

// NewPropagationService creates a new propagation service
func NewPropagationService(useMock bool) *PropagationService {
	return &PropagationService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		useMock: useMock,
	}
}

// GetCurrent returns current propagation conditions
func (s *PropagationService) GetCurrent(ctx context.Context) (*models.PropagationData, error) {
	// Check cache first (15 minute TTL)
	s.cacheMutex.RLock()
	if s.cache != nil && time.Now().Before(s.cacheExpiry) {
		defer s.cacheMutex.RUnlock()
		return s.cache, nil
	}
	s.cacheMutex.RUnlock()

	// Fetch new data
	data, err := s.fetchPropagationData(ctx)
	if err != nil {
		// Return cached data if available, even if expired
		s.cacheMutex.RLock()
		defer s.cacheMutex.RUnlock()
		if s.cache != nil {
			return s.cache, nil
		}
		return nil, err
	}

	// Update cache
	s.cacheMutex.Lock()
	s.cache = data
	s.cacheExpiry = time.Now().Add(15 * time.Minute)
	s.cacheMutex.Unlock()

	return data, nil
}

// GetForecast returns propagation forecast with band conditions
func (s *PropagationService) GetForecast(ctx context.Context) (*models.PropagationForecast, error) {
	current, err := s.GetCurrent(ctx)
	if err != nil {
		return nil, err
	}

	// Analyze band conditions
	bandConditions := s.analyzeBandConditions(current)

	forecast := &models.PropagationForecast{
		Current:        current,
		BandConditions: bandConditions,
		Summary:        s.generateSummary(current, bandConditions),
		LastUpdated:    current.UpdatedAt,
		NextUpdate:     time.Now().Add(15 * time.Minute),
	}

	return forecast, nil
}

// fetchPropagationData fetches data from external API or mock
func (s *PropagationService) fetchPropagationData(ctx context.Context) (*models.PropagationData, error) {
	if s.useMock {
		return s.getMockData(), nil
	}

	// Try HamQSL API first (free, no auth required)
	data, err := s.fetchFromHamQSL(ctx)
	if err == nil {
		return data, nil
	}

	// Fallback to mock data
	return s.getMockData(), nil
}

// fetchFromHamQSL fetches data from HamQSL propagation API
func (s *PropagationService) fetchFromHamQSL(ctx context.Context) (*models.PropagationData, error) {
	// HamQSL provides solar data in JSON format
	url := "https://www.hamqsl.com/solarxml.php"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HamQSL API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse XML response (simplified - would need proper XML parsing)
	// For MVP, we'll use mock data
	// TODO: Implement proper XML parsing
	_ = body

	return s.getMockData(), nil
}

// getMockData returns realistic mock propagation data
func (s *PropagationService) getMockData() *models.PropagationData {
	now := time.Now()

	// Vary data based on time of day for realism
	hour := now.Hour()
	var sfi float64
	var kIndex int

	// Simulate day/night variation
	if hour >= 6 && hour < 18 {
		// Daytime: better conditions
		sfi = 120 + float64(hour%12)*5
		kIndex = 2
	} else {
		// Nighttime: varied conditions
		sfi = 100 + float64(hour%12)*3
		kIndex = 3
	}

	return &models.PropagationData{
		Timestamp:        now,
		SolarFlux:        sfi,
		SunspotNumber:    45,
		AIndex:           8,
		KIndex:           kIndex,
		XRayFlux:         "C1.2",
		SolarWind:        380.5,
		BzComponent:      -2.3,
		ProtonFlux:       0.5,
		ElectronFlux:     1.2,
		GeomagneticStorm: "None",
		RadioBlackout:    "None",
		SolarRadiation:   "None",
		UpdatedAt:        now,
		Source:           "Mock Data",
	}
}

// analyzeBandConditions analyzes conditions for each HF band
func (s *PropagationService) analyzeBandConditions(data *models.PropagationData) []*models.BandConditions {
	hour := time.Now().UTC().Hour()
	isDay := hour >= 6 && hour < 18

	bands := []string{"80m", "40m", "30m", "20m", "17m", "15m", "12m", "10m"}
	conditions := make([]*models.BandConditions, 0, len(bands))

	for _, band := range bands {
		condition := &models.BandConditions{
			Band: band,
		}

		// Rule-based analysis
		switch band {
		case "80m", "40m":
			// Low bands: better at night
			if isDay {
				condition.Day = s.assessCondition(data.SolarFlux, data.KIndex, -1)
				condition.Night = s.assessCondition(data.SolarFlux, data.KIndex, 2)
			} else {
				condition.Day = s.assessCondition(data.SolarFlux, data.KIndex, 1)
				condition.Night = s.assessCondition(data.SolarFlux, data.KIndex, 1)
			}
			condition.Reasoning = "Low bands favor nighttime propagation"

		case "30m":
			condition.Day = s.assessCondition(data.SolarFlux, data.KIndex, 0)
			condition.Night = s.assessCondition(data.SolarFlux, data.KIndex, 1)
			condition.Reasoning = "Good for both day and night with current conditions"

		case "20m", "17m":
			// Mid bands: good during day with solar activity
			condition.Day = s.assessCondition(data.SolarFlux, data.KIndex, 2)
			condition.Night = s.assessCondition(data.SolarFlux, data.KIndex, -1)
			condition.Reasoning = "Daytime DX band, depends on solar flux"

		case "15m", "12m", "10m":
			// High bands: need high solar activity
			modifier := 0
			if data.SolarFlux > 120 {
				modifier = 2
			} else if data.SolarFlux > 100 {
				modifier = 1
			} else {
				modifier = -2
			}
			condition.Day = s.assessCondition(data.SolarFlux, data.KIndex, modifier)
			condition.Night = "poor"
			condition.Reasoning = fmt.Sprintf("High bands need solar flux >120 (current: %.0f)", data.SolarFlux)
		}

		conditions = append(conditions, condition)
	}

	return conditions
}

// assessCondition determines band condition based on solar flux and K-index
func (s *PropagationService) assessCondition(solarFlux float64, kIndex int, modifier int) string {
	score := 0

	// Solar flux contribution
	if solarFlux > 150 {
		score += 3
	} else if solarFlux > 120 {
		score += 2
	} else if solarFlux > 90 {
		score += 1
	}

	// K-index contribution (lower is better)
	if kIndex <= 2 {
		score += 2
	} else if kIndex <= 4 {
		score += 0
	} else {
		score -= 2
	}

	// Apply modifier
	score += modifier

	// Determine condition
	if score >= 4 {
		return "excellent"
	} else if score >= 2 {
		return "good"
	} else if score >= 0 {
		return "fair"
	}
	return "poor"
}

// generateSummary creates a text summary of current conditions
func (s *PropagationService) generateSummary(data *models.PropagationData, bands []*models.BandConditions) string {
	summary := fmt.Sprintf("Solar Flux: %.0f (%s), K-Index: %d (%s). ",
		data.SolarFlux,
		data.SolarActivity(),
		data.KIndex,
		data.GeomagneticActivity(),
	)

	// Find best bands
	bestBands := []string{}
	for _, band := range bands {
		if band.Day == "excellent" || band.Day == "good" {
			bestBands = append(bestBands, band.Band)
		}
	}

	if len(bestBands) > 0 {
		summary += fmt.Sprintf("Best daytime bands: %v. ", bestBands)
	}

	summary += fmt.Sprintf("Overall HF conditions: %s.", data.HFConditions())

	return summary
}

// HamQSLResponse represents the structure of HamQSL API response
type HamQSLResponse struct {
	SolarData struct {
		SolarFlux     string `json:"solarflux"`
		AIndex        string `json:"aindex"`
		KIndex        string `json:"kindex"`
		SunspotNumber string `json:"sunspots"`
		Updated       string `json:"updated"`
	} `json:"solardata"`
}

// parseHamQSL parses HamQSL JSON response
func (s *PropagationService) parseHamQSL(body []byte) (*models.PropagationData, error) {
	var resp HamQSLResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// Convert strings to appropriate types
	// TODO: Add proper parsing with error handling

	return &models.PropagationData{
		UpdatedAt: time.Now(),
		Source:    "HamQSL",
	}, nil
}
