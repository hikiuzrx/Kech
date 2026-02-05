package nats

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/smartwaste/backend/internal/config"
)

// Client represents a NATS client
type Client struct {
	conn *nats.Conn
	js   nats.JetStreamContext
	url  string
}

// NewClient creates a new NATS client
func NewClient(cfg *config.Config) *Client {
	// Check if NATS URL is configured, otherwise default
	url := "nats://localhost:4222"
	if cfg.MQTT.Broker != "" {
		// This is a bit of a hack since we're reusing the config struct which might not have NATS specific fields yet
		// ideally we'd add NATS config to the main config struct
		// reusing a hardcoded default or env var would be better if config isn't updated
	}

	return &Client{
		url: url,
	}
}

// Connect connects to the NATS server
func (c *Client) Connect() error {
	opts := []nats.Option{
		nats.Name("Go Backend Service"),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("Disconnected from NATS: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Reconnected to NATS: %s", nc.ConnectedUrl())
		}),
	}

	nc, err := nats.Connect(c.url, opts...)
	if err != nil {
		return err
	}
	c.conn = nc

	js, err := nc.JetStream()
	if err != nil {
		log.Printf("Warning: Failed to init JetStream: %v", err)
		// We might still be able to use basic NATS
	}
	c.js = js

	log.Println("Connected to NATS")
	return nil
}

// Subscribe subscribes to a subject
func (c *Client) Subscribe(subject string, handler func([]byte)) (*nats.Subscription, error) {
	return c.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
}

// Close closes the connection
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
