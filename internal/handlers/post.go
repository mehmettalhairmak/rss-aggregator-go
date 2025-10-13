package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

type postsResponse struct {
	Posts      []models.Post  `json:"posts"`
	Pagination PaginationData `json:"pagination"`
}

type PaginationData struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
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

	pageStr := r.URL.Query().Get("page")
	page := 1

	if parsedPage, err := strconv.Atoi(pageStr); err == nil {
		if parsedPage > 1 {
			page = parsedPage
		}
	}

	offset := (page - 1) * limit

	totalItems, err := cfg.DB.CountPostsForUser(r.Context(), user.ID)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Couldn't count posts for user")
		return
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(limit)))
	}

	posts, errGetPosts := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if errGetPosts != nil {
		models.RespondWithError(w, http.StatusBadRequest, errGetPosts.Error())
		return
	}

	response := postsResponse{
		Posts: models.DatabaseAllPostToAllPost(posts),
		Pagination: PaginationData{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: page,
			PerPage:     limit,
		},
	}

	models.RespondWithJSON(w, http.StatusOK, response)
}
