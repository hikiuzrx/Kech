package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	MQTTBroker   string
	MQTTPort     string
	BinID        string
	BinHeightCm  float64
	ReadInterval time.Duration
	Simulation   bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	return Config{
		MQTTBroker:   getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		BinID:        getEnv("BIN_ID", "bin-default-01"),
		BinHeightCm:  getEnvFloat("BIN_HEIGHT_CM", 100.0),
		ReadInterval: time.Duration(getEnvInt("READ_INTERVAL_SECONDS", 10)) * time.Second,
		Simulation:   getEnvBool("SIMULATION_MODE", true), // Default to simulation if no hardware
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if value, ok := os.LookupEnv(key); ok {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}
