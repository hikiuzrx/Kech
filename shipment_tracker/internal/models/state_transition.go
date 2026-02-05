package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// StateTransition represents an immutable record of a state change
type StateTransition struct {
	ID              uuid.UUID       `db:"id" json:"id"`
	ShipmentID      uuid.UUID       `db:"shipment_id" json:"shipment_id"`
	FromStatus      *ShipmentStatus `db:"from_status" json:"from_status,omitempty"`
	ToStatus        ShipmentStatus  `db:"to_status" json:"to_status"`
	TriggeredBy     uuid.UUID       `db:"triggered_by" json:"triggered_by"`
	TriggeredByRole string          `db:"triggered_by_role" json:"triggered_by_role"`
	ProofHash       *string         `db:"proof_hash" json:"proof_hash,omitempty"`
	Signature       *string         `db:"signature" json:"signature,omitempty"`
	TxHash          *string         `db:"tx_hash" json:"tx_hash,omitempty"`
	Metadata        json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
}

// TransitionMetadata holds additional data for a state transition
type TransitionMetadata struct {
	ActualWeight *float64  `json:"actual_weight_kg,omitempty"`
	Notes        *string   `json:"notes,omitempty"`
	Location     *Location `json:"location,omitempty"`
	DeviceInfo   *string   `json:"device_info,omitempty"`
	IPAddress    *string   `json:"ip_address,omitempty"`
}

// CreateTransitionRequest represents the request to create a state transition
type CreateTransitionRequest struct {
	ShipmentID      uuid.UUID           `json:"shipment_id"`
	FromStatus      *ShipmentStatus     `json:"from_status"`
	ToStatus        ShipmentStatus      `json:"to_status"`
	TriggeredBy     uuid.UUID           `json:"triggered_by"`
	TriggeredByRole string              `json:"triggered_by_role"`
	ProofHash       *string             `json:"proof_hash"`
	Signature       *string             `json:"signature"`
	Metadata        *TransitionMetadata `json:"metadata"`
}

// TransitionResponse represents the API response for a state transition
type TransitionResponse struct {
	ID              uuid.UUID       `json:"id"`
	ShipmentID      uuid.UUID       `json:"shipment_id"`
	FromStatus      *ShipmentStatus `json:"from_status,omitempty"`
	ToStatus        ShipmentStatus  `json:"to_status"`
	TriggeredBy     uuid.UUID       `json:"triggered_by"`
	TriggeredByRole string          `json:"triggered_by_role"`
	ProofHash       *string         `json:"proof_hash,omitempty"`
	TxHash          *string         `json:"tx_hash,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// ToResponse converts StateTransition to TransitionResponse
func (t *StateTransition) ToResponse() *TransitionResponse {
	return &TransitionResponse{
		ID:              t.ID,
		ShipmentID:      t.ShipmentID,
		FromStatus:      t.FromStatus,
		ToStatus:        t.ToStatus,
		TriggeredBy:     t.TriggeredBy,
		TriggeredByRole: t.TriggeredByRole,
		ProofHash:       t.ProofHash,
		TxHash:          t.TxHash,
		CreatedAt:       t.CreatedAt,
	}
}
