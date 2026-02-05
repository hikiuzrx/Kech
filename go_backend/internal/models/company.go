package models

import (
	"time"

	"github.com/google/uuid"
)

// Company represents a recycling company
type Company struct {
	ID                 uuid.UUID `db:"id" json:"id"`
	Name               string    `db:"name" json:"name"`
	Email              string    `db:"email" json:"email"`
	Phone              *string   `db:"phone" json:"phone,omitempty"`
	Address            *string   `db:"address" json:"address,omitempty"`
	City               *string   `db:"city" json:"city,omitempty"`
	Country            *string   `db:"country" json:"country,omitempty"`
	RegistrationNumber *string   `db:"registration_number" json:"registration_number,omitempty"`
	IsActive           bool      `db:"is_active" json:"is_active"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

// CreateCompanyRequest represents the request to create a new company
type CreateCompanyRequest struct {
	Name               string  `json:"name" binding:"required"`
	Email              string  `json:"email" binding:"required,email"`
	Phone              *string `json:"phone"`
	Address            *string `json:"address"`
	City               *string `json:"city"`
	Country            *string `json:"country"`
	RegistrationNumber *string `json:"registration_number"`
}

// UpdateCompanyRequest represents the request to update a company
type UpdateCompanyRequest struct {
	Name               *string `json:"name"`
	Email              *string `json:"email"`
	Phone              *string `json:"phone"`
	Address            *string `json:"address"`
	City               *string `json:"city"`
	Country            *string `json:"country"`
	RegistrationNumber *string `json:"registration_number"`
	IsActive           *bool   `json:"is_active"`
}

// CompanyResponse represents the API response for a company
type CompanyResponse struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Phone              *string   `json:"phone,omitempty"`
	Address            *string   `json:"address,omitempty"`
	City               *string   `json:"city,omitempty"`
	Country            *string   `json:"country,omitempty"`
	RegistrationNumber *string   `json:"registration_number,omitempty"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ToResponse converts Company to CompanyResponse
func (c *Company) ToResponse() *CompanyResponse {
	return &CompanyResponse{
		ID:                 c.ID,
		Name:               c.Name,
		Email:              c.Email,
		Phone:              c.Phone,
		Address:            c.Address,
		City:               c.City,
		Country:            c.Country,
		RegistrationNumber: c.RegistrationNumber,
		IsActive:           c.IsActive,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}
