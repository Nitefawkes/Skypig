package models

import "time"

type PropagationData struct {
	ID          string    `json:"id" db:"id"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	SolarFlux   float64   `json:"solar_flux" db:"solar_flux"`
	SunspotNumber int     `json:"sunspot_number" db:"sunspot_number"`
	AIndex      int       `json:"a_index" db:"a_index"`
	KIndex      int       `json:"k_index" db:"k_index"`
	XRayFlux    string    `json:"xray_flux" db:"xray_flux"`
	HeliumLine  float64   `json:"helium_line" db:"helium_line"`
	ProtonFlux  int       `json:"proton_flux" db:"proton_flux"`
	ElectronFlux int      `json:"electron_flux" db:"electron_flux"`
	Source      string    `json:"source" db:"source"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type BandCondition struct {
	Band      string  `json:"band"`
	Condition string  `json:"condition"` // good, fair, poor
	Score     float64 `json:"score"`
	DayNight  string  `json:"day_night"`
}

type PropagationForecast struct {
	Timestamp      time.Time       `json:"timestamp"`
	BandConditions []BandCondition `json:"band_conditions"`
	Summary        string          `json:"summary"`
	Confidence     float64         `json:"confidence"`
}
