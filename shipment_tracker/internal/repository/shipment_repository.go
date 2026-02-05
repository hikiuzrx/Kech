package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smartwaste/shipment-tracker/internal/models"
)

// ShipmentRepository handles database operations for shipments
type ShipmentRepository struct {
	db *sqlx.DB
}

// NewShipmentRepository creates a new ShipmentRepository
func NewShipmentRepository(db *sqlx.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

// Create creates a new shipment
func (r *ShipmentRepository) Create(s *models.Shipment) error {
	query := `
		INSERT INTO shipments (
			id, user_id, collection_id, waste_type, estimated_weight_kg,
			price_offered, price_confirmed, status,
			pickup_latitude, pickup_longitude, pickup_address,
			dropoff_latitude, dropoff_longitude, dropoff_address,
			notes, created_at, updated_at
		) VALUES (
			:id, :user_id, :collection_id, :waste_type, :estimated_weight_kg,
			:price_offered, :price_confirmed, :status,
			:pickup_latitude, :pickup_longitude, :pickup_address,
			:dropoff_latitude, :dropoff_longitude, :dropoff_address,
			:notes, :created_at, :updated_at
		)`

	_, err := r.db.NamedExec(query, s)
	return err
}

// GetByID retrieves a shipment by ID
func (r *ShipmentRepository) GetByID(id uuid.UUID) (*models.Shipment, error) {
	var s models.Shipment
	err := r.db.Get(&s, "SELECT * FROM shipments WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	return &s, err
}

// UpdateStatus updates the status of a shipment
func (r *ShipmentRepository) UpdateStatus(id uuid.UUID, status models.ShipmentStatus) error {
	_, err := r.db.Exec("UPDATE shipments SET status = $1 WHERE id = $2", status, id)
	return err
}

// UpdateContractDetails updates the smart contract details for a shipment
func (r *ShipmentRepository) UpdateContractDetails(id uuid.UUID, address, txHash string) error {
	_, err := r.db.Exec("UPDATE shipments SET contract_address = $1, contract_tx_hash = $2 WHERE id = $3", address, txHash, id)
	return err
}

// AssignDriver assigns a driver to a shipment
func (r *ShipmentRepository) AssignDriver(id uuid.UUID, driverID uuid.UUID) error {
	_, err := r.db.Exec("UPDATE shipments SET driver_id = $1, status = $2 WHERE id = $3", driverID, models.StatusDriverAssigned, id)
	return err
}

// UpdateActualWeight updates the actual weight of the shipment
func (r *ShipmentRepository) UpdateActualWeight(id uuid.UUID, weight float64) error {
	_, err := r.db.Exec("UPDATE shipments SET actual_weight_kg = $1 WHERE id = $2", weight, id)
	return err
}

// List retrieves a list of shipments with optional filtering
func (r *ShipmentRepository) List(userID *uuid.UUID, driverID *uuid.UUID, status *models.ShipmentStatus) ([]models.Shipment, error) {
	query := "SELECT * FROM shipments WHERE 1=1"
	args := []interface{}{}
	argID := 1

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argID)
		args = append(args, *userID)
		argID++
	}

	if driverID != nil {
		query += fmt.Sprintf(" AND driver_id = $%d", argID)
		args = append(args, *driverID)
		argID++
	}

	if status != nil {
		query += fmt.Sprintf(" AND status = $%d", argID)
		args = append(args, *status)
		argID++
	}

	query += " ORDER BY created_at DESC"

	var shipments []models.Shipment
	err := r.db.Select(&shipments, query, args...)
	return shipments, err
}
