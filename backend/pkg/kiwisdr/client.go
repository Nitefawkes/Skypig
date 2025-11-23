package kiwisdr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// KiwiSDR public directory API
	DirectoryURL = "http://kiwisdr.com/public/"
	DirectoryJSON = "http://proxy.kiwisdr.com:8073/status"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// KiwiSDREntry represents a single KiwiSDR receiver from the directory
type KiwiSDREntry struct {
	Name       string   `json:"name"`
	Callsign   string   `json:"sdr_hu_ant_callsign"`
	URL        string   `json:"url"`
	Location   string   `json:"sdr_hu_loc"`
	GridSquare string   `json:"sdr_hu_gps"`
	Latitude   float64  `json:"gps_lat"`
	Longitude  float64  `json:"gps_lon"`
	Antenna    string   `json:"sdr_hu_ant"`
	FreqMin    float64  `json:"freq_min_khz"`
	FreqMax    float64  `json:"freq_max_khz"`
	Users      int      `json:"users"`
	UsersMax   int      `json:"users_max"`
	Bands      []string `json:"sdr_hu_bands"`
	Status     string   `json:"status"`
	Offline    bool     `json:"offline"`
}

// DirectoryResponse represents the KiwiSDR directory JSON response
type DirectoryResponse struct {
	SDRs map[string]KiwiSDREntry `json:"sdrs"`
}

// GetDirectory fetches all KiwiSDR receivers from the public directory
func (c *Client) GetDirectory() ([]KiwiSDREntry, error) {
	resp, err := c.httpClient.Get(DirectoryJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch KiwiSDR directory: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("KiwiSDR directory returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// The response is a map of SDR ID -> SDR data
	var directoryMap map[string]KiwiSDREntry
	if err := json.Unmarshal(body, &directoryMap); err != nil {
		return nil, fmt.Errorf("failed to parse directory data: %w", err)
	}

	// Convert map to slice
	sdrs := make([]KiwiSDREntry, 0, len(directoryMap))
	for id, sdr := range directoryMap {
		// Set URL if not provided
		if sdr.URL == "" {
			sdr.URL = "http://" + id
		}

		// Parse bands if available
		if len(sdr.Bands) == 0 && sdr.FreqMin > 0 && sdr.FreqMax > 0 {
			sdr.Bands = determineBands(sdr.FreqMin, sdr.FreqMax)
		}

		sdrs = append(sdrs, sdr)
	}

	return sdrs, nil
}

// determineBands determines which amateur radio bands are covered by frequency range
func determineBands(minKHz, maxKHz float64) []string {
	type band struct {
		name  string
		start float64
		end   float64
	}

	bands := []band{
		{"160m", 1800, 2000},
		{"80m", 3500, 4000},
		{"60m", 5330, 5405},
		{"40m", 7000, 7300},
		{"30m", 10100, 10150},
		{"20m", 14000, 14350},
		{"17m", 18068, 18168},
		{"15m", 21000, 21450},
		{"12m", 24890, 24990},
		{"10m", 28000, 29700},
		{"6m", 50000, 54000},
		{"2m", 144000, 148000},
	}

	var covered []string
	for _, b := range bands {
		// Check if band overlaps with SDR frequency range
		if !(maxKHz < b.start || minKHz > b.end) {
			covered = append(covered, b.name)
		}
	}

	return covered
}

// NormalizeCallsign cleans up callsign formatting
func NormalizeCallsign(callsign string) string {
	if callsign == "" {
		return ""
	}
	callsign = strings.TrimSpace(strings.ToUpper(callsign))
	// Remove common suffixes like /KIWI, /SDR, etc.
	callsign = strings.Split(callsign, "/")[0]
	return callsign
}

// ParseGridSquare validates and normalizes Maidenhead grid square
func ParseGridSquare(grid string) string {
	if grid == "" {
		return ""
	}
	grid = strings.ToUpper(strings.TrimSpace(grid))
	// Basic validation: should be 4 or 6 characters
	if len(grid) != 4 && len(grid) != 6 {
		return ""
	}
	return grid
}
