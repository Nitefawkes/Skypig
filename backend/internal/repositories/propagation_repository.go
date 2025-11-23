package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nitefawkes/ham-radio-cloud/internal/models"
)

type PropagationRepository struct {
	db *sql.DB
}

func NewPropagationRepository(db *sql.DB) *PropagationRepository {
	return &PropagationRepository{db: db}
}

func (r *PropagationRepository) Create(data *models.PropagationData) error {
	query := `
		INSERT INTO propagation_data
		(timestamp, solar_flux, sunspot_number, a_index, k_index, xray_flux,
		 helium_line, proton_flux, electron_flux, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		query,
		data.Timestamp, data.SolarFlux, data.SunspotNumber, data.AIndex,
		data.KIndex, data.XRayFlux, data.HeliumLine, data.ProtonFlux,
		data.ElectronFlux, data.Source,
	).Scan(&data.ID, &data.CreatedAt)
}

func (r *PropagationRepository) GetLatest() (*models.PropagationData, error) {
	query := `
		SELECT id, timestamp, solar_flux, sunspot_number, a_index, k_index,
		       xray_flux, helium_line, proton_flux, electron_flux, source, created_at
		FROM propagation_data
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var data models.PropagationData
	err := r.db.QueryRow(query).Scan(
		&data.ID, &data.Timestamp, &data.SolarFlux, &data.SunspotNumber,
		&data.AIndex, &data.KIndex, &data.XRayFlux, &data.HeliumLine,
		&data.ProtonFlux, &data.ElectronFlux, &data.Source, &data.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no propagation data available")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest propagation data: %w", err)
	}

	return &data, nil
}

func (r *PropagationRepository) GetByTimeRange(start, end time.Time) ([]models.PropagationData, error) {
	query := `
		SELECT id, timestamp, solar_flux, sunspot_number, a_index, k_index,
		       xray_flux, helium_line, proton_flux, electron_flux, source, created_at
		FROM propagation_data
		WHERE timestamp BETWEEN $1 AND $2
		ORDER BY timestamp DESC
	`

	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query propagation data: %w", err)
	}
	defer rows.Close()

	var results []models.PropagationData
	for rows.Next() {
		var data models.PropagationData
		err := rows.Scan(
			&data.ID, &data.Timestamp, &data.SolarFlux, &data.SunspotNumber,
			&data.AIndex, &data.KIndex, &data.XRayFlux, &data.HeliumLine,
			&data.ProtonFlux, &data.ElectronFlux, &data.Source, &data.CreatedAt,
		)
		if err != nil {
			continue
		}
		results = append(results, data)
	}

	return results, nil
}

func (r *PropagationRepository) DeleteOlderThan(cutoff time.Time) error {
	query := `DELETE FROM propagation_data WHERE timestamp < $1`
	_, err := r.db.Exec(query, cutoff)
	return err
}
