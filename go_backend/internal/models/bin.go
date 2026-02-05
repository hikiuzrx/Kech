package models

import (
	"time"

	"github.com/google/uuid"
)

// Bin represents a smart waste bin with IoT sensors
type Bin struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	DeviceID         string     `db:"device_id" json:"device_id"`
	LocationName     *string    `db:"location_name" json:"location_name,omitempty"`
	Latitude         float64    `db:"latitude" json:"latitude"`
	Longitude        float64    `db:"longitude" json:"longitude"`
	FillLevel        int        `db:"fill_level" json:"fill_level"`
	WasteType        string     `db:"waste_type" json:"waste_type"`
	CapacityLiters   int        `db:"capacity_liters" json:"capacity_liters"`
	LastCollectionAt *time.Time `db:"last_collection_at" json:"last_collection_at,omitempty"`
	LastUpdatedAt    time.Time  `db:"last_updated_at" json:"last_updated_at"`
	IsActive         bool       `db:"is_active" json:"is_active"`
	CompanyID        *uuid.UUID `db:"company_id" json:"company_id,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
}

// CreateBinRequest represents the request to register a new bin
type CreateBinRequest struct {
	DeviceID       string     `json:"device_id" binding:"required"`
	LocationName   *string    `json:"location_name"`
	Latitude       float64    `json:"latitude" binding:"required"`
	Longitude      float64    `json:"longitude" binding:"required"`
	WasteType      string     `json:"waste_type" binding:"required"`
	CapacityLiters int        `json:"capacity_liters" binding:"required,gt=0"`
	CompanyID      *uuid.UUID `json:"company_id"`
}

// UpdateBinRequest represents the request to update a bin
type UpdateBinRequest struct {
	LocationName   *string    `json:"location_name"`
	Latitude       *float64   `json:"latitude"`
	Longitude      *float64   `json:"longitude"`
	WasteType      *string    `json:"waste_type"`
	CapacityLiters *int       `json:"capacity_liters"`
	IsActive       *bool      `json:"is_active"`
	CompanyID      *uuid.UUID `json:"company_id"`
}

// BinStatusUpdate represents IoT payload from ESP32
type BinStatusUpdate struct {
	BinID     string `json:"bin_id"`
	FillLevel int    `json:"fill_level"`
}

// BinResponse represents the API response for a bin
type BinResponse struct {
	ID               uuid.UUID  `json:"id"`
	DeviceID         string     `json:"device_id"`
	LocationName     *string    `json:"location_name,omitempty"`
	Latitude         float64    `json:"latitude"`
	Longitude        float64    `json:"longitude"`
	FillLevel        int        `json:"fill_level"`
	WasteType        string     `json:"waste_type"`
	CapacityLiters   int        `json:"capacity_liters"`
	LastCollectionAt *time.Time `json:"last_collection_at,omitempty"`
	LastUpdatedAt    time.Time  `json:"last_updated_at"`
	IsActive         bool       `json:"is_active"`
	CompanyID        *uuid.UUID `json:"company_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// ToResponse converts Bin to BinResponse
func (b *Bin) ToResponse() *BinResponse {
	return &BinResponse{
		ID:               b.ID,
		DeviceID:         b.DeviceID,
		LocationName:     b.LocationName,
		Latitude:         b.Latitude,
		Longitude:        b.Longitude,
		FillLevel:        b.FillLevel,
		WasteType:        b.WasteType,
		CapacityLiters:   b.CapacityLiters,
		LastCollectionAt: b.LastCollectionAt,
		LastUpdatedAt:    b.LastUpdatedAt,
		IsActive:         b.IsActive,
		CompanyID:        b.CompanyID,
		CreatedAt:        b.CreatedAt,
	}
}

// NeedsCollection returns true if the bin fill level exceeds the threshold
func (b *Bin) NeedsCollection(threshold int) bool {
	return b.FillLevel >= threshold
}
