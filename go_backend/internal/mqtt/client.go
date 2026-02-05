package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/smartwaste/backend/internal/config"
	"github.com/smartwaste/backend/internal/models"
	"github.com/smartwaste/backend/internal/repository"
	"github.com/smartwaste/backend/internal/services"
)

// Client wraps the MQTT client
type Client struct {
	client              pahomqtt.Client
	binRepo             *repository.BinRepository
	notificationService *services.NotificationService
	fillLevelThreshold  int
}

// NewClient creates a new MQTT client
func NewClient(cfg *config.MQTTConfig, binRepo *repository.BinRepository, notificationService *services.NotificationService) *Client {
	opts := pahomqtt.NewClientOptions()
	broker := fmt.Sprintf("tcp://%s:%s", cfg.Broker, cfg.Port)
	opts.AddBroker(broker)
	opts.SetClientID(cfg.ClientID)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	// Set connection options
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetCleanSession(true)

	mqttClient := &Client{
		binRepo:             binRepo,
		notificationService: notificationService,
		fillLevelThreshold:  90, // Trigger notification when fill level exceeds 90%
	}

	// Set callbacks
	opts.SetOnConnectHandler(mqttClient.onConnect)
	opts.SetConnectionLostHandler(mqttClient.onConnectionLost)
	opts.SetDefaultPublishHandler(mqttClient.messageHandler)

	mqttClient.client = pahomqtt.NewClient(opts)

	return mqttClient
}

// Connect establishes connection to the MQTT broker
func (c *Client) Connect() error {
	token := c.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}
	log.Println("Connected to MQTT broker successfully")
	return nil
}

// Disconnect closes the MQTT connection
func (c *Client) Disconnect() {
	c.client.Disconnect(250)
	log.Println("Disconnected from MQTT broker")
}

// Subscribe subscribes to the bin status topic
func (c *Client) Subscribe() error {
	// Subscribe to bin status updates from all bins
	// Topic pattern: bins/+/status where + is a wildcard for bin_id
	topic := "bins/+/status"
	token := c.client.Subscribe(topic, 1, c.binStatusHandler)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}
	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

// onConnect is called when the client connects to the broker
func (c *Client) onConnect(client pahomqtt.Client) {
	log.Println("MQTT client connected")
	// Resubscribe after reconnection
	if err := c.Subscribe(); err != nil {
		log.Printf("Failed to resubscribe after reconnection: %v", err)
	}
}

// onConnectionLost is called when the connection to the broker is lost
func (c *Client) onConnectionLost(client pahomqtt.Client, err error) {
	log.Printf("MQTT connection lost: %v", err)
}

// messageHandler is the default message handler
func (c *Client) messageHandler(client pahomqtt.Client, msg pahomqtt.Message) {
	log.Printf("Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))
}

// binStatusHandler processes bin status updates
func (c *Client) binStatusHandler(client pahomqtt.Client, msg pahomqtt.Message) {
	// Process message in a goroutine for concurrent handling
	go c.processBinStatus(msg.Payload())
}

// processBinStatus handles the bin status update logic
func (c *Client) processBinStatus(payload []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse JSON payload
	var status models.BinStatusUpdate
	if err := json.Unmarshal(payload, &status); err != nil {
		log.Printf("Failed to parse bin status payload: %v", err)
		return
	}

	log.Printf("Processing bin status update: BinID=%s, FillLevel=%d%%", status.BinID, status.FillLevel)

	// Validate fill level
	if status.FillLevel < 0 || status.FillLevel > 100 {
		log.Printf("Invalid fill level %d for bin %s", status.FillLevel, status.BinID)
		return
	}

	// Update bin fill level in database
	if err := c.binRepo.UpdateFillLevel(ctx, status.BinID, status.FillLevel); err != nil {
		log.Printf("Failed to update bin fill level: %v", err)
		return
	}

	// Check if bin needs collection (threshold exceeded)
	if status.FillLevel >= c.fillLevelThreshold {
		log.Printf("Bin %s fill level (%d%%) exceeds threshold (%d%%), triggering notification",
			status.BinID, status.FillLevel, c.fillLevelThreshold)

		// Get bin details
		bin, err := c.binRepo.GetByDeviceID(ctx, status.BinID)
		if err != nil || bin == nil {
			log.Printf("Failed to get bin details for notification: %v", err)
			return
		}

		// Trigger notification to nearest driver
		go c.notificationService.NotifyNearestDriver(ctx, bin)
	}
}

// Publish publishes a message to a topic
func (c *Client) Publish(topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	token := c.client.Publish(topic, 1, false, data)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %w", token.Error())
	}

	return nil
}

// IsConnected returns true if the client is connected
func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}
