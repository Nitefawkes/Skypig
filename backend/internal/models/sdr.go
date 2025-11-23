package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type SDRReceiver struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	Name         string         `json:"name" db:"name"`
	Callsign     *string        `json:"callsign,omitempty" db:"callsign"`
	URL          string         `json:"url" db:"url"`
	Type         string         `json:"type" db:"type"`
	Location     *string        `json:"location,omitempty" db:"location"`
	GridSquare   *string        `json:"grid_square,omitempty" db:"grid_square"`
	Latitude     *float64       `json:"latitude,omitempty" db:"latitude"`
	Longitude    *float64       `json:"longitude,omitempty" db:"longitude"`
	Country      *string        `json:"country,omitempty" db:"country"`
	Bands        pq.StringArray `json:"bands" db:"bands"`
	Modes        pq.StringArray `json:"modes" db:"modes"`
	AntennaInfo  *string        `json:"antenna_info,omitempty" db:"antenna_info"`
	FrequencyMin *float64       `json:"frequency_min,omitempty" db:"frequency_min"`
	FrequencyMax *float64       `json:"frequency_max,omitempty" db:"frequency_max"`
	UsersMax     *int           `json:"users_max,omitempty" db:"users_max"`
	Status       string         `json:"status" db:"status"`
	LastSeen     *time.Time     `json:"last_seen,omitempty" db:"last_seen"`
	Description  *string        `json:"description,omitempty" db:"description"`
	AvatarURL    *string        `json:"avatar_url,omitempty" db:"avatar_url"`
	IsPublic     bool           `json:"is_public" db:"is_public"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}

type UserFavoriteSDR struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	SDRID     uuid.UUID `json:"sdr_id" db:"sdr_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SDRFilter represents filter criteria for SDR queries
type SDRFilter struct {
	Type       string   `json:"type,omitempty"`
	Country    string   `json:"country,omitempty"`
	Status     string   `json:"status,omitempty"`
	Bands      []string `json:"bands,omitempty"`
	Search     string   `json:"search,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`
}
