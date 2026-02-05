package models

import (
	"time"

	"github.com/google/uuid"
)

// DisputeStatus represents the status of a dispute
type DisputeStatus string

const (
	DisputeStatusOpen          DisputeStatus = "open"
	DisputeStatusInvestigating DisputeStatus = "investigating"
	DisputeStatusResolved      DisputeStatus = "resolved"
)

// Dispute represents a dispute raised on a shipment
type Dispute struct {
	ID           uuid.UUID     `db:"id" json:"id"`
	ShipmentID   uuid.UUID     `db:"shipment_id" json:"shipment_id"`
	RaisedBy     uuid.UUID     `db:"raised_by" json:"raised_by"`
	RaisedByRole string        `db:"raised_by_role" json:"raised_by_role"`
	Reason       string        `db:"reason" json:"reason"`
	EvidenceHash *string       `db:"evidence_hash" json:"evidence_hash,omitempty"`
	Resolution   *string       `db:"resolution" json:"resolution,omitempty"`
	ResolvedBy   *uuid.UUID    `db:"resolved_by" json:"resolved_by,omitempty"`
	ResolvedAt   *time.Time    `db:"resolved_at" json:"resolved_at,omitempty"`
	Status       DisputeStatus `db:"status" json:"status"`
	CreatedAt    time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at" json:"updated_at"`
}

// ResolveDisputeRequest represents the request to resolve a dispute
type ResolveDisputeRequest struct {
	ResolvedBy uuid.UUID `json:"resolved_by" binding:"required"`
	Resolution string    `json:"resolution" binding:"required"`
	Outcome    string    `json:"outcome" binding:"required"` // "user_wins", "driver_wins", "split"
}

// DisputeResponse represents the API response for a dispute
type DisputeResponse struct {
	ID           uuid.UUID     `json:"id"`
	ShipmentID   uuid.UUID     `json:"shipment_id"`
	RaisedBy     uuid.UUID     `json:"raised_by"`
	RaisedByRole string        `json:"raised_by_role"`
	Reason       string        `json:"reason"`
	EvidenceHash *string       `json:"evidence_hash,omitempty"`
	Resolution   *string       `json:"resolution,omitempty"`
	ResolvedBy   *uuid.UUID    `json:"resolved_by,omitempty"`
	ResolvedAt   *time.Time    `json:"resolved_at,omitempty"`
	Status       DisputeStatus `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
}

// ToResponse converts Dispute to DisputeResponse
func (d *Dispute) ToResponse() *DisputeResponse {
	return &DisputeResponse{
		ID:           d.ID,
		ShipmentID:   d.ShipmentID,
		RaisedBy:     d.RaisedBy,
		RaisedByRole: d.RaisedByRole,
		Reason:       d.Reason,
		EvidenceHash: d.EvidenceHash,
		Resolution:   d.Resolution,
		ResolvedBy:   d.ResolvedBy,
		ResolvedAt:   d.ResolvedAt,
		Status:       d.Status,
		CreatedAt:    d.CreatedAt,
	}
}
