package middleware

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/mehmettalhairmak/rss-aggregator/internal/auth"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

// AuthedHandler is a handler function that requires authentication.
// Unlike a normal http.HandlerFunc, it takes database.User as a third parameter.
// Description: This is a custom handler type for protected endpoints.
type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

// Config holds dependencies for middleware.
type Config struct {
	DB *database.Queries
}

// NewConfig creates a new middleware config.
func NewConfig(db *database.Queries) *Config {
	return &Config{
		DB: db,
	}
}

// Auth wraps an authenticated handler with JWT authentication logic.
// Description: This is the JWT Middleware - it checks every request to protected endpoints.
// Flow:
// 1. Gets the JWT from the "Authorization: Bearer <token>" header.
// 2. Validates the token (checks signature and expiration).
// 3. Extracts the user_id from the token.
// 4. Finds the user in the database.
// 5. Passes the user information to the handler.
func (cfg *Config) Auth(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header.
		// Format: "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			models.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Strip the "Bearer " prefix and get the token.
		token, err := auth.GetBearerToken(authHeader)
		if err != nil {
			models.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid authorization header: %v", err))
			return
		}

		// Validate the JWT token.
		// This function checks the token's:
		// - Signature (is it signed with our secret key?).
		// - Expiration time (has it expired?).
		// - Claims (parses user_id, email, etc.).
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			models.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			return
		}

		// Find the user in the database with the user_id from the token.
		user, err := cfg.DB.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				models.RespondWithError(w, http.StatusUnauthorized, "User not found")
				return
			}
			models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
			return
		}

		// User found! Call the handler and pass the user information.
		// Now, user.ID, user.Email, etc., can be used inside the handler.
		handler(w, r, user)
	}
}
