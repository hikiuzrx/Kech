package nats

const (
	// TopicShipmentCreated is published when a new shipment is created
	TopicShipmentCreated = "shipment.created"
	// TopicPriceConfirmed is published when a price is confirmed
	TopicPriceConfirmed = "shipment.price.confirmed"
	// TopicDriverAssigned is published when a driver is assigned
	TopicDriverAssigned = "shipment.driver.assigned"
	// TopicPickupStarted is published when pickup starts
	TopicPickupStarted = "shipment.pickup.started"
	// TopicPickupConfirmed is published when pickup is confirmed
	TopicPickupConfirmed = "shipment.pickup.confirmed"
	// TopicInTransit is published when shipment is in transit
	TopicInTransit = "shipment.in.transit"
	// TopicDelivered is published when shipment is delivered
	TopicDelivered = "shipment.delivered"
	// TopicCompleted is published when shipment is completed
	TopicCompleted = "shipment.completed"
	// TopicCancelled is published when shipment is cancelled
	TopicCancelled = "shipment.cancelled"
	// TopicDisputed is published when a dispute is raised
	TopicDisputed = "shipment.disputed"
	// TopicResolved is published when a dispute is resolved
	TopicResolved = "shipment.resolved"
	// TopicContractDeployed is published when a smart contract is deployed
	TopicContractDeployed = "shipment.contract.deployed"
)

// EventPayload represents the standard event payload structure
type EventPayload struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}
