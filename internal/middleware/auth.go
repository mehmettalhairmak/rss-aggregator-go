package middleware

import (
	"fmt"
	"net/http"

	"github.com/mehmettalhairmak/rss-aggregator/internal/auth"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

// AuthedHandler is a handler function that requires authentication
// Unlike normal http.HandlerFunc, it takes database.User as a third parameter
type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

// Config holds dependencies for middleware
type Config struct {
	DB *database.Queries
}

// NewConfig creates a new middleware config
func NewConfig(db *database.Queries) *Config {
	return &Config{
		DB: db,
	}
}

// Auth wraps an authenticated handler with authentication logic
// Validates API Key and passes user information to the handler
func (cfg *Config) Auth(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get API Key from HTTP header
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			models.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error: %v", err))
			return
		}

		// Fetch user from database using API Key
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		// If user found, call handler and pass user as parameter
		handler(w, r, user)
	}
}
