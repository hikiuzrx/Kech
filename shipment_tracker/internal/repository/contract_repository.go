package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smartwaste/shipment-tracker/internal/models"
)

// ContractRepository handles database operations for smart contracts
type ContractRepository struct {
	db *sqlx.DB
}

// NewContractRepository creates a new ContractRepository
func NewContractRepository(db *sqlx.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

// Create stores a new smart contract record
func (r *ContractRepository) Create(sc *models.SmartContract) error {
	query := `
		INSERT INTO smart_contracts (
			id, shipment_id, contract_address, deployment_tx_hash,
			chain_id, abi_version, is_active, created_at
		) VALUES (
			:id, :shipment_id, :contract_address, :deployment_tx_hash,
			:chain_id, :abi_version, :is_active, :created_at
		)`

	_, err := r.db.NamedExec(query, sc)
	return err
}

// GetByShipmentID gets the smart contract for a shipment
func (r *ContractRepository) GetByShipmentID(shipmentID uuid.UUID) (*models.SmartContract, error) {
	var sc models.SmartContract
	err := r.db.Get(&sc, "SELECT * FROM smart_contracts WHERE shipment_id = $1", shipmentID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &sc, err
}
