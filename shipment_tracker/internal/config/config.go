package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the service
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	NATS       NATSConfig
	Blockchain BlockchainConfig
	Service    ServiceConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NATSConfig holds NATS messaging configuration
type NATSConfig struct {
	URL       string
	ClusterID string
}

// BlockchainConfig holds blockchain configuration
type BlockchainConfig struct {
	RPCURL          string
	ChainID         int64
	PrivateKey      string
	ContractAddress string
}

// ServiceConfig holds service-specific configuration
type ServiceConfig struct {
	Name     string
	LogLevel string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults
	viper.SetDefault("SERVER_PORT", "8082")
	viper.SetDefault("SERVER_MODE", "debug")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "smartwaste_shipments")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("NATS_URL", "nats://localhost:4222")
	viper.SetDefault("NATS_CLUSTER_ID", "smartwaste-cluster")
	viper.SetDefault("BLOCKCHAIN_CHAIN_ID", 80001) // Polygon Mumbai
	viper.SetDefault("SERVICE_NAME", "shipment-tracker")
	viper.SetDefault("LOG_LEVEL", "debug")

	cfg := &Config{
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
		NATS: NATSConfig{
			URL:       viper.GetString("NATS_URL"),
			ClusterID: viper.GetString("NATS_CLUSTER_ID"),
		},
		Blockchain: BlockchainConfig{
			RPCURL:          viper.GetString("BLOCKCHAIN_RPC_URL"),
			ChainID:         viper.GetInt64("BLOCKCHAIN_CHAIN_ID"),
			PrivateKey:      viper.GetString("BLOCKCHAIN_PRIVATE_KEY"),
			ContractAddress: viper.GetString("CONTRACT_ADDRESS"),
		},
		Service: ServiceConfig{
			Name:     viper.GetString("SERVICE_NAME"),
			LogLevel: viper.GetString("LOG_LEVEL"),
		},
	}

	log.Printf("Configuration loaded for service: %s", cfg.Service.Name)
	return cfg
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}
