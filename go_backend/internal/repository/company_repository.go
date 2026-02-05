package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smartwaste/backend/internal/models"
)

// CompanyRepository handles company data operations
type CompanyRepository struct {
	db *sqlx.DB
}

// NewCompanyRepository creates a new CompanyRepository instance
func NewCompanyRepository(db *sqlx.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// Create creates a new company
func (r *CompanyRepository) Create(ctx context.Context, company *models.Company) error {
	query := `
		INSERT INTO companies (name, email, phone, address, city, country, registration_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, is_active, created_at, updated_at`

	return r.db.QueryRowxContext(ctx, query,
		company.Name,
		company.Email,
		company.Phone,
		company.Address,
		company.City,
		company.Country,
		company.RegistrationNumber,
	).Scan(&company.ID, &company.IsActive, &company.CreatedAt, &company.UpdatedAt)
}

// GetByID retrieves a company by ID
func (r *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	var company models.Company
	query := `SELECT * FROM companies WHERE id = $1`

	err := r.db.GetContext(ctx, &company, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &company, err
}

// GetByEmail retrieves a company by email
func (r *CompanyRepository) GetByEmail(ctx context.Context, email string) (*models.Company, error) {
	var company models.Company
	query := `SELECT * FROM companies WHERE email = $1`

	err := r.db.GetContext(ctx, &company, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &company, err
}

// Update updates a company
func (r *CompanyRepository) Update(ctx context.Context, company *models.Company) error {
	query := `
		UPDATE companies
		SET name = $1, email = $2, phone = $3, address = $4, city = $5, country = $6, registration_number = $7, is_active = $8
		WHERE id = $9
		RETURNING updated_at`

	return r.db.QueryRowxContext(ctx, query,
		company.Name,
		company.Email,
		company.Phone,
		company.Address,
		company.City,
		company.Country,
		company.RegistrationNumber,
		company.IsActive,
		company.ID,
	).Scan(&company.UpdatedAt)
}

// List retrieves all companies with pagination
func (r *CompanyRepository) List(ctx context.Context, limit, offset int) ([]models.Company, error) {
	var companies []models.Company
	query := `SELECT * FROM companies WHERE is_active = true ORDER BY name ASC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &companies, query, limit, offset)
	return companies, err
}

// Delete deletes a company (soft delete)
func (r *CompanyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE companies SET is_active = false WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
