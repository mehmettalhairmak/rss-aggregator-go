package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

// HandlerCreateFeedFollow creates a new feed follow relationship
// User starts following a feed
func (cfg *Config) HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	// Create feed follow relationship
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Create feed follow failed: %v", err))
		return
	}

	models.RespondWithJSON(w, http.StatusCreated, models.DatabaseFeedFollowToFeedFollow(feedFollow))
}

// HandlerGetFeedFollow returns all feeds the user follows
func (cfg *Config) HandlerGetFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Get feed follows failed: %v", err))
		return
	}

	models.RespondWithJSON(w, http.StatusOK, models.DatabaseAllFeedFollowToAllFeedFollow(feedFollows))
}

// HandlerDeleteFeedFollow deletes a feed follow relationship
// User stops following a feed
func (cfg *Config) HandlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	// Get feed_follow_id parameter from URL (via chi router)
	feedFollowIDString := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDString)
	if err != nil {
		models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid feed follow ID: %v", err))
		return
	}

	// Delete feed follow relationship
	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Delete feed follow failed: %v", err))
		return
	}

	models.RespondWithJSON(w, http.StatusNoContent, struct{}{})
}
