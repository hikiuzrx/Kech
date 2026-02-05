package models

import (
	"time"

	"github.com/google/uuid"
)

// CollectionStatus represents the status of a collection
type CollectionStatus string

const (
	CollectionStatusPending    CollectionStatus = "pending"
	CollectionStatusInProgress CollectionStatus = "in_progress"
	CollectionStatusCompleted  CollectionStatus = "completed"
	CollectionStatusCancelled  CollectionStatus = "cancelled"
)

// Collection represents a waste collection event
type Collection struct {
	ID              uuid.UUID        `db:"id" json:"id"`
	BinID           uuid.UUID        `db:"bin_id" json:"bin_id"`
	DriverID        uuid.UUID        `db:"driver_id" json:"driver_id"`
	FillLevelBefore int              `db:"fill_level_before" json:"fill_level_before"`
	FillLevelAfter  int              `db:"fill_level_after" json:"fill_level_after"`
	WeightKg        *float64         `db:"weight_kg" json:"weight_kg,omitempty"`
	QRCodeVerified  bool             `db:"qr_code_verified" json:"qr_code_verified"`
	Notes           *string          `db:"notes" json:"notes,omitempty"`
	StartedAt       time.Time        `db:"started_at" json:"started_at"`
	CompletedAt     *time.Time       `db:"completed_at" json:"completed_at,omitempty"`
	Status          CollectionStatus `db:"status" json:"status"`
}

// CreateCollectionRequest represents the request to create a new collection
type CreateCollectionRequest struct {
	BinID    uuid.UUID `json:"bin_id" binding:"required"`
	DriverID uuid.UUID `json:"driver_id" binding:"required"`
}

// UpdateCollectionRequest represents the request to update a collection
type UpdateCollectionRequest struct {
	FillLevelAfter *int     `json:"fill_level_after"`
	WeightKg       *float64 `json:"weight_kg"`
	Notes          *string  `json:"notes"`
	Status         *string  `json:"status"`
}

// CompleteCollectionRequest represents the request to complete a collection
type CompleteCollectionRequest struct {
	FillLevelAfter int      `json:"fill_level_after" binding:"required,gte=0,lte=100"`
	WeightKg       *float64 `json:"weight_kg"`
	Notes          *string  `json:"notes"`
}

// CollectionResponse represents the API response for a collection
type CollectionResponse struct {
	ID              uuid.UUID        `json:"id"`
	BinID           uuid.UUID        `json:"bin_id"`
	DriverID        uuid.UUID        `json:"driver_id"`
	FillLevelBefore int              `json:"fill_level_before"`
	FillLevelAfter  int              `json:"fill_level_after"`
	WeightKg        *float64         `json:"weight_kg,omitempty"`
	QRCodeVerified  bool             `json:"qr_code_verified"`
	Notes           *string          `json:"notes,omitempty"`
	StartedAt       time.Time        `json:"started_at"`
	CompletedAt     *time.Time       `json:"completed_at,omitempty"`
	Status          CollectionStatus `json:"status"`
}

// ToResponse converts Collection to CollectionResponse
func (c *Collection) ToResponse() *CollectionResponse {
	return &CollectionResponse{
		ID:              c.ID,
		BinID:           c.BinID,
		DriverID:        c.DriverID,
		FillLevelBefore: c.FillLevelBefore,
		FillLevelAfter:  c.FillLevelAfter,
		WeightKg:        c.WeightKg,
		QRCodeVerified:  c.QRCodeVerified,
		Notes:           c.Notes,
		StartedAt:       c.StartedAt,
		CompletedAt:     c.CompletedAt,
		Status:          c.Status,
	}
}
