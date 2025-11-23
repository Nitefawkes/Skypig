package repositories

import (
	"database/sql"
	"fmt"

	"github.com/nitefawkes/ham-radio-cloud/internal/models"
)

type QSORepository struct {
	db *sql.DB
}

func NewQSORepository(db *sql.DB) *QSORepository {
	return &QSORepository{db: db}
}

func (r *QSORepository) GetByUserID(userID string, filter *models.QSOFilter) ([]models.QSO, error) {
	query := `
		SELECT id, user_id, callsign, frequency, band, mode, rst_sent, rst_received,
			   qso_date, time_on, time_off, grid_square, country, state, county,
			   comment, contest_id, propagation_mode, satellite_name, tx_power,
			   lotw_sent, lotw_confirmed, created_at, updated_at
		FROM qsos
		WHERE user_id = $1
		ORDER BY time_on DESC
		LIMIT $2 OFFSET $3
	`

	limit := 100
	offset := 0
	if filter != nil {
		if filter.Limit > 0 {
			limit = filter.Limit
		}
		if filter.Offset > 0 {
			offset = filter.Offset
		}
	}

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query QSOs: %w", err)
	}
	defer rows.Close()

	var qsos []models.QSO
	for rows.Next() {
		var qso models.QSO
		err := rows.Scan(
			&qso.ID, &qso.UserID, &qso.Callsign, &qso.Frequency, &qso.Band, &qso.Mode,
			&qso.RST_Sent, &qso.RST_Received, &qso.QSODate, &qso.TimeOn, &qso.TimeOff,
			&qso.GridSquare, &qso.Country, &qso.State, &qso.County, &qso.Comment,
			&qso.ContestID, &qso.PropagationMode, &qso.SatelliteName, &qso.TXPower,
			&qso.LoTWSent, &qso.LoTWConfirmed, &qso.CreatedAt, &qso.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan QSO: %w", err)
		}
		qsos = append(qsos, qso)
	}

	return qsos, nil
}

func (r *QSORepository) Create(qso *models.QSO) error {
	query := `
		INSERT INTO qsos (user_id, callsign, frequency, band, mode, rst_sent, rst_received,
						  qso_date, time_on, time_off, grid_square, country, state, county,
						  comment, tx_power)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		qso.UserID, qso.Callsign, qso.Frequency, qso.Band, qso.Mode,
		qso.RST_Sent, qso.RST_Received, qso.QSODate, qso.TimeOn, qso.TimeOff,
		qso.GridSquare, qso.Country, qso.State, qso.County, qso.Comment, qso.TXPower,
	).Scan(&qso.ID, &qso.CreatedAt, &qso.UpdatedAt)
}

func (r *QSORepository) Delete(id, userID string) error {
	query := `DELETE FROM qsos WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete QSO: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("QSO not found")
	}

	return nil
}
