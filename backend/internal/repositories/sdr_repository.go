package repositories

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nitefawkes/ham-radio-cloud/internal/models"
)

type SDRRepository struct {
	db *sqlx.DB
}

func NewSDRRepository(db *sqlx.DB) *SDRRepository {
	return &SDRRepository{db: db}
}

func (r *SDRRepository) Create(sdr *models.SDRReceiver) error {
	query := `
		INSERT INTO sdr_receivers (
			name, callsign, url, type, location, grid_square,
			latitude, longitude, country, bands, modes, antenna_info,
			frequency_min, frequency_max, users_max, status, last_seen,
			description, avatar_url, is_public
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		) RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		sdr.Name, sdr.Callsign, sdr.URL, sdr.Type, sdr.Location, sdr.GridSquare,
		sdr.Latitude, sdr.Longitude, sdr.Country, sdr.Bands, sdr.Modes, sdr.AntennaInfo,
		sdr.FrequencyMin, sdr.FrequencyMax, sdr.UsersMax, sdr.Status, sdr.LastSeen,
		sdr.Description, sdr.AvatarURL, sdr.IsPublic,
	).Scan(&sdr.ID, &sdr.CreatedAt, &sdr.UpdatedAt)
}

func (r *SDRRepository) Upsert(sdr *models.SDRReceiver) error {
	query := `
		INSERT INTO sdr_receivers (
			name, callsign, url, type, location, grid_square,
			latitude, longitude, country, bands, modes, antenna_info,
			frequency_min, frequency_max, users_max, status, last_seen,
			description, avatar_url, is_public
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)
		ON CONFLICT (url)
		DO UPDATE SET
			name = EXCLUDED.name,
			callsign = EXCLUDED.callsign,
			type = EXCLUDED.type,
			location = EXCLUDED.location,
			grid_square = EXCLUDED.grid_square,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			country = EXCLUDED.country,
			bands = EXCLUDED.bands,
			modes = EXCLUDED.modes,
			antenna_info = EXCLUDED.antenna_info,
			frequency_min = EXCLUDED.frequency_min,
			frequency_max = EXCLUDED.frequency_max,
			users_max = EXCLUDED.users_max,
			status = EXCLUDED.status,
			last_seen = EXCLUDED.last_seen,
			description = EXCLUDED.description,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		sdr.Name, sdr.Callsign, sdr.URL, sdr.Type, sdr.Location, sdr.GridSquare,
		sdr.Latitude, sdr.Longitude, sdr.Country, sdr.Bands, sdr.Modes, sdr.AntennaInfo,
		sdr.FrequencyMin, sdr.FrequencyMax, sdr.UsersMax, sdr.Status, sdr.LastSeen,
		sdr.Description, sdr.AvatarURL, sdr.IsPublic,
	).Scan(&sdr.ID, &sdr.CreatedAt, &sdr.UpdatedAt)
}

