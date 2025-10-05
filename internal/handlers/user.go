package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

// HandlerCreateUser creates a new user
func (cfg *Config) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// parameters struct for parsing JSON from request body
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// Add new user to database
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	// Convert database model to API response model
	models.RespondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
}

// HandlerGetUser returns the authenticated user's information
func (cfg *Config) HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	models.RespondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}
