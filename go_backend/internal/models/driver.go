package models

import (
	"time"

	"github.com/google/uuid"
)

// Driver represents a driver in the system
type Driver struct {
	ID               uuid.UUID `db:"id" json:"id"`
	Email            string    `db:"email" json:"email"`
	PasswordHash     string    `db:"password_hash" json:"-"`
	FullName         string    `db:"full_name" json:"full_name"`
	Phone            string    `db:"phone" json:"phone"`
	LicenseNumber    string    `db:"license_number" json:"license_number"`
	VehicleType      *string   `db:"vehicle_type" json:"vehicle_type,omitempty"`
	VehiclePlate     *string   `db:"vehicle_plate" json:"vehicle_plate,omitempty"`
	Latitude         *float64  `db:"latitude" json:"latitude,omitempty"`
	Longitude        *float64  `db:"longitude" json:"longitude,omitempty"`
	IsAvailable      bool      `db:"is_available" json:"is_available"`
	TotalCollections int       `db:"total_collections" json:"total_collections"`
	AverageRating    float64   `db:"average_rating" json:"average_rating"`
	FCMToken         *string   `db:"fcm_token" json:"-"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// CreateDriverRequest represents the request to create a new driver
type CreateDriverRequest struct {
	Email         string  `json:"email" binding:"required,email"`
	Password      string  `json:"password" binding:"required,min=8"`
	FullName      string  `json:"full_name" binding:"required"`
	Phone         string  `json:"phone" binding:"required"`
	LicenseNumber string  `json:"license_number" binding:"required"`
	VehicleType   *string `json:"vehicle_type"`
	VehiclePlate  *string `json:"vehicle_plate"`
}

// UpdateDriverRequest represents the request to update a driver
type UpdateDriverRequest struct {
	FullName     *string `json:"full_name"`
	Phone        *string `json:"phone"`
	VehicleType  *string `json:"vehicle_type"`
	VehiclePlate *string `json:"vehicle_plate"`
	IsAvailable  *bool   `json:"is_available"`
}

// UpdateDriverLocationRequest represents the request to update driver location
type UpdateDriverLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// DriverResponse represents the API response for a driver
type DriverResponse struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	FullName         string    `json:"full_name"`
	Phone            string    `json:"phone"`
	LicenseNumber    string    `json:"license_number"`
	VehicleType      *string   `json:"vehicle_type,omitempty"`
	VehiclePlate     *string   `json:"vehicle_plate,omitempty"`
	Latitude         *float64  `json:"latitude,omitempty"`
	Longitude        *float64  `json:"longitude,omitempty"`
	IsAvailable      bool      `json:"is_available"`
	TotalCollections int       `json:"total_collections"`
	AverageRating    float64   `json:"average_rating"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// VerifyTaskRequest represents the request to verify a task via QR code
type VerifyTaskRequest struct {
	QRCode       string `json:"qr_code" binding:"required"`
	CollectionID string `json:"collection_id" binding:"required,uuid"`
}

// ToResponse converts Driver to DriverResponse
func (d *Driver) ToResponse() *DriverResponse {
	return &DriverResponse{
		ID:               d.ID,
		Email:            d.Email,
		FullName:         d.FullName,
		Phone:            d.Phone,
		LicenseNumber:    d.LicenseNumber,
		VehicleType:      d.VehicleType,
		VehiclePlate:     d.VehiclePlate,
		Latitude:         d.Latitude,
		Longitude:        d.Longitude,
		IsAvailable:      d.IsAvailable,
		TotalCollections: d.TotalCollections,
		AverageRating:    d.AverageRating,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}
