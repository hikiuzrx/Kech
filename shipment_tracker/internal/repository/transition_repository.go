package repository

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smartwaste/shipment-tracker/internal/models"
)

// TransitionRepository handles database operations for state transitions
type TransitionRepository struct {
	db *sqlx.DB
}

// NewTransitionRepository creates a new TransitionRepository
func NewTransitionRepository(db *sqlx.DB) *TransitionRepository {
	return &TransitionRepository{db: db}
}

// Create creates a new state transition record
func (r *TransitionRepository) Create(t *models.StateTransition) error {
	// Ensure Metadata is valid JSON if nil
	if t.Metadata == nil {
		t.Metadata = json.RawMessage("{}")
	}

	query := `
		INSERT INTO state_transitions (
			id, shipment_id, from_status, to_status,
			triggered_by, triggered_by_role,
			proof_hash, signature, tx_hash, metadata, created_at
		) VALUES (
			:id, :shipment_id, :from_status, :to_status,
			:triggered_by, :triggered_by_role,
			:proof_hash, :signature, :tx_hash, :metadata, :created_at
		)`

	_, err := r.db.NamedExec(query, t)
	return err
}

// GetByShipmentID retrieves all transitions for a shipment
func (r *TransitionRepository) GetByShipmentID(shipmentID uuid.UUID) ([]models.StateTransition, error) {
	var transitions []models.StateTransition
	err := r.db.Select(&transitions, "SELECT * FROM state_transitions WHERE shipment_id = $1 ORDER BY created_at ASC", shipmentID)
	return transitions, err
}
