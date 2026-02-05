package models

import (
	"time"

	"github.com/google/uuid"
)

// SmartContract represents a deployed smart contract for a shipment
type SmartContract struct {
	ID               uuid.UUID `db:"id" json:"id"`
	ShipmentID       uuid.UUID `db:"shipment_id" json:"shipment_id"`
	ContractAddress  string    `db:"contract_address" json:"contract_address"`
	DeploymentTxHash string    `db:"deployment_tx_hash" json:"deployment_tx_hash"`
	ChainID          int       `db:"chain_id" json:"chain_id"`
	ABIVersion       string    `db:"abi_version" json:"abi_version"`
	IsActive         bool      `db:"is_active" json:"is_active"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

// ContractDeployRequest represents a request to deploy a smart contract
type ContractDeployRequest struct {
	ShipmentID uuid.UUID `json:"shipment_id"`
	UserID     uuid.UUID `json:"user_id"`
	Price      float64   `json:"price"`
	WasteType  string    `json:"waste_type"`
}

// ContractEventType represents types of blockchain events
type ContractEventType string

const (
	EventShipmentCreated   ContractEventType = "ShipmentCreated"
	EventPriceConfirmed    ContractEventType = "PriceConfirmed"
	EventDriverAssigned    ContractEventType = "DriverAssigned"
	EventStatusChanged     ContractEventType = "StatusChanged"
	EventDisputeRaised     ContractEventType = "DisputeRaised"
	EventShipmentCompleted ContractEventType = "ShipmentCompleted"
)

// ContractEvent represents a blockchain event
type ContractEvent struct {
	EventType       ContractEventType      `json:"event_type"`
	ContractAddress string                 `json:"contract_address"`
	TxHash          string                 `json:"tx_hash"`
	BlockNumber     uint64                 `json:"block_number"`
	ShipmentID      string                 `json:"shipment_id"`
	Data            map[string]interface{} `json:"data"`
	Timestamp       time.Time              `json:"timestamp"`
}

// SmartContractResponse represents the API response for a smart contract
type SmartContractResponse struct {
	ID               uuid.UUID `json:"id"`
	ShipmentID       uuid.UUID `json:"shipment_id"`
	ContractAddress  string    `json:"contract_address"`
	DeploymentTxHash string    `json:"deployment_tx_hash"`
	ChainID          int       `json:"chain_id"`
	ABIVersion       string    `json:"abi_version"`
	IsActive         bool      `json:"is_active"`
	ExplorerURL      string    `json:"explorer_url"`
	CreatedAt        time.Time `json:"created_at"`
}

// ToResponse converts SmartContract to SmartContractResponse
func (sc *SmartContract) ToResponse() *SmartContractResponse {
	// Generate explorer URL based on chain ID
	explorerURL := ""
	switch sc.ChainID {
	case 80001: // Polygon Mumbai
		explorerURL = "https://mumbai.polygonscan.com/address/" + sc.ContractAddress
	case 137: // Polygon Mainnet
		explorerURL = "https://polygonscan.com/address/" + sc.ContractAddress
	case 1: // Ethereum Mainnet
		explorerURL = "https://etherscan.io/address/" + sc.ContractAddress
	case 11155111: // Sepolia
		explorerURL = "https://sepolia.etherscan.io/address/" + sc.ContractAddress
	}

	return &SmartContractResponse{
		ID:               sc.ID,
		ShipmentID:       sc.ShipmentID,
		ContractAddress:  sc.ContractAddress,
		DeploymentTxHash: sc.DeploymentTxHash,
		ChainID:          sc.ChainID,
		ABIVersion:       sc.ABIVersion,
		IsActive:         sc.IsActive,
		ExplorerURL:      explorerURL,
		CreatedAt:        sc.CreatedAt,
	}
}
