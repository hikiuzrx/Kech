package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FullName     string    `db:"full_name" json:"full_name"`
	Phone        *string   `db:"phone" json:"phone,omitempty"`
	Address      *string   `db:"address" json:"address,omitempty"`
	RewardPoints int       `db:"reward_points" json:"reward_points"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=8"`
	FullName string  `json:"full_name" binding:"required"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FullName *string `json:"full_name"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
}

// UserResponse represents the API response for a user
type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Phone        *string   `json:"phone,omitempty"`
	Address      *string   `json:"address,omitempty"`
	RewardPoints int       `json:"reward_points"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AddRewardPointsRequest represents the request to add reward points
type AddRewardPointsRequest struct {
	Points int    `json:"points" binding:"required,gt=0"`
	Reason string `json:"reason" binding:"required"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		Email:        u.Email,
		FullName:     u.FullName,
		Phone:        u.Phone,
		Address:      u.Address,
		RewardPoints: u.RewardPoints,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
