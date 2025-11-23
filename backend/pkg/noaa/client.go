package noaa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	SolarWindURL      = "https://services.swpc.noaa.gov/products/solar-wind/mag-1-day.json"
	SolarFluxURL      = "https://services.swpc.noaa.gov/json/f107_cm_flux.json"
	PlanetaryKURL     = "https://services.swpc.noaa.gov/products/noaa-planetary-k-index.json"
	SunspotNumberURL  = "https://services.swpc.noaa.gov/json/solar-cycle/observed-solar-cycle-indices.json"
	XRayFluxURL       = "https://services.swpc.noaa.gov/json/goes/primary/xrays-6-hour.json"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type SolarFluxData struct {
	TimeTag    string  `json:"time_tag"`
	Flux       float64 `json:"flux"`
}

type PlanetaryKData struct {
	TimeTag string `json:"time_tag"`
	KIndex  int    `json:"Kp"`
	AIndex  int    `json:"ap,omitempty"`
}

type SunspotData struct {
	TimeTag        string  `json:"time-tag"`
	SunspotNumber  float64 `json:"ssn"`
}

type XRayData struct {
	TimeTag string  `json:"time_tag"`
	Flux    float64 `json:"flux"`
	Energy  string  `json:"energy"`
}

func (c *Client) GetSolarFlux() (float64, error) {
	resp, err := c.httpClient.Get(SolarFluxURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch solar flux: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("NOAA API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var data []SolarFluxData
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("failed to parse solar flux data: %w", err)
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("no solar flux data available")
	}

	// Return most recent value
	return data[len(data)-1].Flux, nil
}

func (c *Client) GetPlanetaryK() (int, int, error) {
	resp, err := c.httpClient.Get(PlanetaryKURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch planetary K: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("NOAA API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	var data [][]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, 0, fmt.Errorf("failed to parse K-index data: %w", err)
	}

	if len(data) < 2 {
		return 0, 0, fmt.Errorf("no K-index data available")
	}

	// Get most recent K-index (last row, 3rd column)
	lastRow := data[len(data)-1]
	if len(lastRow) < 3 {
		return 0, 0, fmt.Errorf("invalid K-index data format")
	}

	kIndex := int(lastRow[1].(float64))

	// A-index is typically 3rd column if available
	aIndex := 0
	if len(lastRow) > 3 {
		aIndex = int(lastRow[3].(float64))
	}

	return kIndex, aIndex, nil
}

func (c *Client) GetSunspotNumber() (int, error) {
	resp, err := c.httpClient.Get(SunspotNumberURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch sunspot number: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("NOAA API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var data []SunspotData
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("failed to parse sunspot data: %w", err)
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("no sunspot data available")
	}

	// Return most recent value
	return int(data[len(data)-1].SunspotNumber), nil
}

func (c *Client) GetXRayFlux() (string, error) {
	resp, err := c.httpClient.Get(XRayFluxURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch X-ray flux: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("NOAA API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var data []XRayData
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("failed to parse X-ray data: %w", err)
	}

	if len(data) == 0 {
		return "N/A", nil
	}

	// Get most recent X-ray flux and classify
	flux := data[len(data)-1].Flux
	return classifyXRayFlux(flux), nil
}

func classifyXRayFlux(flux float64) string {
	// X-ray flux classification (W/m²)
	switch {
	case flux >= 1e-4:
		return "X" + fmt.Sprintf("%.1f", flux/1e-4)
	case flux >= 1e-5:
		return "M" + fmt.Sprintf("%.1f", flux/1e-5)
	case flux >= 1e-6:
		return "C" + fmt.Sprintf("%.1f", flux/1e-6)
	case flux >= 1e-7:
		return "B" + fmt.Sprintf("%.1f", flux/1e-7)
	default:
		return "A" + fmt.Sprintf("%.1f", flux/1e-8)
	}
}
