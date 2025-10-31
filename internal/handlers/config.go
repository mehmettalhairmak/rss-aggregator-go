package handlers

import (
	"database/sql"

	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/realtime"
	"github.com/rs/zerolog"
)

// Config holds the dependencies for all handlers
type Config struct {
	DB     *database.Queries
	DBConn *sql.DB
	Logger zerolog.Logger
	Hub    *realtime.Hub
}

// NewConfig creates a new handler config
// Constructor pattern - used to create Config instances
func NewConfig(queries *database.Queries, db *sql.DB, logger zerolog.Logger, hub *realtime.Hub) *Config {
	return &Config{
		DB:     queries,
		DBConn: db,
		Logger: logger,
		Hub:    hub,
	}
}
