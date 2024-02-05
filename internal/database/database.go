package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open function opens a connection to the PostgreSQL database using the provided configuration
func Open(config PostgresConfig) (*sql.DB, error) {
	// Open a database connection using the "pgx" driver and the configuration string
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

// DefaultPostgresConfig returns a default configuration for a PostgreSQL database
func DefaultPostgresConfig() PostgresConfig {

	return PostgresConfig{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Database: os.Getenv("DATABASE"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}
}

// PostgresConfig represents the configuration parameters for a PostgreSQL database connection
type PostgresConfig struct {
	Host     string // Database host
	Port     string // Database port
	User     string // Database user
	Password string // Database password
	Database string // Database name
	SSLMode  string // SSL mode for the connection
}

// String method converts the PostgresConfig to a connection string
func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}
