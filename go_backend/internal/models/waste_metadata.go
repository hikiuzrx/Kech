package models

import (
	"time"

	"github.com/google/uuid"
)

// WasteMetadata represents AI-detected waste classification data
type WasteMetadata struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	CollectionID    *uuid.UUID `db:"collection_id" json:"collection_id,omitempty"`
	WasteType       string     `db:"waste_type" json:"waste_type"`
	Condition       string     `db:"condition" json:"condition"`
	ConfidenceScore *float64   `db:"confidence_score" json:"confidence_score,omitempty"`
	ImageURL        *string    `db:"image_url" json:"image_url,omitempty"`
	DetectedAt      time.Time  `db:"detected_at" json:"detected_at"`
	ValuatedPrice   *float64   `db:"valuated_price" json:"valuated_price,omitempty"`
	PricingRuleID   *uuid.UUID `db:"pricing_rule_id" json:"pricing_rule_id,omitempty"`
}

// CreateWasteMetadataRequest represents the request to create waste metadata
type CreateWasteMetadataRequest struct {
	CollectionID    *uuid.UUID `json:"collection_id"`
	WasteType       string     `json:"waste_type" binding:"required"`
	Condition       string     `json:"condition" binding:"required"`
	ConfidenceScore *float64   `json:"confidence_score"`
	ImageURL        *string    `json:"image_url"`
}

// WasteMetadataResponse represents the API response for waste metadata
type WasteMetadataResponse struct {
	ID              uuid.UUID  `json:"id"`
	CollectionID    *uuid.UUID `json:"collection_id,omitempty"`
	WasteType       string     `json:"waste_type"`
	Condition       string     `json:"condition"`
	ConfidenceScore *float64   `json:"confidence_score,omitempty"`
	ImageURL        *string    `json:"image_url,omitempty"`
	DetectedAt      time.Time  `json:"detected_at"`
	ValuatedPrice   *float64   `json:"valuated_price,omitempty"`
	PricingRuleID   *uuid.UUID `json:"pricing_rule_id,omitempty"`
}

// ValuationRequest represents the request to valuate waste
type ValuationRequest struct {
	WasteType string  `json:"waste_type" binding:"required"`
	Condition string  `json:"condition" binding:"required"`
	WeightKg  float64 `json:"weight_kg" binding:"required,gt=0"`
}

// ValuationResponse represents the response for waste valuation
type ValuationResponse struct {
	WasteType     string   `json:"waste_type"`
	Condition     string   `json:"condition"`
	WeightKg      float64  `json:"weight_kg"`
	PricePerKg    float64  `json:"price_per_kg"`
	TotalPrice    float64  `json:"total_price"`
	Currency      string   `json:"currency"`
	PricingRuleID *string  `json:"pricing_rule_id,omitempty"`
	Message       string   `json:"message,omitempty"`
}

// ToResponse converts WasteMetadata to WasteMetadataResponse
func (w *WasteMetadata) ToResponse() *WasteMetadataResponse {
	return &WasteMetadataResponse{
		ID:              w.ID,
		CollectionID:    w.CollectionID,
		WasteType:       w.WasteType,
		Condition:       w.Condition,
		ConfidenceScore: w.ConfidenceScore,
		ImageURL:        w.ImageURL,
		DetectedAt:      w.DetectedAt,
		ValuatedPrice:   w.ValuatedPrice,
		PricingRuleID:   w.PricingRuleID,
	}
}
