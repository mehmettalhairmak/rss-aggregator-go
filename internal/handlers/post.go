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

// @Summary     Get user posts
// @Description Get posts from all followed feeds with cursor-based pagination
// @Tags        posts
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Param       limit   query     int     false  "Number of posts to return (max 100)"  default(20)
// @Param       cursor  query     string  false  "Cursor for pagination (RFC3339 timestamp)"
// @Success     200     {object}  object  "List of posts"
// @Failure     400     {object}  object  "Invalid parameters"
// @Router      /v1/posts [get]
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
