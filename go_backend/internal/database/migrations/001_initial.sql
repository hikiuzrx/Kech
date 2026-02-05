-- Smart Waste Management Database Schema
-- PostgreSQL Migration: 001_initial.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table: stores user profiles and reward points
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    reward_points INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Drivers table: stores driver profiles and performance metrics
CREATE TABLE drivers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    license_number VARCHAR(50) NOT NULL,
    vehicle_type VARCHAR(50),
    vehicle_plate VARCHAR(20),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    is_available BOOLEAN DEFAULT true,
    total_collections INTEGER DEFAULT 0,
    average_rating DECIMAL(3, 2) DEFAULT 0.00,
    fcm_token VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Recycling companies table
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    registration_number VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Smart bins table: stores bin status and location
CREATE TABLE bins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id VARCHAR(100) UNIQUE NOT NULL,
    location_name VARCHAR(255),
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    fill_level INTEGER DEFAULT 0 CHECK (fill_level >= 0 AND fill_level <= 100),
    waste_type VARCHAR(50) DEFAULT 'general',
    capacity_liters INTEGER DEFAULT 240,
    last_collection_at TIMESTAMP WITH TIME ZONE,
    last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    company_id UUID REFERENCES companies(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Collection history table
CREATE TABLE collections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bin_id UUID NOT NULL REFERENCES bins(id) ON DELETE CASCADE,
    driver_id UUID NOT NULL REFERENCES drivers(id) ON DELETE CASCADE,
    fill_level_before INTEGER NOT NULL,
    fill_level_after INTEGER DEFAULT 0,
    weight_kg DECIMAL(10, 2),
    qr_code_verified BOOLEAN DEFAULT false,
    notes TEXT,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'pending'
);

-- Pricing rules table: for waste valuation engine
CREATE TABLE pricing_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    waste_type VARCHAR(50) NOT NULL,
    condition VARCHAR(50) NOT NULL,
    price_per_kg DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    min_weight_kg DECIMAL(10, 2) DEFAULT 0,
    max_weight_kg DECIMAL(10, 2),
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Waste metadata table: AI-detected waste classification
CREATE TABLE waste_metadata (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    collection_id UUID REFERENCES collections(id) ON DELETE CASCADE,
    waste_type VARCHAR(50) NOT NULL,
    condition VARCHAR(50) NOT NULL,
    confidence_score DECIMAL(5, 4),
    image_url TEXT,
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    valuated_price DECIMAL(10, 2),
    pricing_rule_id UUID REFERENCES pricing_rules(id) ON DELETE SET NULL
);

-- Notifications table: notification history
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    driver_id UUID REFERENCES drivers(id) ON DELETE CASCADE,
    bin_id UUID REFERENCES bins(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP WITH TIME ZONE
);

-- Driver routes table: stores optimized routes
CREATE TABLE driver_routes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    driver_id UUID NOT NULL REFERENCES drivers(id) ON DELETE CASCADE,
    waypoints JSONB NOT NULL,
    total_distance_km DECIMAL(10, 2),
    estimated_duration_minutes INTEGER,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for performance
CREATE INDEX idx_bins_fill_level ON bins(fill_level);
CREATE INDEX idx_bins_company ON bins(company_id);
CREATE INDEX idx_bins_device_id ON bins(device_id);
CREATE INDEX idx_collections_bin ON collections(bin_id);
CREATE INDEX idx_collections_driver ON collections(driver_id);
CREATE INDEX idx_collections_status ON collections(status);
CREATE INDEX idx_drivers_location ON drivers(latitude, longitude);
CREATE INDEX idx_drivers_available ON drivers(is_available);
CREATE INDEX idx_notifications_driver ON notifications(driver_id);
CREATE INDEX idx_notifications_unread ON notifications(driver_id, is_read) WHERE NOT is_read;
CREATE INDEX idx_pricing_rules_type ON pricing_rules(waste_type, condition);

-- Update timestamp trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply update triggers
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_drivers_updated_at BEFORE UPDATE ON drivers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pricing_rules_updated_at BEFORE UPDATE ON pricing_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
