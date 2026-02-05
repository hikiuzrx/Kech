package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	MQTT     MQTTConfig
	Google   GoogleConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Mode string // debug, release, test
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// MQTTConfig holds MQTT broker configuration
type MQTTConfig struct {
	Broker   string
	Port     string
	ClientID string
	Username string
	Password string
}

// GoogleConfig holds Google API configuration
type GoogleConfig struct {
	MapsAPIKey string
}

var (
	cfg  *Config
	once sync.Once
)

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/app")

		// Set defaults
		viper.SetDefault("SERVER_PORT", "8080")
		viper.SetDefault("SERVER_MODE", "debug")
		viper.SetDefault("DB_HOST", "postgres")
		viper.SetDefault("DB_PORT", "5432")
		viper.SetDefault("DB_USER", "postgres")
		viper.SetDefault("DB_PASSWORD", "postgres")
		viper.SetDefault("DB_NAME", "smartwaste")
		viper.SetDefault("DB_SSLMODE", "disable")
		viper.SetDefault("MQTT_BROKER", "mosquitto")
		viper.SetDefault("MQTT_PORT", "1883")
		viper.SetDefault("MQTT_CLIENT_ID", "smartwaste-backend")
		viper.SetDefault("GOOGLE_MAPS_API_KEY", "")

		// Read from environment variables
		viper.AutomaticEnv()

		cfg = &Config{
			Server: ServerConfig{
				Port: viper.GetString("SERVER_PORT"),
				Mode: viper.GetString("SERVER_MODE"),
			},
			Database: DatabaseConfig{
				Host:     viper.GetString("DB_HOST"),
				Port:     viper.GetString("DB_PORT"),
				User:     viper.GetString("DB_USER"),
				Password: viper.GetString("DB_PASSWORD"),
				DBName:   viper.GetString("DB_NAME"),
				SSLMode:  viper.GetString("DB_SSLMODE"),
			},
			MQTT: MQTTConfig{
				Broker:   viper.GetString("MQTT_BROKER"),
				Port:     viper.GetString("MQTT_PORT"),
				ClientID: viper.GetString("MQTT_CLIENT_ID"),
				Username: viper.GetString("MQTT_USERNAME"),
				Password: viper.GetString("MQTT_PASSWORD"),
			},
			Google: GoogleConfig{
				MapsAPIKey: viper.GetString("GOOGLE_MAPS_API_KEY"),
			},
		}

		log.Printf("Configuration loaded: Server Port=%s, DB Host=%s, MQTT Broker=%s",
			cfg.Server.Port, cfg.Database.Host, cfg.MQTT.Broker)
	})

	return cfg
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	if cfg == nil {
		return LoadConfig()
	}
	return cfg
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}