func (r *SDRRepository) List(filter *models.SDRFilter) ([]models.SDRReceiver, error) {
	query := `SELECT * FROM sdr_receivers WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	// Build dynamic query based on filters
	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argPos)
		args = append(args, filter.Type)
		argPos++
	}

	if filter.Country != "" {
		query += fmt.Sprintf(" AND country = $%d", argPos)
		args = append(args, filter.Country)
		argPos++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, filter.Status)
		argPos++
	}

	if len(filter.Bands) > 0 {
		query += fmt.Sprintf(" AND bands && $%d", argPos)
		args = append(args, pq.Array(filter.Bands))
		argPos++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR location ILIKE $%d OR callsign ILIKE $%d)", argPos, argPos, argPos)
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	// Default to public SDRs only
	query += " AND is_public = true"

	// Order by status (online first) and name
	query += " ORDER BY CASE WHEN status = 'online' THEN 0 ELSE 1 END, name ASC"

	// Pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filter.Limit)
		argPos++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, filter.Offset)
		argPos++
	}

	var sdrs []models.SDRReceiver
	if err := r.db.Select(&sdrs, query, args...); err != nil {
		return nil, err
	}

	return sdrs, nil
}

func (r *SDRRepository) GetByID(id uuid.UUID) (*models.SDRReceiver, error) {
	var sdr models.SDRReceiver
	query := `SELECT * FROM sdr_receivers WHERE id = $1`
	if err := r.db.Get(&sdr, query, id); err != nil {
		return nil, err
	}
	return &sdr, nil
}

func (r *SDRRepository) GetByURL(url string) (*models.SDRReceiver, error) {
	var sdr models.SDRReceiver
	query := `SELECT * FROM sdr_receivers WHERE url = $1`
	if err := r.db.Get(&sdr, query, url); err != nil {
		return nil, err
	}
	return &sdr, nil
}

func (r *SDRRepository) Update(sdr *models.SDRReceiver) error {
	query := `
		UPDATE sdr_receivers SET
			name = $1, callsign = $2, type = $3, location = $4, grid_square = $5,
			latitude = $6, longitude = $7, country = $8, bands = $9, modes = $10,
			antenna_info = $11, frequency_min = $12, frequency_max = $13, users_max = $14,
			status = $15, last_seen = $16, description = $17, avatar_url = $18, is_public = $19
		WHERE id = $20`

	_, err := r.db.Exec(
		query,
		sdr.Name, sdr.Callsign, sdr.Type, sdr.Location, sdr.GridSquare,
		sdr.Latitude, sdr.Longitude, sdr.Country, sdr.Bands, sdr.Modes,
		sdr.AntennaInfo, sdr.FrequencyMin, sdr.FrequencyMax, sdr.UsersMax,
		sdr.Status, sdr.LastSeen, sdr.Description, sdr.AvatarURL, sdr.IsPublic,
		sdr.ID,
	)

	return err
}

func (r *SDRRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM sdr_receivers WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SDRRepository) Count(filter *models.SDRFilter) (int, error) {
	query := `SELECT COUNT(*) FROM sdr_receivers WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argPos)
		args = append(args, filter.Type)
		argPos++
	}

	if filter.Country != "" {
		query += fmt.Sprintf(" AND country = $%d", argPos)
		args = append(args, filter.Country)
		argPos++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, filter.Status)
		argPos++
	}

	if len(filter.Bands) > 0 {
		query += fmt.Sprintf(" AND bands && $%d", argPos)
		args = append(args, pq.Array(filter.Bands))
		argPos++
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query += fmt.Sprintf(" AND (name ILIKE $%d OR location ILIKE $%d OR callsign ILIKE $%d)", argPos, argPos+1, argPos+2)
		args = append(args, searchPattern, searchPattern, searchPattern)
		argPos += 3
	}

	query += " AND is_public = true"

	var count int
	err := r.db.Get(&count, query, args...)
	return count, err
}

func (r *SDRRepository) AddFavorite(userID, sdrID uuid.UUID) error {
	query := `INSERT INTO user_favorite_sdrs (user_id, sdr_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, userID, sdrID)
	return err
}

func (r *SDRRepository) RemoveFavorite(userID, sdrID uuid.UUID) error {
	query := `DELETE FROM user_favorite_sdrs WHERE user_id = $1 AND sdr_id = $2`
	_, err := r.db.Exec(query, userID, sdrID)
	return err
}

func (r *SDRRepository) GetUserFavorites(userID uuid.UUID) ([]models.SDRReceiver, error) {
	query := `
		SELECT s.* FROM sdr_receivers s
		JOIN user_favorite_sdrs f ON s.id = f.sdr_id
		WHERE f.user_id = $1
		ORDER BY s.name ASC`

	var sdrs []models.SDRReceiver
	if err := r.db.Select(&sdrs, query, userID); err != nil {
		return nil, err
	}

	return sdrs, nil
}

// BulkUpsert efficiently inserts or updates multiple SDRs
func (r *SDRRepository) BulkUpsert(sdrs []models.SDRReceiver) error {
	if len(sdrs) == 0 {
		return nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(`
		INSERT INTO sdr_receivers (
			name, callsign, url, type, location, grid_square,
			latitude, longitude, country, bands, modes, antenna_info,
			frequency_min, frequency_max, users_max, status, last_seen,
			description, avatar_url, is_public
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)
		ON CONFLICT (url)
		DO UPDATE SET
			name = EXCLUDED.name,
			callsign = EXCLUDED.callsign,
			type = EXCLUDED.type,
			location = EXCLUDED.location,
			grid_square = EXCLUDED.grid_square,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			country = EXCLUDED.country,
			bands = EXCLUDED.bands,
			modes = EXCLUDED.modes,
			antenna_info = EXCLUDED.antenna_info,
			frequency_min = EXCLUDED.frequency_min,
			frequency_max = EXCLUDED.frequency_max,
			users_max = EXCLUDED.users_max,
			status = EXCLUDED.status,
			last_seen = EXCLUDED.last_seen,
			description = EXCLUDED.description,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, sdr := range sdrs {
		_, err := stmt.Exec(
			sdr.Name, sdr.Callsign, sdr.URL, sdr.Type, sdr.Location, sdr.GridSquare,
			sdr.Latitude, sdr.Longitude, sdr.Country, sdr.Bands, sdr.Modes, sdr.AntennaInfo,
			sdr.FrequencyMin, sdr.FrequencyMax, sdr.UsersMax, sdr.Status, sdr.LastSeen,
			sdr.Description, sdr.AvatarURL, sdr.IsPublic,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
