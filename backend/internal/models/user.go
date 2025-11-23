package models

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	Callsign  string    `json:"callsign" db:"callsign"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	QRZUserID string    `json:"qrz_user_id" db:"qrz_user_id"`
	Tier      string    `json:"tier" db:"tier"` // free, operator, contester
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserSettings struct {
	UserID              string `json:"user_id" db:"user_id"`
	LoTWEnabled         bool   `json:"lotw_enabled" db:"lotw_enabled"`
	LoTWUsername        string `json:"lotw_username" db:"lotw_username"`
	PropagationAlerts   bool   `json:"propagation_alerts" db:"propagation_alerts"`
	PreferredBands      string `json:"preferred_bands" db:"preferred_bands"`
	GridSquare          string `json:"grid_square" db:"grid_square"`
	TimeZone            string `json:"time_zone" db:"time_zone"`
}
