package models

import "time"

// PropagationData represents current space weather and propagation conditions
type PropagationData struct {
	ID              int64     `json:"id"`
	Timestamp       time.Time `json:"timestamp"`
	SolarFlux       float64   `json:"solar_flux"`        // SFI (Solar Flux Index)
	SunspotNumber   int       `json:"sunspot_number"`    // SSN
	AIndex          int       `json:"a_index"`           // Planetary A-index (24h)
	KIndex          int       `json:"k_index"`           // Planetary K-index (3h)
	XRayFlux        string    `json:"x_ray_flux"`        // e.g., "M1.2", "C5.4"
	SolarWind       float64   `json:"solar_wind"`        // km/s
	BzComponent     float64   `json:"bz_component"`      // nT (Interplanetary Magnetic Field)
	ProtonFlux      float64   `json:"proton_flux"`       // particles/cm²/s/sr
	ElectronFlux    float64   `json:"electron_flux"`     // particles/cm²/s/sr
	GeomagneticStorm string   `json:"geomagnetic_storm"` // "None", "Minor", "Moderate", "Strong", "Severe"
	RadioBlackout   string    `json:"radio_blackout"`    // "None", "Minor", "Moderate", "Strong", "Severe"
	SolarRadiation  string    `json:"solar_radiation"`   // "None", "Minor", "Moderate", "Strong", "Severe"
	UpdatedAt       time.Time `json:"updated_at"`
	Source          string    `json:"source"` // Data source (e.g., "NOAA", "HamQSL")
}

// BandConditions represents HF band conditions
type BandConditions struct {
	Band      string `json:"band"`      // e.g., "80m", "40m", "20m"
	Day       string `json:"day"`       // "poor", "fair", "good", "excellent"
	Night     string `json:"night"`     // "poor", "fair", "good", "excellent"
	Reasoning string `json:"reasoning"` // Brief explanation
}

// PropagationForecast represents propagation conditions with band analysis
type PropagationForecast struct {
	Current        *PropagationData   `json:"current"`
	BandConditions []*BandConditions  `json:"band_conditions"`
	Summary        string             `json:"summary"`
	LastUpdated    time.Time          `json:"last_updated"`
	NextUpdate     time.Time          `json:"next_update"`
}

// SolarActivity returns a text description of current solar activity
func (p *PropagationData) SolarActivity() string {
	if p.SolarFlux < 70 {
		return "Very Low"
	} else if p.SolarFlux < 100 {
		return "Low"
	} else if p.SolarFlux < 150 {
		return "Moderate"
	} else if p.SolarFlux < 200 {
		return "High"
	}
	return "Very High"
}

// GeomagneticActivity returns a text description of geomagnetic activity
func (p *PropagationData) GeomagneticActivity() string {
	if p.KIndex <= 2 {
		return "Quiet"
	} else if p.KIndex <= 4 {
		return "Unsettled"
	} else if p.KIndex <= 6 {
		return "Active"
	} else if p.KIndex <= 8 {
		return "Storm"
	}
	return "Severe Storm"
}

// HFConditions returns overall HF condition assessment
func (p *PropagationData) HFConditions() string {
	// Simple rule-based assessment
	// Good: High SFI, Low K
	// Poor: Low SFI, High K

	score := 0

	// Solar flux contribution
	if p.SolarFlux > 150 {
		score += 3
	} else if p.SolarFlux > 100 {
		score += 2
	} else if p.SolarFlux > 70 {
		score += 1
	}

	// K-index contribution (inverse - lower is better)
	if p.KIndex <= 2 {
		score += 3
	} else if p.KIndex <= 4 {
		score += 1
	} else {
		score -= 2
	}

	// A-index contribution (inverse)
	if p.AIndex <= 7 {
		score += 1
	} else if p.AIndex > 20 {
		score -= 1
	}

	// Interpret score
	if score >= 6 {
		return "excellent"
	} else if score >= 4 {
		return "good"
	} else if score >= 2 {
		return "fair"
	}
	return "poor"
}
