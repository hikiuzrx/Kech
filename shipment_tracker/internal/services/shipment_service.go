package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/smartwaste/shipment-tracker/internal/models"
	"github.com/smartwaste/shipment-tracker/internal/nats"
	"github.com/smartwaste/shipment-tracker/internal/repository"
)

// ShipmentService handles shipment business logic
type ShipmentService struct {
	shipmentRepo   *repository.ShipmentRepository
	transitionRepo *repository.TransitionRepository
	natsClient     *nats.Client
}

// NewShipmentService creates a new ShipmentService
func NewShipmentService(
	shipmentRepo *repository.ShipmentRepository,
	transitionRepo *repository.TransitionRepository,
	natsClient *nats.Client,
) *ShipmentService {
	return &ShipmentService{
		shipmentRepo:   shipmentRepo,
		transitionRepo: transitionRepo,
		natsClient:     natsClient,
	}
}

// CreateShipment creates a new shipment and logs the transition
func (s *ShipmentService) CreateShipment(req *models.CreateShipmentRequest) (*models.Shipment, error) {
	id := uuid.New()
	now := time.Now()

	shipment := &models.Shipment{
		ID:                id,
		UserID:            req.UserID,
		CollectionID:      req.CollectionID,
		WasteType:         req.WasteType,
		EstimatedWeightKg: req.EstimatedWeightKg,
		PriceOffered:      req.PriceOffered,
		Status:            models.StatusCreated,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if req.PickupLocation != nil {
		shipment.PickupLatitude = &req.PickupLocation.Latitude
		shipment.PickupLongitude = &req.PickupLocation.Longitude
		shipment.PickupAddress = &req.PickupLocation.Address
	}

	if req.DropoffLocation != nil {
		shipment.DropoffLatitude = &req.DropoffLocation.Latitude
		shipment.DropoffLongitude = &req.DropoffLocation.Longitude
		shipment.DropoffAddress = &req.DropoffLocation.Address
	}

	if req.Notes != nil {
		shipment.Notes = req.Notes
	}

	// 1. Save shipment to DB
	if err := s.shipmentRepo.Create(shipment); err != nil {
		return nil, err
	}

	// 2. Create initial state transition
	transition := &models.StateTransition{
		ID:              uuid.New(),
		ShipmentID:      id,
		FromStatus:      nil, // Initial state
		ToStatus:        models.StatusCreated,
		TriggeredBy:     req.UserID,
		TriggeredByRole: "user",
		CreatedAt:       now,
	}
	if err := s.transitionRepo.Create(transition); err != nil {
		// Log error but don't fail, we successfully created the shipment
		// In production, might want transactional integrity here
	}

	// 3. Publish event to NATS
	s.publishEvent(nats.TopicShipmentCreated, shipment)

	return shipment, nil
}

// GetShipment retrieves a shipment by ID
func (s *ShipmentService) GetShipment(id uuid.UUID) (*models.Shipment, error) {
	return s.shipmentRepo.GetByID(id)
}

// AssignDriver assigns a driver to the shipment
func (s *ShipmentService) AssignDriver(shipmentID uuid.UUID, driverID uuid.UUID) error {
	shipment, err := s.shipmentRepo.GetByID(shipmentID)
	if err != nil {
		return err
	}
	if shipment == nil {
		return fmt.Errorf("shipment not found")
	}

	// Validate transition
	if !shipment.CanTransitionTo(models.StatusDriverAssigned) {
		return fmt.Errorf("cannot transition from %s to %s", shipment.Status, models.StatusDriverAssigned)
	}

	// Update DB
	if err := s.shipmentRepo.AssignDriver(shipmentID, driverID); err != nil {
		return err
	}

	// Record transition
	now := time.Now()
	fromStatus := shipment.Status
	transition := &models.StateTransition{
		ID:              uuid.New(),
		ShipmentID:      shipmentID,
		FromStatus:      &fromStatus,
		ToStatus:        models.StatusDriverAssigned,
		TriggeredBy:     driverID, // Assuming driver requests assignment or system does
		TriggeredByRole: "driver", // or system
		CreatedAt:       now,
	}
	s.transitionRepo.Create(transition)

	// Publish event
	s.publishEvent(nats.TopicDriverAssigned, map[string]interface{}{
		"shipment_id": shipmentID,
		"driver_id":   driverID,
	})

	return nil
}

// Helper to update shipment status and record transition
func (s *ShipmentService) updateStatusAndRecord(
	shipment *models.Shipment,
	newStatus models.ShipmentStatus,
	triggeredBy uuid.UUID,
	role string,
	proofHash *string,
	signature *string,
	metadata map[string]interface{},
) error {
	// 1. Validate Transition
	if !shipment.CanTransitionTo(newStatus) {
		return fmt.Errorf("invalid transition from %s to %s", shipment.Status, newStatus)
	}

	// 2. Update Shipment Status
	if err := s.shipmentRepo.UpdateStatus(shipment.ID, newStatus); err != nil {
		return err
	}

	// 3. Record Transition
	mdBytes, _ := json.Marshal(metadata)
	fromStatus := shipment.Status
	transition := &models.StateTransition{
		ID:              uuid.New(),
		ShipmentID:      shipment.ID,
		FromStatus:      &fromStatus,
		ToStatus:        newStatus,
		TriggeredBy:     triggeredBy,
		TriggeredByRole: role,
		ProofHash:       proofHash,
		Signature:       signature,
		Metadata:        json.RawMessage(mdBytes),
		CreatedAt:       time.Now(),
	}
	if err := s.transitionRepo.Create(transition); err != nil {
		return err
	}

	// 4. Publish Event
	topic := s.getTopicForStatus(newStatus)
	s.publishEvent(topic, map[string]interface{}{
		"shipment_id": shipment.ID,
		"status":      newStatus,
		"updated_by":  triggeredBy,
	})

	return nil
}

func (s *ShipmentService) getTopicForStatus(status models.ShipmentStatus) string {
	switch status {
	case models.StatusPriceConfirmed:
		return nats.TopicPriceConfirmed
	case models.StatusPickupStarted:
		return nats.TopicPickupStarted
	case models.StatusInTransit:
		return nats.TopicInTransit
	case models.StatusDelivered:
		return nats.TopicDelivered
	case models.StatusCompleted:
		return nats.TopicCompleted
	case models.StatusCancelled:
		return nats.TopicCancelled
	case models.StatusDisputed:
		return nats.TopicDisputed
	case models.StatusResolved:
		return nats.TopicResolved
	default:
		return "shipment.status.changed"
	}
}

func (s *ShipmentService) publishEvent(topic string, data interface{}) {
	payload := nats.EventPayload{
		EventID:   uuid.New().String(),
		EventType: topic,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
	if err := s.natsClient.Publish(topic, payload); err != nil {
		// Log error
		fmt.Printf("Failed to publish event %s: %v\n", topic, err)
	}
}
