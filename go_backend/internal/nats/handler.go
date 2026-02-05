package nats

import (
	"encoding/json"
	"log"

	"github.com/smartwaste/backend/internal/services"
)

// EventPayload matches the payload structure from shipment_tracker
type EventPayload struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// EventHandler handles incoming NATS events
type EventHandler struct {
	notificationSvc *services.NotificationService
}

// NewEventHandler creates a new event handler
func NewEventHandler(notificationSvc *services.NotificationService) *EventHandler {
	return &EventHandler{
		notificationSvc: notificationSvc,
	}
}

// HandleShipmentCreated handles shipment creation events
func (h *EventHandler) HandleShipmentCreated(data []byte) {
	var payload EventPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("Error unmarshalling shipment created event: %v", err)
		return
	}
	log.Printf("Received Shipment Created Event: %v", payload.EventID)
	// TODO: Notify admin or update local state
}

// HandlePriceConfirmed handles price confirmation events
func (h *EventHandler) HandlePriceConfirmed(data []byte) {
	var payload EventPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("Error unmarshalling price confirmed event: %v", err)
		return
	}
	log.Printf("Received Price Confirmed Event: %v", payload.EventID)
	// Example: Notify driver that price is confirmed and they can proceed
}

// HandlePickupStarted handles pickup started events
func (h *EventHandler) HandlePickupStarted(data []byte) {
	var payload EventPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("Error unmarshalling pickup started event: %v", err)
		return
	}
	log.Printf("Received Pickup Started Event: %v", payload.EventID)
	// Notify user that driver has started pickup
}

// HandleDeliveryCompleted handles delivery completion events
func (h *EventHandler) HandleDeliveryCompleted(data []byte) {
	var payload EventPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("Error unmarshalling delivery completed event: %v", err)
		return
	}
	log.Printf("Received Delivery Completed Event: %v", payload.EventID)
	// Process payment, update user stats, etc.
}
