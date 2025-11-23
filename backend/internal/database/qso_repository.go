package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Nitefawkes/Skypig/backend/internal/models"
)

// QSORepository handles QSO database operations
type QSORepository struct {
	db *DB
}

// NewQSORepository creates a new QSO repository
func NewQSORepository(db *DB) *QSORepository {
	return &QSORepository{db: db}
}

// Create inserts a new QSO into the database
func (r *QSORepository) Create(ctx context.Context, qso *models.QSO) error {
	query := `
		INSERT INTO qsos (
			user_id, callsign, operator_call, station_callsign,
			qso_date, time_on, time_off,
			band, band_rx, freq, freq_rx,
			mode, submode, rst_sent, rst_rcvd,
			name, qth, gridsquare, country, dxcc, state, county,
			comment, notes, tx_pwr, rx_pwr,
			prop_mode, sat_name, sat_mode,
			contest_id, stx, srx
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32
		) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		qso.UserID, qso.Callsign, nullString(qso.OperatorCall), nullString(qso.StationCallsign),
		qso.QSODate, qso.TimeOn, nullTime(qso.TimeOff),
		nullString(qso.Band), nullString(qso.BandRX), nullFloat64(qso.Freq), nullFloat64(qso.FreqRX),
		nullString(qso.Mode), nullString(qso.Submode), nullString(qso.RSTSent), nullString(qso.RSTRcvd),
		nullString(qso.Name), nullString(qso.QTH), nullString(qso.GridSquare),
		nullString(qso.Country), nullInt(qso.DXCC), nullString(qso.State), nullString(qso.County),
		nullString(qso.Comment), nullString(qso.Notes), nullFloat64(qso.TXPower), nullFloat64(qso.RXPower),
		nullString(qso.PropagationMode), nullString(qso.SatName), nullString(qso.SatMode),
		nullString(qso.Contest), nullInt(qso.STX), nullInt(qso.SRX),
	).Scan(&qso.ID, &qso.CreatedAt, &qso.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create QSO: %w", err)
	}

	return nil
}

// GetByID retrieves a QSO by ID
func (r *QSORepository) GetByID(ctx context.Context, id, userID int64) (*models.QSO, error) {
	query := `
		SELECT id, user_id, callsign, operator_call, station_callsign,
			qso_date, time_on, time_off,
			band, band_rx, freq, freq_rx,
			mode, submode, rst_sent, rst_rcvd,
			name, qth, gridsquare, country, dxcc, state, county,
			comment, notes, tx_pwr, rx_pwr,
			prop_mode, sat_name, sat_mode,
			contest_id, stx, srx,
			lotw_qsl_sent, lotw_qsl_rcvd, lotw_qslrdate,
			eqsl_qsl_sent, eqsl_qsl_rcvd,
			created_at, updated_at
		FROM qsos
		WHERE id = $1 AND user_id = $2
	`

	qso := &models.QSO{}
	var operatorCall, stationCall, timeOff, bandRx, submode, rstSent, rstRcvd sql.NullString
	var freq, freqRx, txPwr, rxPwr sql.NullFloat64
	var name, qth, grid, country, state, county, comment, notes sql.NullString
	var band, mode, propMode, satName, satMode, contest sql.NullString
	var dxcc, stx, srx sql.NullInt64
	var lotwSent, lotwRcvd, lotwDate, eqslSent, eqslRcvd sql.NullString

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&qso.ID, &qso.UserID, &qso.Callsign, &operatorCall, &stationCall,
		&qso.QSODate, &qso.TimeOn, &timeOff,
		&band, &bandRx, &freq, &freqRx,
		&mode, &submode, &rstSent, &rstRcvd,
		&name, &qth, &grid, &country, &dxcc, &state, &county,
		&comment, &notes, &txPwr, &rxPwr,
		&propMode, &satName, &satMode,
		&contest, &stx, &srx,
		&lotwSent, &lotwRcvd, &lotwDate,
		&eqslSent, &eqslRcvd,
		&qso.CreatedAt, &qso.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("QSO not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get QSO: %w", err)
	}

	// Map nullable fields
	qso.OperatorCall = operatorCall.String
	qso.StationCallsign = stationCall.String
	if timeOff.Valid {
		t, _ := time.Parse(time.RFC3339, timeOff.String)
		qso.TimeOff = t
	}
	qso.Band = band.String
	qso.BandRX = bandRx.String
	if freq.Valid {
		qso.Freq = freq.Float64
	}
	if freqRx.Valid {
		qso.FreqRX = freqRx.Float64
	}
	qso.Mode = mode.String
	qso.Submode = submode.String
	qso.RSTSent = rstSent.String
	qso.RSTRcvd = rstRcvd.String
	qso.Name = name.String
	qso.QTH = qth.String
	qso.GridSquare = grid.String
	qso.Country = country.String
	if dxcc.Valid {
		qso.DXCC = int(dxcc.Int64)
	}
	qso.State = state.String
	qso.County = county.String
	qso.Comment = comment.String
	qso.Notes = notes.String
	if txPwr.Valid {
		qso.TXPower = txPwr.Float64
	}
	if rxPwr.Valid {
		qso.RXPower = rxPwr.Float64
	}
	qso.PropagationMode = propMode.String
	qso.SatName = satName.String
	qso.SatMode = satMode.String
	qso.Contest = contest.String
	if stx.Valid {
		qso.STX = int(stx.Int64)
	}
	if srx.Valid {
		qso.SRX = int(srx.Int64)
	}
	qso.LoTWQSLSent = lotwSent.String
	qso.LoTWQSLRcvd = lotwRcvd.String
	qso.EQSLQSLSent = eqslSent.String
	qso.EQSLQSLRcvd = eqslRcvd.String

	return qso, nil
}

// List retrieves QSOs with optional filtering
func (r *QSORepository) List(ctx context.Context, userID int64, filter models.QSOFilter) ([]*models.QSO, int, error) {
	// Build query dynamically based on filters
	query := `
		SELECT id, user_id, callsign, operator_call, station_callsign,
			qso_date, time_on, time_off,
			band, band_rx, freq, freq_rx,
			mode, submode, rst_sent, rst_rcvd,
			name, qth, gridsquare, country, dxcc, state, county,
			comment, notes, tx_pwr, rx_pwr,
			prop_mode, sat_name, sat_mode,
			contest_id, stx, srx,
			lotw_qsl_sent, lotw_qsl_rcvd, lotw_qslrdate,
			eqsl_qsl_sent, eqsl_qsl_rcvd,
			created_at, updated_at
		FROM qsos
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argPos := 2

	if filter.Callsign != "" {
		query += fmt.Sprintf(" AND callsign ILIKE $%d", argPos)
		args = append(args, "%"+filter.Callsign+"%")
		argPos++
	}

	if filter.Band != "" {
		query += fmt.Sprintf(" AND band = $%d", argPos)
		args = append(args, filter.Band)
		argPos++
	}

	if filter.Mode != "" {
		query += fmt.Sprintf(" AND mode = $%d", argPos)
		args = append(args, filter.Mode)
		argPos++
	}

	if filter.Country != "" {
		query += fmt.Sprintf(" AND country ILIKE $%d", argPos)
		args = append(args, "%"+filter.Country+"%")
		argPos++
	}

	if !filter.DateFrom.IsZero() {
		query += fmt.Sprintf(" AND qso_date >= $%d", argPos)
		args = append(args, filter.DateFrom)
		argPos++
	}

	if !filter.DateTo.IsZero() {
		query += fmt.Sprintf(" AND qso_date <= $%d", argPos)
		args = append(args, filter.DateTo)
		argPos++
	}

	// Count total matching records
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS count_query"
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count QSOs: %w", err)
	}

	// Add ordering and pagination
	query += " ORDER BY time_on DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filter.Limit)
		argPos++
	} else {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, 50) // Default limit
		argPos++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list QSOs: %w", err)
	}
	defer rows.Close()

	qsos := []*models.QSO{}
	for rows.Next() {
		qso := &models.QSO{}
		var operatorCall, stationCall, timeOff, bandRx, submode, rstSent, rstRcvd sql.NullString
		var freq, freqRx, txPwr, rxPwr sql.NullFloat64
		var name, qth, grid, country, state, county, comment, notes sql.NullString
		var band, mode, propMode, satName, satMode, contest sql.NullString
		var dxcc, stx, srx sql.NullInt64
		var lotwSent, lotwRcvd, lotwDate, eqslSent, eqslRcvd sql.NullString

		err := rows.Scan(
			&qso.ID, &qso.UserID, &qso.Callsign, &operatorCall, &stationCall,
			&qso.QSODate, &qso.TimeOn, &timeOff,
			&band, &bandRx, &freq, &freqRx,
			&mode, &submode, &rstSent, &rstRcvd,
			&name, &qth, &grid, &country, &dxcc, &state, &county,
			&comment, &notes, &txPwr, &rxPwr,
			&propMode, &satName, &satMode,
			&contest, &stx, &srx,
			&lotwSent, &lotwRcvd, &lotwDate,
			&eqslSent, &eqslRcvd,
			&qso.CreatedAt, &qso.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan QSO: %w", err)
		}

		// Map nullable fields
		qso.OperatorCall = operatorCall.String
		qso.StationCallsign = stationCall.String
		if timeOff.Valid {
			t, _ := time.Parse(time.RFC3339, timeOff.String)
			qso.TimeOff = t
		}
		qso.Band = band.String
		qso.BandRX = bandRx.String
		if freq.Valid {
			qso.Freq = freq.Float64
		}
		if freqRx.Valid {
			qso.FreqRX = freqRx.Float64
		}
		qso.Mode = mode.String
		qso.Submode = submode.String
		qso.RSTSent = rstSent.String
		qso.RSTRcvd = rstRcvd.String
		qso.Name = name.String
		qso.QTH = qth.String
		qso.GridSquare = grid.String
		qso.Country = country.String
		if dxcc.Valid {
			qso.DXCC = int(dxcc.Int64)
		}
		qso.State = state.String
		qso.County = county.String
		qso.Comment = comment.String
		qso.Notes = notes.String
		if txPwr.Valid {
			qso.TXPower = txPwr.Float64
		}
		if rxPwr.Valid {
			qso.RXPower = rxPwr.Float64
		}
		qso.PropagationMode = propMode.String
		qso.SatName = satName.String
		qso.SatMode = satMode.String
		qso.Contest = contest.String
		if stx.Valid {
			qso.STX = int(stx.Int64)
		}
		if srx.Valid {
			qso.SRX = int(srx.Int64)
		}
		qso.LoTWQSLSent = lotwSent.String
		qso.LoTWQSLRcvd = lotwRcvd.String
		qso.EQSLQSLSent = eqslSent.String
		qso.EQSLQSLRcvd = eqslRcvd.String

		qsos = append(qsos, qso)
	}

	return qsos, total, nil
}

