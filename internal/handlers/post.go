package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

type postsResponse struct {
	Posts      []models.Post `json:"posts"`
	NextCursor string        `json:"next_cursor"`
}

func (cfg *Config) HandlerGetUserPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.URL.Query().Get("limit")
	limit := 20

	if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = parsedLimit
	}

	if limit > 100 {
		limit = 100
	}

	cursor := time.Now().UTC()
	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		parsedCursor, err := time.Parse(time.RFC3339, cursorStr)
		if err != nil {
			models.RespondWithError(w, http.StatusBadRequest, "Invalid cursor format")
			return
		}
		cursor = parsedCursor
	}

	posts, errGetPosts := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID:      user.ID,
		Limit:       int32(limit),
		PublishedAt: cursor,
	})

	if errGetPosts != nil {
		models.RespondWithError(w, http.StatusBadRequest, errGetPosts.Error())
		return
	}

	nextCursor := ""
	if len(posts) > 0 {
		lastPost := posts[len(posts)-1]
		nextCursor = lastPost.PublishedAt.Format(time.RFC3339)
	}

	response := postsResponse{
		Posts:      models.DatabaseAllPostToAllPost(posts),
		NextCursor: nextCursor,
	}

	models.RespondWithJSON(w, http.StatusOK, response)
}
