package handlers

import (
	"database/sql"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
)

// Config holds the dependencies for all handlers
type Config struct {
	DB     *database.Queries
	DBConn *sql.DB
}

// NewConfig creates a new handler config
// Constructor pattern - used to create Config instances
func NewConfig(queries *database.Queries, db *sql.DB) *Config {
	return &Config{
		DB:     queries,
		DBConn: db,
	}
}
