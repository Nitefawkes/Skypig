package models

import (
	"time"
)

// QSO represents a single contact/QSO record
type QSO struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Callsign        string    `json:"callsign"`        // Contact's callsign
	OperatorCall    string    `json:"operator_call"`   // User's callsign at time of QSO
	StationCallsign string    `json:"station_callsign"` // Station callsign if different
	QSODate         time.Time `json:"qso_date"`
	TimeOn          time.Time `json:"time_on"`
	TimeOff         time.Time `json:"time_off"`
	Band            string    `json:"band"`     // e.g., "20m", "40m"
	BandRX          string    `json:"band_rx"`  // Receive band if split
	Freq            float64   `json:"freq"`     // Frequency in MHz
	FreqRX          float64   `json:"freq_rx"`  // RX frequency if split
	Mode            string    `json:"mode"`     // e.g., "SSB", "CW", "FT8"
	Submode         string    `json:"submode"`  // e.g., "USB", "LSB"
	RSTSent         string    `json:"rst_sent"` // Signal report sent
	RSTRcvd         string    `json:"rst_rcvd"` // Signal report received
	Name            string    `json:"name"`     // Contact's name
	QTH             string    `json:"qth"`      // Contact's location
	GridSquare      string    `json:"gridsquare"` // Maidenhead grid square
	Country         string    `json:"country"`
	DXCC            int       `json:"dxcc"` // DXCC entity number
	State           string    `json:"state"`
	County          string    `json:"county"`
	Comment         string    `json:"comment"`
	Notes           string    `json:"notes"`
	TXPower         float64   `json:"tx_pwr"` // Transmit power in watts
	RXPower         float64   `json:"rx_pwr"` // Receive power in watts
	PropagationMode string    `json:"prop_mode"` // e.g., "F2", "Es", "SAT"
	SatName         string    `json:"sat_name"`  // Satellite name if applicable
	SatMode         string    `json:"sat_mode"`  // Satellite mode
	Contest         string    `json:"contest_id"`
	STX             int       `json:"stx"` // Serial number transmitted
	SRX             int       `json:"srx"` // Serial number received
	LoTWQSLSent     string    `json:"lotw_qsl_sent"`   // "Y", "N", "R" (requested)
	LoTWQSLRcvd     string    `json:"lotw_qsl_rcvd"`   // "Y", "N"
	LoTWQSLRDate    time.Time `json:"lotw_qslrdate"`
	EQSLQSLSent     string    `json:"eqsl_qsl_sent"`
	EQSLQSLRcvd     string    `json:"eqsl_qsl_rcvd"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// QSOFilter represents filtering options for QSO queries
type QSOFilter struct {
	Callsign   string
	Band       string
	Mode       string
	Country    string
	DateFrom   time.Time
	DateTo     time.Time
	Limit      int
	Offset     int
}

// QSOStats represents QSO statistics for a user
type QSOStats struct {
	TotalQSOs     int `json:"total_qsos"`
	QSOLimit      int `json:"qso_limit"`
	RemainingQSOs int `json:"remaining_qsos"`
}
