package models

import (
	"time"

	"github.com/google/uuid"
)

// ShipmentStatus represents the status of a shipment
type ShipmentStatus string

const (
	StatusCreated        ShipmentStatus = "created"
	StatusPriceConfirmed ShipmentStatus = "price_confirmed"
	StatusDriverAssigned ShipmentStatus = "driver_assigned"
	StatusPickupStarted  ShipmentStatus = "pickup_started"
	StatusInTransit      ShipmentStatus = "in_transit"
	StatusDelivered      ShipmentStatus = "delivered"
	StatusCompleted      ShipmentStatus = "completed"
	StatusCancelled      ShipmentStatus = "cancelled"
	StatusDisputed       ShipmentStatus = "disputed"
	StatusResolved       ShipmentStatus = "resolved"
)

// ValidTransitions defines valid state transitions
var ValidTransitions = map[ShipmentStatus][]ShipmentStatus{
	StatusCreated:        {StatusPriceConfirmed, StatusCancelled},
	StatusPriceConfirmed: {StatusDriverAssigned, StatusCancelled},
	StatusDriverAssigned: {StatusPickupStarted, StatusDisputed, StatusCancelled},
	StatusPickupStarted:  {StatusInTransit, StatusDisputed},
	StatusInTransit:      {StatusDelivered, StatusDisputed},
	StatusDelivered:      {StatusCompleted, StatusDisputed},
	StatusDisputed:       {StatusResolved},
	StatusResolved:       {StatusCompleted},
}

// Location represents a geographic location
type Location struct {
	Latitude  float64 `db:"latitude" json:"latitude"`
	Longitude float64 `db:"longitude" json:"longitude"`
	Address   string  `db:"address" json:"address,omitempty"`
}

// Shipment represents a waste shipment
type Shipment struct {
	ID                uuid.UUID      `db:"id" json:"id"`
	UserID            uuid.UUID      `db:"user_id" json:"user_id"`
	DriverID          *uuid.UUID     `db:"driver_id" json:"driver_id,omitempty"`
	CollectionID      uuid.UUID      `db:"collection_id" json:"collection_id"`
	WasteType         string         `db:"waste_type" json:"waste_type"`
	EstimatedWeightKg float64        `db:"estimated_weight_kg" json:"estimated_weight_kg"`
	ActualWeightKg    *float64       `db:"actual_weight_kg" json:"actual_weight_kg,omitempty"`
	PriceOffered      float64        `db:"price_offered" json:"price_offered"`
	PriceConfirmed    bool           `db:"price_confirmed" json:"price_confirmed"`
	ContractAddress   *string        `db:"contract_address" json:"contract_address,omitempty"`
	ContractTxHash    *string        `db:"contract_tx_hash" json:"contract_tx_hash,omitempty"`
	Status            ShipmentStatus `db:"status" json:"status"`
	PickupLatitude    *float64       `db:"pickup_latitude" json:"pickup_latitude,omitempty"`
	PickupLongitude   *float64       `db:"pickup_longitude" json:"pickup_longitude,omitempty"`
	PickupAddress     *string        `db:"pickup_address" json:"pickup_address,omitempty"`
	DropoffLatitude   *float64       `db:"dropoff_latitude" json:"dropoff_latitude,omitempty"`
	DropoffLongitude  *float64       `db:"dropoff_longitude" json:"dropoff_longitude,omitempty"`
	DropoffAddress    *string        `db:"dropoff_address" json:"dropoff_address,omitempty"`
	Notes             *string        `db:"notes" json:"notes,omitempty"`
	CreatedAt         time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at" json:"updated_at"`
}

// CreateShipmentRequest represents the request to create a new shipment
type CreateShipmentRequest struct {
	UserID            uuid.UUID `json:"user_id" binding:"required"`
	CollectionID      uuid.UUID `json:"collection_id" binding:"required"`
	WasteType         string    `json:"waste_type" binding:"required"`
	EstimatedWeightKg float64   `json:"estimated_weight_kg" binding:"required,gt=0"`
	PriceOffered      float64   `json:"price_offered" binding:"required,gt=0"`
	PickupLocation    *Location `json:"pickup_location"`
	DropoffLocation   *Location `json:"dropoff_location"`
	Notes             *string   `json:"notes"`
}

