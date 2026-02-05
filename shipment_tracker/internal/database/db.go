package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/smartwaste/shipment-tracker/internal/config"
)

var db *sqlx.DB

// InitDB initializes the database connection
func InitDB(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	var err error
	db, err = sqlx.Connect("postgres", cfg.GetDSN())
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connection established")
	return db, nil
}

// GetDB returns the database connection
func GetDB() *sqlx.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}
}
