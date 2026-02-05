package models

import (
	"time"

	"github.com/google/uuid"
)

// PricingRule represents a pricing rule for waste valuation
type PricingRule struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	WasteType   string     `db:"waste_type" json:"waste_type"`
	Condition   string     `db:"condition" json:"condition"`
	PricePerKg  float64    `db:"price_per_kg" json:"price_per_kg"`
	Currency    string     `db:"currency" json:"currency"`
	MinWeightKg float64    `db:"min_weight_kg" json:"min_weight_kg"`
	MaxWeightKg *float64   `db:"max_weight_kg" json:"max_weight_kg,omitempty"`
	CompanyID   *uuid.UUID `db:"company_id" json:"company_id,omitempty"`
	IsActive    bool       `db:"is_active" json:"is_active"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// CreatePricingRuleRequest represents the request to create a pricing rule
type CreatePricingRuleRequest struct {
	WasteType   string     `json:"waste_type" binding:"required"`
	Condition   string     `json:"condition" binding:"required"`
	PricePerKg  float64    `json:"price_per_kg" binding:"required,gt=0"`
	Currency    string     `json:"currency" binding:"required,len=3"`
	MinWeightKg float64    `json:"min_weight_kg"`
	MaxWeightKg *float64   `json:"max_weight_kg"`
	CompanyID   *uuid.UUID `json:"company_id"`
}

// UpdatePricingRuleRequest represents the request to update a pricing rule
type UpdatePricingRuleRequest struct {
	WasteType   *string  `json:"waste_type"`
	Condition   *string  `json:"condition"`
	PricePerKg  *float64 `json:"price_per_kg"`
	Currency    *string  `json:"currency"`
	MinWeightKg *float64 `json:"min_weight_kg"`
	MaxWeightKg *float64 `json:"max_weight_kg"`
	IsActive    *bool    `json:"is_active"`
}

// PricingRuleResponse represents the API response for a pricing rule
type PricingRuleResponse struct {
	ID          uuid.UUID  `json:"id"`
	WasteType   string     `json:"waste_type"`
	Condition   string     `json:"condition"`
	PricePerKg  float64    `json:"price_per_kg"`
	Currency    string     `json:"currency"`
	MinWeightKg float64    `json:"min_weight_kg"`
	MaxWeightKg *float64   `json:"max_weight_kg,omitempty"`
	CompanyID   *uuid.UUID `json:"company_id,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToResponse converts PricingRule to PricingRuleResponse
func (p *PricingRule) ToResponse() *PricingRuleResponse {
	return &PricingRuleResponse{
		ID:          p.ID,
		WasteType:   p.WasteType,
		Condition:   p.Condition,
		PricePerKg:  p.PricePerKg,
		Currency:    p.Currency,
		MinWeightKg: p.MinWeightKg,
		MaxWeightKg: p.MaxWeightKg,
		CompanyID:   p.CompanyID,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