// AssignDriverRequest represents the request to assign a driver
type AssignDriverRequest struct {
	DriverID uuid.UUID `json:"driver_id" binding:"required"`
}

// ConfirmPickupRequest represents the request to confirm pickup
type ConfirmPickupRequest struct {
	ConfirmedBy  uuid.UUID `json:"confirmed_by" binding:"required"`
	Role         string    `json:"role" binding:"required"` // "user" or "driver"
	ProofHash    *string   `json:"proof_hash"`
	ActualWeight *float64  `json:"actual_weight_kg"`
	Signature    string    `json:"signature" binding:"required"`
}

// ConfirmDeliveryRequest represents the request to confirm delivery
type ConfirmDeliveryRequest struct {
	ConfirmedBy uuid.UUID `json:"confirmed_by" binding:"required"`
	Role        string    `json:"role" binding:"required"`
	ProofHash   *string   `json:"proof_hash"`
	Signature   string    `json:"signature" binding:"required"`
}

// RaiseDisputeRequest represents the request to raise a dispute
type RaiseDisputeRequest struct {
	RaisedBy     uuid.UUID `json:"raised_by" binding:"required"`
	Role         string    `json:"role" binding:"required"`
	Reason       string    `json:"reason" binding:"required"`
	EvidenceHash *string   `json:"evidence_hash"`
}

// ShipmentResponse represents the API response for a shipment
type ShipmentResponse struct {
	ID                uuid.UUID      `json:"id"`
	UserID            uuid.UUID      `json:"user_id"`
	DriverID          *uuid.UUID     `json:"driver_id,omitempty"`
	CollectionID      uuid.UUID      `json:"collection_id"`
	WasteType         string         `json:"waste_type"`
	EstimatedWeightKg float64        `json:"estimated_weight_kg"`
	ActualWeightKg    *float64       `json:"actual_weight_kg,omitempty"`
	PriceOffered      float64        `json:"price_offered"`
	PriceConfirmed    bool           `json:"price_confirmed"`
	ContractAddress   *string        `json:"contract_address,omitempty"`
	Status            ShipmentStatus `json:"status"`
	PickupLocation    *Location      `json:"pickup_location,omitempty"`
	DropoffLocation   *Location      `json:"dropoff_location,omitempty"`
	Notes             *string        `json:"notes,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// ToResponse converts Shipment to ShipmentResponse
func (s *Shipment) ToResponse() *ShipmentResponse {
	resp := &ShipmentResponse{
		ID:                s.ID,
		UserID:            s.UserID,
		DriverID:          s.DriverID,
		CollectionID:      s.CollectionID,
		WasteType:         s.WasteType,
		EstimatedWeightKg: s.EstimatedWeightKg,
		ActualWeightKg:    s.ActualWeightKg,
		PriceOffered:      s.PriceOffered,
		PriceConfirmed:    s.PriceConfirmed,
		ContractAddress:   s.ContractAddress,
		Status:            s.Status,
		Notes:             s.Notes,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}

	if s.PickupLatitude != nil && s.PickupLongitude != nil {
		resp.PickupLocation = &Location{
			Latitude:  *s.PickupLatitude,
			Longitude: *s.PickupLongitude,
		}
		if s.PickupAddress != nil {
			resp.PickupLocation.Address = *s.PickupAddress
		}
	}

	if s.DropoffLatitude != nil && s.DropoffLongitude != nil {
		resp.DropoffLocation = &Location{
			Latitude:  *s.DropoffLatitude,
			Longitude: *s.DropoffLongitude,
		}
		if s.DropoffAddress != nil {
			resp.DropoffLocation.Address = *s.DropoffAddress
		}
	}

	return resp
}

// CanTransitionTo checks if the shipment can transition to the given status
func (s *Shipment) CanTransitionTo(newStatus ShipmentStatus) bool {
	validNextStates, exists := ValidTransitions[s.Status]
	if !exists {
		return false
	}

	for _, validStatus := range validNextStates {
		if validStatus == newStatus {
			return true
		}
	}
	return false
}
