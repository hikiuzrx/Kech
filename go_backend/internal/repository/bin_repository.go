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

// BinRepository handles bin data operations
type BinRepository struct {
	db *sqlx.DB
}

// NewBinRepository creates a new BinRepository instance
func NewBinRepository(db *sqlx.DB) *BinRepository {
	return &BinRepository{db: db}
}

// Create creates a new bin
func (r *BinRepository) Create(ctx context.Context, bin *models.Bin) error {
	query := `
		INSERT INTO bins (device_id, location_name, latitude, longitude, waste_type, capacity_liters, company_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, fill_level, last_updated_at, is_active, created_at`

	return r.db.QueryRowxContext(ctx, query,
		bin.DeviceID,
		bin.LocationName,
		bin.Latitude,
		bin.Longitude,
		bin.WasteType,
		bin.CapacityLiters,
		bin.CompanyID,
	).Scan(&bin.ID, &bin.FillLevel, &bin.LastUpdatedAt, &bin.IsActive, &bin.CreatedAt)
}

// GetByID retrieves a bin by ID
func (r *BinRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Bin, error) {
	var bin models.Bin
	query := `SELECT * FROM bins WHERE id = $1`

	err := r.db.GetContext(ctx, &bin, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &bin, err
}

// GetByDeviceID retrieves a bin by device ID
func (r *BinRepository) GetByDeviceID(ctx context.Context, deviceID string) (*models.Bin, error) {
	var bin models.Bin
	query := `SELECT * FROM bins WHERE device_id = $1`

	err := r.db.GetContext(ctx, &bin, query, deviceID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &bin, err
}

// Update updates a bin
func (r *BinRepository) Update(ctx context.Context, bin *models.Bin) error {
	query := `
		UPDATE bins
		SET location_name = $1, latitude = $2, longitude = $3, waste_type = $4, capacity_liters = $5, is_active = $6, company_id = $7
		WHERE id = $8`

	_, err := r.db.ExecContext(ctx, query,
		bin.LocationName,
		bin.Latitude,
		bin.Longitude,
		bin.WasteType,
		bin.CapacityLiters,
		bin.IsActive,
		bin.CompanyID,
		bin.ID,
	)
	return err
}

// UpdateFillLevel updates a bin's fill level
func (r *BinRepository) UpdateFillLevel(ctx context.Context, deviceID string, fillLevel int) error {
	query := `UPDATE bins SET fill_level = $1, last_updated_at = CURRENT_TIMESTAMP WHERE device_id = $2`
	_, err := r.db.ExecContext(ctx, query, fillLevel, deviceID)
	return err
}

// MarkCollected marks a bin as collected
func (r *BinRepository) MarkCollected(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE bins SET fill_level = 0, last_collection_at = $1, last_updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// GetBinsNeedingCollection retrieves bins with fill level above threshold
func (r *BinRepository) GetBinsNeedingCollection(ctx context.Context, threshold int) ([]models.Bin, error) {
	var bins []models.Bin
	query := `SELECT * FROM bins WHERE is_active = true AND fill_level >= $1 ORDER BY fill_level DESC`
	err := r.db.SelectContext(ctx, &bins, query, threshold)
	return bins, err
}

// List retrieves all bins with pagination
func (r *BinRepository) List(ctx context.Context, limit, offset int) ([]models.Bin, error) {
	var bins []models.Bin
	query := `SELECT * FROM bins WHERE is_active = true ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &bins, query, limit, offset)
	return bins, err
}

// ListByCompany retrieves bins for a specific company
func (r *BinRepository) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Bin, error) {
	var bins []models.Bin
	query := `SELECT * FROM bins WHERE company_id = $1 AND is_active = true ORDER BY fill_level DESC`
	err := r.db.SelectContext(ctx, &bins, query, companyID)
	return bins, err
}

// Delete deletes a bin (soft delete by setting is_active = false)
func (r *BinRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE bins SET is_active = false WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetStatistics retrieves bin statistics
func (r *BinRepository) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total bins
	var totalBins int
	err := r.db.GetContext(ctx, &totalBins, `SELECT COUNT(*) FROM bins WHERE is_active = true`)
	if err != nil {
		return nil, err
	}
	stats["total_bins"] = totalBins

	// Bins needing collection (>80%)
	var needsCollection int
	err = r.db.GetContext(ctx, &needsCollection, `SELECT COUNT(*) FROM bins WHERE is_active = true AND fill_level >= 80`)
	if err != nil {
		return nil, err
	}
	stats["needs_collection"] = needsCollection

	// Average fill level
	var avgFillLevel float64
	err = r.db.GetContext(ctx, &avgFillLevel, `SELECT COALESCE(AVG(fill_level), 0) FROM bins WHERE is_active = true`)
	if err != nil {
		return nil, err
	}
	stats["average_fill_level"] = avgFillLevel

	return stats, nil
}
