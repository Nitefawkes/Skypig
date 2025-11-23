package repositories

import (
	"database/sql"
	"fmt"
	"strings"

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
	`

	args := []interface{}{userID}
	argCount := 1

	// Apply filters
	if filter != nil {
		if filter.Callsign != "" {
			argCount++
			query += fmt.Sprintf(" AND UPPER(callsign) LIKE UPPER($%d)", argCount)
			args = append(args, "%"+filter.Callsign+"%")
		}
		if filter.Band != "" {
			argCount++
			query += fmt.Sprintf(" AND band = $%d", argCount)
			args = append(args, filter.Band)
		}
		if filter.Mode != "" {
			argCount++
			query += fmt.Sprintf(" AND mode = $%d", argCount)
			args = append(args, filter.Mode)
		}
		if filter.StartDate != nil {
			argCount++
			query += fmt.Sprintf(" AND qso_date >= $%d", argCount)
			args = append(args, filter.StartDate)
		}
		if filter.EndDate != nil {
			argCount++
			query += fmt.Sprintf(" AND qso_date <= $%d", argCount)
			args = append(args, filter.EndDate)
		}
	}

	query += " ORDER BY time_on DESC"

	// Apply pagination
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

	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, offset)

	rows, err := r.db.Query(query, args...)
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

	if qsos == nil {
		qsos = []models.QSO{}
	}

	return qsos, nil
}

func (r *QSORepository) GetByID(id, userID string) (*models.QSO, error) {
	query := `
		SELECT id, user_id, callsign, frequency, band, mode, rst_sent, rst_received,
			   qso_date, time_on, time_off, grid_square, country, state, county,
			   comment, contest_id, propagation_mode, satellite_name, tx_power,
			   lotw_sent, lotw_confirmed, created_at, updated_at
		FROM qsos
		WHERE id = $1 AND user_id = $2
	`

	var qso models.QSO
	err := r.db.QueryRow(query, id, userID).Scan(
		&qso.ID, &qso.UserID, &qso.Callsign, &qso.Frequency, &qso.Band, &qso.Mode,
		&qso.RST_Sent, &qso.RST_Received, &qso.QSODate, &qso.TimeOn, &qso.TimeOff,
		&qso.GridSquare, &qso.Country, &qso.State, &qso.County, &qso.Comment,
		&qso.ContestID, &qso.PropagationMode, &qso.SatelliteName, &qso.TXPower,
		&qso.LoTWSent, &qso.LoTWConfirmed, &qso.CreatedAt, &qso.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("QSO not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get QSO: %w", err)
	}

	return &qso, nil
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
		qso.UserID, strings.ToUpper(qso.Callsign), qso.Frequency, qso.Band, qso.Mode,
		qso.RST_Sent, qso.RST_Received, qso.QSODate, qso.TimeOn, qso.TimeOff,
		qso.GridSquare, qso.Country, qso.State, qso.County, qso.Comment, qso.TXPower,
	).Scan(&qso.ID, &qso.CreatedAt, &qso.UpdatedAt)
}

func (r *QSORepository) Update(qso *models.QSO) error {
	query := `
		UPDATE qsos
		SET callsign = $1, frequency = $2, band = $3, mode = $4, rst_sent = $5,
		    rst_received = $6, qso_date = $7, time_on = $8, time_off = $9,
		    grid_square = $10, country = $11, state = $12, county = $13,
		    comment = $14, tx_power = $15, updated_at = CURRENT_TIMESTAMP
		WHERE id = $16 AND user_id = $17
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		strings.ToUpper(qso.Callsign), qso.Frequency, qso.Band, qso.Mode,
		qso.RST_Sent, qso.RST_Received, qso.QSODate, qso.TimeOn, qso.TimeOff,
		qso.GridSquare, qso.Country, qso.State, qso.County, qso.Comment,
		qso.TXPower, qso.ID, qso.UserID,
	).Scan(&qso.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("QSO not found or unauthorized")
	}

	return err
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
		return fmt.Errorf("QSO not found or unauthorized")
	}

	return nil
}

func (r *QSORepository) GetCount(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM qsos WHERE user_id = $1`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *QSORepository) BulkCreate(qsos []models.QSO) (int, error) {
	if len(qsos) == 0 {
		return 0, nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO qsos (user_id, callsign, frequency, band, mode, rst_sent, rst_received,
						  qso_date, time_on, time_off, grid_square, country, state, county,
						  comment, tx_power)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	imported := 0
	for _, qso := range qsos {
		_, err := stmt.Exec(
			qso.UserID, strings.ToUpper(qso.Callsign), qso.Frequency, qso.Band, qso.Mode,
			qso.RST_Sent, qso.RST_Received, qso.QSODate, qso.TimeOn, qso.TimeOff,
			qso.GridSquare, qso.Country, qso.State, qso.County, qso.Comment, qso.TXPower,
		)
		if err != nil {
			// Log error but continue with other QSOs
			continue
		}
		imported++
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return imported, nil
}
