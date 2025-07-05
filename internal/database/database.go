package database

import (
	"database/sql"
	"emergency-response-backend/internal/config"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func New(cfg config.DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable PostGIS extension
	if err := enablePostGIS(db); err != nil {
		return nil, fmt.Errorf("failed to enable PostGIS: %w", err)
	}

	return &DB{db}, nil
}

func enablePostGIS(db *sql.DB) error {
	extensions := []string{
		"CREATE EXTENSION IF NOT EXISTS postgis;",
		"CREATE EXTENSION IF NOT EXISTS postgis_topology;",
		"CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;",
		"CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder;",
	}

	for _, ext := range extensions {
		if _, err := db.Exec(ext); err != nil {
			return fmt.Errorf("failed to create extension: %w", err)
		}
	}

	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
