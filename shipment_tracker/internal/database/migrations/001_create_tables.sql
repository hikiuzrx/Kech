-- Shipment Tracker Database Schema
-- Migration: 001_create_tables.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Shipment status enum
CREATE TYPE shipment_status AS ENUM (
    'created',
    'price_confirmed',
    'driver_assigned',
    'pickup_started',
    'in_transit',
    'delivered',
    'completed',
    'cancelled',
    'disputed',
    'resolved'
);

-- Shipments table
CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    driver_id UUID,
    collection_id UUID NOT NULL,
    waste_type VARCHAR(100) NOT NULL,
    estimated_weight_kg DECIMAL(10, 2) NOT NULL,
    actual_weight_kg DECIMAL(10, 2),
    price_offered DECIMAL(12, 2) NOT NULL,
    price_confirmed BOOLEAN DEFAULT FALSE,
    contract_address VARCHAR(66),
    contract_tx_hash VARCHAR(66),
    status shipment_status NOT NULL DEFAULT 'created',
    pickup_latitude DECIMAL(10, 8),
    pickup_longitude DECIMAL(11, 8),
    pickup_address TEXT,
    dropoff_latitude DECIMAL(10, 8),
    dropoff_longitude DECIMAL(11, 8),
    dropoff_address TEXT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- State transitions table (immutable audit log)
CREATE TABLE IF NOT EXISTS state_transitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shipment_id UUID NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    from_status shipment_status,
    to_status shipment_status NOT NULL,
    triggered_by UUID NOT NULL,
    triggered_by_role VARCHAR(50) NOT NULL, -- 'user', 'driver', 'system', 'admin'
    proof_hash VARCHAR(66), -- IPFS or blockchain hash of evidence
    signature VARCHAR(132), -- Digital signature
    tx_hash VARCHAR(66), -- Blockchain transaction hash
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Smart contract records table
CREATE TABLE IF NOT EXISTS smart_contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shipment_id UUID NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    contract_address VARCHAR(66) NOT NULL,
    deployment_tx_hash VARCHAR(66) NOT NULL,
    chain_id INTEGER NOT NULL,
    abi_version VARCHAR(20) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Disputes table
CREATE TABLE IF NOT EXISTS disputes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shipment_id UUID NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    raised_by UUID NOT NULL,
    raised_by_role VARCHAR(50) NOT NULL,
    reason TEXT NOT NULL,
    evidence_hash VARCHAR(66),
    resolution TEXT,
    resolved_by UUID,
    resolved_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'open', -- open, investigating, resolved
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_shipments_user_id ON shipments(user_id);
CREATE INDEX idx_shipments_driver_id ON shipments(driver_id);
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_shipments_created_at ON shipments(created_at);
CREATE INDEX idx_state_transitions_shipment_id ON state_transitions(shipment_id);
CREATE INDEX idx_state_transitions_created_at ON state_transitions(created_at);
CREATE INDEX idx_smart_contracts_shipment_id ON smart_contracts(shipment_id);
CREATE INDEX idx_disputes_shipment_id ON disputes(shipment_id);
CREATE INDEX idx_disputes_status ON disputes(status);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_shipments_updated_at
    BEFORE UPDATE ON shipments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_disputes_updated_at
    BEFORE UPDATE ON disputes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
