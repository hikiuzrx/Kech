package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/smartwaste/iot-sensor/pkg/config"
	"github.com/smartwaste/iot-sensor/pkg/sensor"
)

// Payload represents the data sent to the backend
type Payload struct {
	BinID     string `json:"bin_id"`
	FillLevel int    `json:"fill_level"`
	Battery   int    `json:"battery_level,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()
	log.Printf("Starting IoT Sensor Service for Bin: %s", cfg.BinID)

	// 2. Setup Sensor
	var s sensor.Sensor

	if cfg.Simulation {
		log.Println("Mode: Simulation")
		s = sensor.NewSimulator(cfg.BinHeightCm)
	} else {
		// Hardware initialization would go here (requires build tags for TinyGo/Hardware)
		// For now we default to simulator if not strictly configured for hardware
		// In a real RPi deployment, we'd initialize the GPIO here
		log.Println("Mode: Hardware (Falling back to simulator for this cross-platform build example)")
		// In actual tinygo code:
		// trig := machine.GPIO23
		// echo := machine.GPIO24
		// s = sensor.NewHCSR04(trig, echo)
		s = sensor.NewSimulator(cfg.BinHeightCm)
	}
	defer s.Close()

	// 3. Setup MQTT
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTTBroker)
	opts.SetClientID("iot-sensor-" + cfg.BinID)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT: %v", token.Error())
	}
	log.Printf("Connected to MQTT Broker: %s", cfg.MQTTBroker)

	// 4. Main Loop
	topic := fmt.Sprintf("bins/%s/status", cfg.BinID)

	for {
		// Read Sensor
		distanceCm, err := s.ReadDistance()
		if err != nil {
			log.Printf("Error reading sensor: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Calculate Fill Level
		// Distance = Gap from top to waste.
		// Fill Height = BinHeight - Distance
		// Fill % = (Fill Height / BinHeight) * 100
		fillHeight := cfg.BinHeightCm - distanceCm
		if fillHeight < 0 {
			fillHeight = 0
		}

		fillLevel := int((fillHeight / cfg.BinHeightCm) * 100)
		if fillLevel > 100 {
			fillLevel = 100
		} else if fillLevel < 0 {
			fillLevel = 0
		}

		// Create Payload
		payload := Payload{
			BinID:     cfg.BinID,
			FillLevel: fillLevel,
			Timestamp: time.Now().Unix(),
		}

		data, _ := json.Marshal(payload)

		// Publish
		token := client.Publish(topic, 0, false, data)
		token.Wait()

		log.Printf("Published to %s: %s (Distance: %.1fcm)", topic, string(data), distanceCm)

		time.Sleep(cfg.ReadInterval)
	}
}
