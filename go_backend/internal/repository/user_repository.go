package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smartwaste/backend/internal/models"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, phone, address, reward_points)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowxContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Phone,
		user.Address,
		user.RewardPoints,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET full_name = $1, phone = $2, address = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at`

	return r.db.QueryRowxContext(ctx, query,
		user.FullName,
		user.Phone,
		user.Address,
		user.ID,
	).Scan(&user.UpdatedAt)
}

// UpdateRewardPoints updates a user's reward points
func (r *UserRepository) UpdateRewardPoints(ctx context.Context, id uuid.UUID, points int) error {
	query := `UPDATE users SET reward_points = reward_points + $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, points, id)
	return err
}

// GetRewardPoints retrieves a user's reward points
func (r *UserRepository) GetRewardPoints(ctx context.Context, id uuid.UUID) (int, error) {
	var points int
	query := `SELECT reward_points FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &points, query, id)
	return points, err
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves all users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	return users, err
}
