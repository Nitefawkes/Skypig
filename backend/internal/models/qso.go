package models

import "time"

type QSO struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	Callsign       string    `json:"callsign" db:"callsign"`
	Frequency      float64   `json:"frequency" db:"frequency"`
	Band           string    `json:"band" db:"band"`
	Mode           string    `json:"mode" db:"mode"`
	RST_Sent       string    `json:"rst_sent" db:"rst_sent"`
	RST_Received   string    `json:"rst_received" db:"rst_received"`
	QSODate        time.Time `json:"qso_date" db:"qso_date"`
	TimeOn         time.Time `json:"time_on" db:"time_on"`
	TimeOff        *time.Time `json:"time_off,omitempty" db:"time_off"`
	GridSquare     string    `json:"grid_square" db:"grid_square"`
	Country        string    `json:"country" db:"country"`
	State          string    `json:"state" db:"state"`
	County         string    `json:"county" db:"county"`
	Comment        string    `json:"comment" db:"comment"`
	ContestID      string    `json:"contest_id" db:"contest_id"`
	PropagationMode string   `json:"propagation_mode" db:"propagation_mode"`
	SatelliteName  string    `json:"satellite_name" db:"satellite_name"`
	TXPower        float64   `json:"tx_power" db:"tx_power"`
	LoTWSent       bool      `json:"lotw_sent" db:"lotw_sent"`
	LoTWConfirmed  bool      `json:"lotw_confirmed" db:"lotw_confirmed"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type QSOFilter struct {
	StartDate  *time.Time
	EndDate    *time.Time
	Callsign   string
	Band       string
	Mode       string
	Limit      int
	Offset     int
}
