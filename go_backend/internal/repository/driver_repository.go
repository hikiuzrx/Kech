package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smartwaste/backend/internal/models"
)

// DriverRepository handles driver data operations
type DriverRepository struct {
	db *sqlx.DB
}

// NewDriverRepository creates a new DriverRepository instance
func NewDriverRepository(db *sqlx.DB) *DriverRepository {
	return &DriverRepository{db: db}
}

// Create creates a new driver
func (r *DriverRepository) Create(ctx context.Context, driver *models.Driver) error {
	query := `
		INSERT INTO drivers (email, password_hash, full_name, phone, license_number, vehicle_type, vehicle_plate)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowxContext(ctx, query,
		driver.Email,
		driver.PasswordHash,
		driver.FullName,
		driver.Phone,
		driver.LicenseNumber,
		driver.VehicleType,
		driver.VehiclePlate,
	).Scan(&driver.ID, &driver.CreatedAt, &driver.UpdatedAt)
}

// GetByID retrieves a driver by ID
func (r *DriverRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Driver, error) {
	var driver models.Driver
	query := `SELECT * FROM drivers WHERE id = $1`

	err := r.db.GetContext(ctx, &driver, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &driver, err
}

// GetByEmail retrieves a driver by email
func (r *DriverRepository) GetByEmail(ctx context.Context, email string) (*models.Driver, error) {
	var driver models.Driver
	query := `SELECT * FROM drivers WHERE email = $1`

	err := r.db.GetContext(ctx, &driver, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &driver, err
}

// Update updates a driver
func (r *DriverRepository) Update(ctx context.Context, driver *models.Driver) error {
	query := `
		UPDATE drivers
		SET full_name = $1, phone = $2, vehicle_type = $3, vehicle_plate = $4, is_available = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING updated_at`

	return r.db.QueryRowxContext(ctx, query,
		driver.FullName,
		driver.Phone,
		driver.VehicleType,
		driver.VehiclePlate,
		driver.IsAvailable,
		driver.ID,
	).Scan(&driver.UpdatedAt)
}

// UpdateLocation updates a driver's location
func (r *DriverRepository) UpdateLocation(ctx context.Context, id uuid.UUID, lat, lng float64) error {
	query := `UPDATE drivers SET latitude = $1, longitude = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, lat, lng, id)
	return err
}

// UpdateFCMToken updates a driver's FCM token
func (r *DriverRepository) UpdateFCMToken(ctx context.Context, id uuid.UUID, token string) error {
	query := `UPDATE drivers SET fcm_token = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, token, id)
	return err
}

// IncrementCollections increments a driver's total collections
func (r *DriverRepository) IncrementCollections(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE drivers SET total_collections = total_collections + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetAvailableDrivers retrieves all available drivers
func (r *DriverRepository) GetAvailableDrivers(ctx context.Context) ([]models.Driver, error) {
	var drivers []models.Driver
	query := `SELECT * FROM drivers WHERE is_available = true ORDER BY average_rating DESC`
	err := r.db.SelectContext(ctx, &drivers, query)
	return drivers, err
}

// GetNearestDriver finds the nearest available driver to a given location
func (r *DriverRepository) GetNearestDriver(ctx context.Context, lat, lng float64) (*models.Driver, error) {
	var driver models.Driver
	// Using Haversine formula approximation for distance calculation
	query := `
		SELECT *,
			(6371 * acos(cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + sin(radians($1)) * sin(radians(latitude)))) AS distance
		FROM drivers
		WHERE is_available = true AND latitude IS NOT NULL AND longitude IS NOT NULL
		ORDER BY distance ASC
		LIMIT 1`

	err := r.db.GetContext(ctx, &driver, query, lat, lng)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &driver, err
}

// List retrieves all drivers with pagination
func (r *DriverRepository) List(ctx context.Context, limit, offset int) ([]models.Driver, error) {
	var drivers []models.Driver
	query := `SELECT * FROM drivers ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &drivers, query, limit, offset)
	return drivers, err
}

// Delete deletes a driver
func (r *DriverRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM drivers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
