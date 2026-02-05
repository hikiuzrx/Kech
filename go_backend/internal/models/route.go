package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// RouteStatus represents the status of a driver route
type RouteStatus string

const (
	RouteStatusPending    RouteStatus = "pending"
	RouteStatusInProgress RouteStatus = "in_progress"
	RouteStatusCompleted  RouteStatus = "completed"
	RouteStatusCancelled  RouteStatus = "cancelled"
)

// Waypoint represents a single point in a route
type Waypoint struct {
	BinID       uuid.UUID `json:"bin_id"`
	DeviceID    string    `json:"device_id"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	FillLevel   int       `json:"fill_level"`
	Order       int       `json:"order"`
	IsCompleted bool      `json:"is_completed"`
}

// DriverRoute represents an optimized route for a driver
type DriverRoute struct {
	ID                       uuid.UUID       `db:"id" json:"id"`
	DriverID                 uuid.UUID       `db:"driver_id" json:"driver_id"`
	Waypoints                json.RawMessage `db:"waypoints" json:"-"`
	WaypointsList            []Waypoint      `db:"-" json:"waypoints"`
	TotalDistanceKm          *float64        `db:"total_distance_km" json:"total_distance_km,omitempty"`
	EstimatedDurationMinutes *int            `db:"estimated_duration_minutes" json:"estimated_duration_minutes,omitempty"`
	Status                   RouteStatus     `db:"status" json:"status"`
	CreatedAt                time.Time       `db:"created_at" json:"created_at"`
	StartedAt                *time.Time      `db:"started_at" json:"started_at,omitempty"`
	CompletedAt              *time.Time      `db:"completed_at" json:"completed_at,omitempty"`
}

// CreateRouteRequest represents the request to create a route
type CreateRouteRequest struct {
	DriverID   uuid.UUID   `json:"driver_id" binding:"required"`
	BinIDs     []uuid.UUID `json:"bin_ids" binding:"required,min=1"`
	OptimizeBy string      `json:"optimize_by"` // "distance" or "fill_level"
}

// RouteResponse represents the API response for a route
type RouteResponse struct {
	ID                       uuid.UUID   `json:"id"`
	DriverID                 uuid.UUID   `json:"driver_id"`
	Waypoints                []Waypoint  `json:"waypoints"`
	TotalDistanceKm          *float64    `json:"total_distance_km,omitempty"`
	EstimatedDurationMinutes *int        `json:"estimated_duration_minutes,omitempty"`
	Status                   RouteStatus `json:"status"`
	CreatedAt                time.Time   `json:"created_at"`
	StartedAt                *time.Time  `json:"started_at,omitempty"`
	CompletedAt              *time.Time  `json:"completed_at,omitempty"`
}

// ParseWaypoints parses the JSON waypoints into the WaypointsList
func (r *DriverRoute) ParseWaypoints() error {
	if len(r.Waypoints) > 0 {
		return json.Unmarshal(r.Waypoints, &r.WaypointsList)
	}
	return nil
}

// ToResponse converts DriverRoute to RouteResponse
func (r *DriverRoute) ToResponse() *RouteResponse {
	_ = r.ParseWaypoints()
	return &RouteResponse{
		ID:                       r.ID,
		DriverID:                 r.DriverID,
		Waypoints:                r.WaypointsList,
		TotalDistanceKm:          r.TotalDistanceKm,
		EstimatedDurationMinutes: r.EstimatedDurationMinutes,
		Status:                   r.Status,
		CreatedAt:                r.CreatedAt,
		StartedAt:                r.StartedAt,
		CompletedAt:              r.CompletedAt,
	}
}
