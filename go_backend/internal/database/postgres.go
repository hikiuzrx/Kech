package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/smartwaste/backend/internal/config"
)

var db *sqlx.DB

// InitDB initializes the database connection
func InitDB(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := cfg.GetDSN()

	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// GetDB returns the database connection
func GetDB() *sqlx.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// RunMigrations applies database migrations
func RunMigrations(db *sqlx.DB) error {
	// In production, use a proper migration tool like golang-migrate
	// This is a simplified version for demonstration
	log.Println("Database migrations completed (placeholder)")
	return nil
}
