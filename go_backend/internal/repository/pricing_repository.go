package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smartwaste/backend/internal/models"
)

// PricingRepository handles pricing rule data operations
type PricingRepository struct {
	db *sqlx.DB
}

// NewPricingRepository creates a new PricingRepository instance
func NewPricingRepository(db *sqlx.DB) *PricingRepository {
	return &PricingRepository{db: db}
}

// Create creates a new pricing rule
func (r *PricingRepository) Create(ctx context.Context, rule *models.PricingRule) error {
	query := `
		INSERT INTO pricing_rules (waste_type, condition, price_per_kg, currency, min_weight_kg, max_weight_kg, company_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, is_active, created_at, updated_at`

	return r.db.QueryRowxContext(ctx, query,
		rule.WasteType,
		rule.Condition,
		rule.PricePerKg,
		rule.Currency,
		rule.MinWeightKg,
		rule.MaxWeightKg,
		rule.CompanyID,
	).Scan(&rule.ID, &rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt)
}

// GetByID retrieves a pricing rule by ID
func (r *PricingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.PricingRule, error) {
	var rule models.PricingRule
	query := `SELECT * FROM pricing_rules WHERE id = $1`

	err := r.db.GetContext(ctx, &rule, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &rule, err
}

// GetByTypeAndCondition retrieves pricing rules by waste type and condition
func (r *PricingRepository) GetByTypeAndCondition(ctx context.Context, wasteType, condition string) (*models.PricingRule, error) {
	var rule models.PricingRule
	query := `SELECT * FROM pricing_rules WHERE waste_type = $1 AND condition = $2 AND is_active = true ORDER BY created_at DESC LIMIT 1`

	err := r.db.GetContext(ctx, &rule, query, wasteType, condition)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &rule, err
}

// Update updates a pricing rule
func (r *PricingRepository) Update(ctx context.Context, rule *models.PricingRule) error {
	query := `
		UPDATE pricing_rules
		SET waste_type = $1, condition = $2, price_per_kg = $3, currency = $4, min_weight_kg = $5, max_weight_kg = $6, is_active = $7
		WHERE id = $8
		RETURNING updated_at`

	return r.db.QueryRowxContext(ctx, query,
		rule.WasteType,
		rule.Condition,
		rule.PricePerKg,
		rule.Currency,
		rule.MinWeightKg,
		rule.MaxWeightKg,
		rule.IsActive,
		rule.ID,
	).Scan(&rule.UpdatedAt)
}

// List retrieves all pricing rules with pagination
func (r *PricingRepository) List(ctx context.Context, limit, offset int) ([]models.PricingRule, error) {
	var rules []models.PricingRule
	query := `SELECT * FROM pricing_rules WHERE is_active = true ORDER BY waste_type, condition LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &rules, query, limit, offset)
	return rules, err
}

// ListByCompany retrieves pricing rules for a specific company
func (r *PricingRepository) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.PricingRule, error) {
	var rules []models.PricingRule
	query := `SELECT * FROM pricing_rules WHERE company_id = $1 AND is_active = true ORDER BY waste_type, condition`
	err := r.db.SelectContext(ctx, &rules, query, companyID)
	return rules, err
}

// Delete deletes a pricing rule (soft delete)
func (r *PricingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE pricing_rules SET is_active = false WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
