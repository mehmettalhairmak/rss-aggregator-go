package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
	"github.com/mmcdole/gofeed"
)

// HandlerCreateFeed creates a new RSS feed
// @Summary     Create RSS feed
// @Description Creates a new RSS feed and automatically follows it
// @Tags        feeds
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Param       feed  body      object  true  "Feed data"
// @Success     201   {object}  object  "Feed created"
// @Failure     400   {object}  object  "Invalid input"
// @Failure     500   {object}  object  "Server error"
// @Router      /v1/feed [post]
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

	gf := gofeed.NewParser()
	parsedFeed, errParseUrl := gf.ParseURL(params.URL)
	if errParseUrl != nil {
		models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request URL: %v", errParseUrl))
		return
	}

	// Extract metadata from parsed feed
	description := parsedFeed.Description
	logoUrl := ""
	if parsedFeed.Image != nil {
		logoUrl = parsedFeed.Image.URL
	}

	tx, errTx := cfg.DBConn.BeginTx(r.Context(), nil)
	if errTx != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error starting transaction: %v", errTx))
		return
	}

	defer tx.Rollback()

	qtx := cfg.DB.WithTx(tx)

	// Add new feed to database with metadata
	var descriptionNullStr, logoUrlNullStr sql.NullString

	if description != "" {
		descriptionNullStr = sql.NullString{String: description, Valid: true}
	}
	if logoUrl != "" {
		logoUrlNullStr = sql.NullString{String: logoUrl, Valid: true}
	}

	feed, errCreateFeed := qtx.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:          uuid.New(),
		Name:        params.Name,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Url:         params.URL,
		UserID:      user.ID,
		Description: descriptionNullStr,
		LogoUrl:     logoUrlNullStr,
		Priority:    3, // Default priority
	})
	if errCreateFeed != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Create Feed failed: %v", errCreateFeed))
		return
	}

	feedFollow, errCreateFeedFollow := qtx.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if errCreateFeedFollow != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Create Feed Follow failed: %v", errCreateFeedFollow))
		return
	}

	if err := tx.Commit(); err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error committing transaction: %v", err))
		return
	}

	type response struct {
		Feed       models.Feed       `json:"feed"`
		FeedFollow models.FeedFollow `json:"feed_follow"`
	}

	models.RespondWithJSON(w, http.StatusCreated, response{
		Feed:       models.DatabaseFeedToFeed(feed),
		FeedFollow: models.DatabaseFeedFollowToFeedFollow(feedFollow),
	})
}

// HandlerGetFeed returns all feeds
// @Summary     Get all feeds
// @Description Get a list of all RSS feeds
// @Tags        feeds
// @Accept      json
// @Produce     json
// @Success     200  {object}  object  "List of feeds"
// @Router      /v1/feed [get]
func (cfg *Config) HandlerGetFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Get Feed failed: %v", err))
		return
	}

	models.RespondWithJSON(w, http.StatusOK, models.DatabaseAllFeedToAllFeed(feeds))
}
