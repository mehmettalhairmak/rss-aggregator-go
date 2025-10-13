package models

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithError sends an error response in JSON format
func RespondWithError(w http.ResponseWriter, code int, message string) {
	// Log 5xx errors (server errors)
	if code > 499 {
		log.Println("Responding with 5XX error:", message)
	}

	// Error response struct'ı
	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errorResponse{Error: message})
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set response header
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	// Write JSON
	_, responseError := w.Write(data)
	if responseError != nil {
		log.Printf("Failed to write response: %v", responseError)
	}
}
