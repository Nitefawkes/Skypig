package models

import (
	"time"
)

// User represents a ham radio operator account
type User struct {
	ID           int64     `json:"id"`
	Callsign     string    `json:"callsign"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	GridSquare   string    `json:"grid_square"`
	QRZVerified  bool      `json:"qrz_verified"`
	Tier         string    `json:"tier"` // "free", "operator", "contester"
	StripeID     string    `json:"-"`
	QSOLimit     int       `json:"qso_limit"`
	QSOCount     int       `json:"qso_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
}

// CanCreateQSO checks if user has capacity to create a new QSO
func (u *User) CanCreateQSO() bool {
	if u.Tier == "contester" {
		return true // Unlimited
	}
	return u.QSOCount < u.QSOLimit
}

// GetQSOLimit returns the QSO limit for the user's tier
func (u *User) GetQSOLimit() int {
	switch u.Tier {
	case "free":
		return 500
	case "operator":
		return 20000
	case "contester":
		return -1 // Unlimited
	default:
		return 500
	}
}
