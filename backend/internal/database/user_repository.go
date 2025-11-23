package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Nitefawkes/Skypig/backend/internal/models"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, callsign, email, name, grid_square, qrz_verified,
			tier, stripe_id, qso_limit, qso_count,
			created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	var name, gridSquare, stripeID sql.NullString
	var lastLogin sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Callsign, &user.Email, &name, &gridSquare, &user.QRZVerified,
		&user.Tier, &stripeID, &user.QSOLimit, &user.QSOCount,
		&user.CreatedAt, &user.UpdatedAt, &lastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Name = name.String
	user.GridSquare = gridSquare.String
	user.StripeID = stripeID.String
	if lastLogin.Valid {
		user.LastLoginAt = lastLogin.Time
	}

	return user, nil
}

// GetByCallsign retrieves a user by callsign
func (r *UserRepository) GetByCallsign(ctx context.Context, callsign string) (*models.User, error) {
	query := `
		SELECT id, callsign, email, name, grid_square, qrz_verified,
			tier, stripe_id, qso_limit, qso_count,
			created_at, updated_at, last_login_at
		FROM users
		WHERE callsign = $1
	`

	user := &models.User{}
	var name, gridSquare, stripeID sql.NullString
	var lastLogin sql.NullTime

	err := r.db.QueryRowContext(ctx, query, callsign).Scan(
		&user.ID, &user.Callsign, &user.Email, &name, &gridSquare, &user.QRZVerified,
		&user.Tier, &stripeID, &user.QSOLimit, &user.QSOCount,
		&user.CreatedAt, &user.UpdatedAt, &lastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Name = name.String
	user.GridSquare = gridSquare.String
	user.StripeID = stripeID.String
	if lastLogin.Valid {
		user.LastLoginAt = lastLogin.Time
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, callsign, email, name, grid_square, qrz_verified,
			tier, stripe_id, qso_limit, qso_count,
			created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	var name, gridSquare, stripeID sql.NullString
	var lastLogin sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Callsign, &user.Email, &name, &gridSquare, &user.QRZVerified,
		&user.Tier, &stripeID, &user.QSOLimit, &user.QSOCount,
		&user.CreatedAt, &user.UpdatedAt, &lastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Name = name.String
	user.GridSquare = gridSquare.String
	user.StripeID = stripeID.String
	if lastLogin.Valid {
		user.LastLoginAt = lastLogin.Time
	}

	return user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (callsign, email, name, grid_square, qrz_verified, tier, qso_limit)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at, qso_count
	`

	err := r.db.QueryRowContext(ctx, query,
		user.Callsign, user.Email, nullString(user.Name),
		nullString(user.GridSquare), user.QRZVerified,
		user.Tier, user.GetQSOLimit(),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.QSOCount)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.QSOLimit = user.GetQSOLimit()

	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET
			callsign = $1, email = $2, name = $3, grid_square = $4,
			qrz_verified = $5, tier = $6, stripe_id = $7,
			updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		user.Callsign, user.Email, nullString(user.Name),
		nullString(user.GridSquare), user.QRZVerified,
		user.Tier, nullString(user.StripeID), user.ID,
	).Scan(&user.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("user not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}
