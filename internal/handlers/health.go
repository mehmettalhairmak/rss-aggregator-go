package handlers

import (
	"net/http"

	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

// HandlerReadiness checks if the server is ready
// Health check endpoint - verifies the server is running
// @Summary     Health check
// @Description Checks if the server is ready to handle requests
// @Tags        health
// @Accept      json
// @Produce     json
// @Success     200  {object}  map[string]interface{}
// @Router      /v1/ready [get]
func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	models.RespondWithJSON(w, http.StatusOK, struct{}{})
}

// HandlerErr is a test error handler
// Test endpoint for error handling
func HandlerErr(w http.ResponseWriter, r *http.Request) {
	models.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
}
