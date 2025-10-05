package handlers

import (
	"net/http"

	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

// HandlerReadiness checks if the server is ready
// Health check endpoint - verifies the server is running
func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	models.RespondWithJSON(w, http.StatusOK, struct{}{})
}

// HandlerErr is a test error handler
// Test endpoint for error handling
func HandlerErr(w http.ResponseWriter, r *http.Request) {
	models.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
}
