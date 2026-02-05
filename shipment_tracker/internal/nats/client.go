package nats

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/smartwaste/shipment-tracker/internal/config"
)

// Client represents a NATS client
type Client struct {
	conn *nats.Conn
	js   nats.JetStreamContext
	url  string
}

// NewClient creates a new NATS client
func NewClient(cfg *config.NATSConfig) *Client {
	return &Client{
		url: cfg.URL,
	}
}

// Connect connects to the NATS server and initializes JetStream
func (c *Client) Connect() error {
	opts := []nats.Option{
		nats.Name("Shipment Tracker Service"),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("Disconnected from NATS: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Reconnected to NATS: %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("NATS connection closed")
		}),
	}

	nc, err := nats.Connect(c.url, opts...)
	if err != nil {
		return err
	}
	c.conn = nc

	// Create JetStream Context
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return err
	}
	c.js = js

	log.Println("Connected to NATS and JetStream initialized")

	// Ensure streams exist
	if err := c.createStreams(); err != nil {
		log.Printf("Warning: Could not create streams: %v", err)
	}

	return nil
}

// Close closes the NATS connection
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// Publish publishes a message to a subject
func (c *Client) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Use JetStream publish for persistence if needed, or standard publish for fire-and-forget
	// Here we use standard publish for simplicity, but could switch to js.Publish
	return c.conn.Publish(subject, payload)
}

// Subscribe subscribes to a subject
func (c *Client) Subscribe(subject string, handler func([]byte)) (*nats.Subscription, error) {
	return c.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
}

// createStreams creates necessary JetStream streams
func (c *Client) createStreams() error {
	// Define Shipment Stream
	_, err := c.js.AddStream(&nats.StreamConfig{
		Name:     "SHIPMENTS",
		Subjects: []string{"shipment.>"},
		Storage:  nats.FileStorage,
	})
	if err != nil {
		return err
	}
	return nil
}
