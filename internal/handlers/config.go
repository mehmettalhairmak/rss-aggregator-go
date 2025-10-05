package handlers

import (
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
)

// Config holds the dependencies for all handlers
type Config struct {
	DB *database.Queries
}

// NewConfig creates a new handler config
// Constructor pattern - used to create Config instances
func NewConfig(db *database.Queries) *Config {
	return &Config{
		DB: db,
	}
}
