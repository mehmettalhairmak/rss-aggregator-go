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

// HandlerCreateFeed creates a new RSS feed
func (cfg *Config) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	// Add new feed to database
	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Create Feed failed: %v", err))
		return
	}

	models.RespondWithJSON(w, http.StatusCreated, models.DatabaseFeedToFeed(feed))
}

// HandlerGetFeed returns all feeds
func (cfg *Config) HandlerGetFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Get Feed failed: %v", err))
		return
	}

	models.RespondWithJSON(w, http.StatusOK, models.DatabaseAllFeedToAllFeed(feeds))
}
