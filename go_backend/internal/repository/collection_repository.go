package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smartwaste/backend/internal/models"
)

// CollectionRepository handles collection data operations
type CollectionRepository struct {
	db *sqlx.DB
}

// NewCollectionRepository creates a new CollectionRepository instance
func NewCollectionRepository(db *sqlx.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

// Create creates a new collection
func (r *CollectionRepository) Create(ctx context.Context, collection *models.Collection) error {
	query := `
		INSERT INTO collections (bin_id, driver_id, fill_level_before, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, started_at`

	return r.db.QueryRowxContext(ctx, query,
		collection.BinID,
		collection.DriverID,
		collection.FillLevelBefore,
		models.CollectionStatusPending,
	).Scan(&collection.ID, &collection.StartedAt)
}

// GetByID retrieves a collection by ID
func (r *CollectionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	var collection models.Collection
	query := `SELECT * FROM collections WHERE id = $1`

	err := r.db.GetContext(ctx, &collection, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &collection, err
}

// Update updates a collection
func (r *CollectionRepository) Update(ctx context.Context, collection *models.Collection) error {
	query := `
		UPDATE collections
		SET fill_level_after = $1, weight_kg = $2, qr_code_verified = $3, notes = $4, status = $5, completed_at = $6
		WHERE id = $7`

	_, err := r.db.ExecContext(ctx, query,
		collection.FillLevelAfter,
		collection.WeightKg,
		collection.QRCodeVerified,
		collection.Notes,
		collection.Status,
		collection.CompletedAt,
		collection.ID,
	)
	return err
}

// Complete marks a collection as completed
func (r *CollectionRepository) Complete(ctx context.Context, id uuid.UUID, fillLevelAfter int, weightKg *float64, notes *string) error {
	now := time.Now()
	query := `
		UPDATE collections
		SET fill_level_after = $1, weight_kg = $2, notes = $3, status = $4, completed_at = $5
		WHERE id = $6`

	_, err := r.db.ExecContext(ctx, query, fillLevelAfter, weightKg, notes, models.CollectionStatusCompleted, now, id)
	return err
}

// VerifyQRCode verifies QR code for a collection
func (r *CollectionRepository) VerifyQRCode(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE collections SET qr_code_verified = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves all collections with pagination
func (r *CollectionRepository) List(ctx context.Context, limit, offset int) ([]models.Collection, error) {
	var collections []models.Collection
	query := `SELECT * FROM collections ORDER BY started_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &collections, query, limit, offset)
	return collections, err
}

// ListByDriver retrieves collections for a specific driver
func (r *CollectionRepository) ListByDriver(ctx context.Context, driverID uuid.UUID, limit, offset int) ([]models.Collection, error) {
	var collections []models.Collection
	query := `SELECT * FROM collections WHERE driver_id = $1 ORDER BY started_at DESC LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &collections, query, driverID, limit, offset)
	return collections, err
}

// ListByBin retrieves collections for a specific bin
func (r *CollectionRepository) ListByBin(ctx context.Context, binID uuid.UUID, limit, offset int) ([]models.Collection, error) {
	var collections []models.Collection
	query := `SELECT * FROM collections WHERE bin_id = $1 ORDER BY started_at DESC LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &collections, query, binID, limit, offset)
	return collections, err
}

// GetDriverStats retrieves driver performance statistics
func (r *CollectionRepository) GetDriverStats(ctx context.Context, driverID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total collections
	var total int
	err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM collections WHERE driver_id = $1`, driverID)
	if err != nil {
		return nil, err
	}
	stats["total_collections"] = total

	// Completed collections
	var completed int
	err = r.db.GetContext(ctx, &completed, `SELECT COUNT(*) FROM collections WHERE driver_id = $1 AND status = 'completed'`, driverID)
	if err != nil {
		return nil, err
	}
	stats["completed_collections"] = completed

	// Total weight collected
	var totalWeight sql.NullFloat64
	err = r.db.GetContext(ctx, &totalWeight, `SELECT COALESCE(SUM(weight_kg), 0) FROM collections WHERE driver_id = $1 AND status = 'completed'`, driverID)
	if err != nil {
		return nil, err
	}
	stats["total_weight_kg"] = totalWeight.Float64

	return stats, nil
}

// GetCollectionStats retrieves overall collection statistics
func (r *CollectionRepository) GetCollectionStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total collections today
	var todayCollections int
	err := r.db.GetContext(ctx, &todayCollections,
		`SELECT COUNT(*) FROM collections WHERE DATE(started_at) = CURRENT_DATE`)
	if err != nil {
		return nil, err
	}
	stats["today_collections"] = todayCollections

	// Total weight today
	var todayWeight sql.NullFloat64
	err = r.db.GetContext(ctx, &todayWeight,
		`SELECT COALESCE(SUM(weight_kg), 0) FROM collections WHERE DATE(started_at) = CURRENT_DATE AND status = 'completed'`)
	if err != nil {
		return nil, err
	}
	stats["today_weight_kg"] = todayWeight.Float64

	// Total collections this month
	var monthCollections int
	err = r.db.GetContext(ctx, &monthCollections,
		`SELECT COUNT(*) FROM collections WHERE DATE_TRUNC('month', started_at) = DATE_TRUNC('month', CURRENT_DATE)`)
	if err != nil {
		return nil, err
	}
	stats["month_collections"] = monthCollections

	return stats, nil
}
