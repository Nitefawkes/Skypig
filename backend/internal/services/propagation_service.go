package services

import (
	"fmt"
	"log"
	"time"

	"github.com/nitefawkes/ham-radio-cloud/internal/models"
	"github.com/nitefawkes/ham-radio-cloud/internal/repositories"
	"github.com/nitefawkes/ham-radio-cloud/pkg/noaa"
)

type PropagationService struct {
	repo       *repositories.PropagationRepository
	noaaClient *noaa.Client
}

func NewPropagationService(repo *repositories.PropagationRepository) *PropagationService {
	return &PropagationService{
		repo:       repo,
		noaaClient: noaa.NewClient(),
	}
}

func (s *PropagationService) FetchAndStore() error {
	// Fetch data from NOAA
	solarFlux, err := s.noaaClient.GetSolarFlux()
	if err != nil {
		log.Printf("Warning: Failed to fetch solar flux: %v", err)
		solarFlux = 0
	}

	kIndex, aIndex, err := s.noaaClient.GetPlanetaryK()
	if err != nil {
		log.Printf("Warning: Failed to fetch K-index: %v", err)
		kIndex, aIndex = 0, 0
	}

	sunspotNumber, err := s.noaaClient.GetSunspotNumber()
	if err != nil {
		log.Printf("Warning: Failed to fetch sunspot number: %v", err)
		sunspotNumber = 0
	}

	xrayFlux, err := s.noaaClient.GetXRayFlux()
	if err != nil {
		log.Printf("Warning: Failed to fetch X-ray flux: %v", err)
		xrayFlux = "N/A"
	}

	// Create propagation data record
	data := &models.PropagationData{
		Timestamp:     time.Now().UTC(),
		SolarFlux:     solarFlux,
		SunspotNumber: sunspotNumber,
		AIndex:        aIndex,
		KIndex:        kIndex,
		XRayFlux:      xrayFlux,
		HeliumLine:    0, // Not available from NOAA free API
		ProtonFlux:    0, // Not available from NOAA free API
		ElectronFlux:  0, // Not available from NOAA free API
		Source:        "NOAA SWPC",
	}

	if err := s.repo.Create(data); err != nil {
		return fmt.Errorf("failed to store propagation data: %w", err)
	}

	log.Printf("✅ Propagation data updated: SFI=%.0f, SSN=%d, K=%d, A=%d",
		solarFlux, sunspotNumber, kIndex, aIndex)

	return nil
}

func (s *PropagationService) GetCurrent() (*models.PropagationData, error) {
	data, err := s.repo.GetLatest()
	if err != nil {
		// If no data exists, try to fetch it
		if fetchErr := s.FetchAndStore(); fetchErr != nil {
			return nil, fmt.Errorf("no propagation data available and fetch failed: %w", fetchErr)
		}
		// Try again after fetching
		data, err = s.repo.GetLatest()
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (s *PropagationService) GetBandConditions(data *models.PropagationData) []models.BandCondition {
	if data == nil {
		return []models.BandCondition{}
	}

	// Determine if it's day or night (simplified - would need user location for accuracy)
	hour := time.Now().UTC().Hour()
	isDaytime := hour >= 6 && hour < 18
	dayNight := "day"
	if !isDaytime {
		dayNight = "night"
	}

	conditions := []models.BandCondition{
		s.calculateBandCondition("160m", data, dayNight),
		s.calculateBandCondition("80m", data, dayNight),
		s.calculateBandCondition("40m", data, dayNight),
		s.calculateBandCondition("30m", data, dayNight),
		s.calculateBandCondition("20m", data, dayNight),
		s.calculateBandCondition("17m", data, dayNight),
		s.calculateBandCondition("15m", data, dayNight),
		s.calculateBandCondition("12m", data, dayNight),
		s.calculateBandCondition("10m", data, dayNight),
		s.calculateBandCondition("6m", data, dayNight),
	}

	return conditions
}

func (s *PropagationService) calculateBandCondition(band string, data *models.PropagationData, dayNight string) models.BandCondition {
	// Simplified band condition calculator
	// Real implementation would use VOACAP or more sophisticated models

	score := 5.0 // Base score

	// Solar flux impact (higher is better for HF)
	if data.SolarFlux > 150 {
		score += 2
	} else if data.SolarFlux > 100 {
		score += 1
	} else if data.SolarFlux < 70 {
		score -= 1
	}

	// K-index impact (lower is better)
	if data.KIndex <= 2 {
		score += 1
	} else if data.KIndex >= 5 {
		score -= 2
	} else if data.KIndex >= 4 {
		score -= 1
	}

	// A-index impact (lower is better)
	if data.AIndex <= 7 {
		score += 0.5
	} else if data.AIndex > 20 {
		score -= 1.5
	}

	// Band-specific adjustments
	switch band {
	case "160m", "80m":
		// Low bands: better at night
		if dayNight == "night" {
			score += 2
		} else {
			score -= 3
		}
		// Less dependent on solar flux
		if data.SolarFlux < 70 {
			score += 0.5
		}

	case "40m", "30m":
		// Mid bands: work day and night
		if dayNight == "night" {
			score += 1
		}

	case "20m", "17m":
		// Classic DX bands: need good conditions
		if dayNight == "day" {
			score += 1
		}
		// Very dependent on solar flux
		if data.SolarFlux > 120 {
			score += 1
		} else if data.SolarFlux < 80 {
			score -= 2
		}

	case "15m", "12m", "10m":
		// High bands: need high solar flux
		if data.SolarFlux > 150 {
			score += 2
		} else if data.SolarFlux > 120 {
			score += 1
		} else if data.SolarFlux < 100 {
			score -= 3
		}
		// Daytime is better
		if dayNight == "day" {
			score += 1
		} else {
			score -= 1
		}

	case "6m":
		// Sporadic E and solar dependent
		if data.SolarFlux > 150 {
			score += 1
		}
	}

	// Clamp score to 0-10 range
	if score > 10 {
		score = 10
	} else if score < 0 {
		score = 0
	}

	// Determine condition category
	condition := "fair"
	if score >= 7 {
		condition = "good"
	} else if score < 4 {
		condition = "poor"
	}

	return models.BandCondition{
		Band:      band,
		Condition: condition,
		Score:     score,
		DayNight:  dayNight,
	}
}

func (s *PropagationService) CleanupOldData(daysToKeep int) error {
	cutoff := time.Now().AddDate(0, 0, -daysToKeep)
	return s.repo.DeleteOlderThan(cutoff)
}