// Update updates an existing QSO
func (r *QSORepository) Update(ctx context.Context, qso *models.QSO) error {
	query := `
		UPDATE qsos SET
			callsign = $1, operator_call = $2, station_callsign = $3,
			qso_date = $4, time_on = $5, time_off = $6,
			band = $7, band_rx = $8, freq = $9, freq_rx = $10,
			mode = $11, submode = $12, rst_sent = $13, rst_rcvd = $14,
			name = $15, qth = $16, gridsquare = $17, country = $18,
			dxcc = $19, state = $20, county = $21,
			comment = $22, notes = $23, tx_pwr = $24, rx_pwr = $25,
			prop_mode = $26, sat_name = $27, sat_mode = $28,
			contest_id = $29, stx = $30, srx = $31,
			updated_at = NOW()
		WHERE id = $32 AND user_id = $33
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		qso.Callsign, nullString(qso.OperatorCall), nullString(qso.StationCallsign),
		qso.QSODate, qso.TimeOn, nullTime(qso.TimeOff),
		nullString(qso.Band), nullString(qso.BandRX), nullFloat64(qso.Freq), nullFloat64(qso.FreqRX),
		nullString(qso.Mode), nullString(qso.Submode), nullString(qso.RSTSent), nullString(qso.RSTRcvd),
		nullString(qso.Name), nullString(qso.QTH), nullString(qso.GridSquare),
		nullString(qso.Country), nullInt(qso.DXCC), nullString(qso.State), nullString(qso.County),
		nullString(qso.Comment), nullString(qso.Notes), nullFloat64(qso.TXPower), nullFloat64(qso.RXPower),
		nullString(qso.PropagationMode), nullString(qso.SatName), nullString(qso.SatMode),
		nullString(qso.Contest), nullInt(qso.STX), nullInt(qso.SRX),
		qso.ID, qso.UserID,
	).Scan(&qso.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("QSO not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update QSO: %w", err)
	}

	return nil
}

// Delete removes a QSO
func (r *QSORepository) Delete(ctx context.Context, id, userID int64) error {
	query := `DELETE FROM qsos WHERE id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
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

// Helper functions for nullable fields

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func nullInt(i int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: i != 0}
}

func nullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: f != 0}
}

func nullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: !t.IsZero()}
}
