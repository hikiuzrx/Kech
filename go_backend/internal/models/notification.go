package models

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeBinFull        NotificationType = "bin_full"
	NotificationTypeRouteAssigned  NotificationType = "route_assigned"
	NotificationTypeTaskCompleted  NotificationType = "task_completed"
	NotificationTypeSystemAlert    NotificationType = "system_alert"
)

// Notification represents a notification sent to a driver
type Notification struct {
	ID       uuid.UUID         `db:"id" json:"id"`
	DriverID *uuid.UUID        `db:"driver_id" json:"driver_id,omitempty"`
	BinID    *uuid.UUID        `db:"bin_id" json:"bin_id,omitempty"`
	Type     NotificationType  `db:"type" json:"type"`
	Title    string            `db:"title" json:"title"`
	Message  string            `db:"message" json:"message"`
	IsRead   bool              `db:"is_read" json:"is_read"`
	SentAt   time.Time         `db:"sent_at" json:"sent_at"`
	ReadAt   *time.Time        `db:"read_at" json:"read_at,omitempty"`
}

// CreateNotificationRequest represents the request to create a notification
type CreateNotificationRequest struct {
	DriverID *uuid.UUID       `json:"driver_id"`
	BinID    *uuid.UUID       `json:"bin_id"`
	Type     NotificationType `json:"type" binding:"required"`
	Title    string           `json:"title" binding:"required"`
	Message  string           `json:"message" binding:"required"`
}

// NotificationResponse represents the API response for a notification
type NotificationResponse struct {
	ID       uuid.UUID        `json:"id"`
	DriverID *uuid.UUID       `json:"driver_id,omitempty"`
	BinID    *uuid.UUID       `json:"bin_id,omitempty"`
	Type     NotificationType `json:"type"`
	Title    string           `json:"title"`
	Message  string           `json:"message"`
	IsRead   bool             `json:"is_read"`
	SentAt   time.Time        `json:"sent_at"`
	ReadAt   *time.Time       `json:"read_at,omitempty"`
}

// ToResponse converts Notification to NotificationResponse
func (n *Notification) ToResponse() *NotificationResponse {
	return &NotificationResponse{
		ID:       n.ID,
		DriverID: n.DriverID,
		BinID:    n.BinID,
		Type:     n.Type,
		Title:    n.Title,
		Message:  n.Message,
		IsRead:   n.IsRead,
		SentAt:   n.SentAt,
		ReadAt:   n.ReadAt,
	}
}
