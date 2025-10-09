package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mehmettalhairmak/rss-aggregator/internal/auth"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// HandlerRegister handles new user registration (sign up).
//
// Flow:
//  1. Parse and validate request body (name, email, password)
//  2. Hash password using bcrypt for secure storage
//  3. Create user record in database
//  4. Generate JWT token for immediate authentication
//  5. Return user data and authentication token
//
// Security:
//   - Password is hashed with bcrypt (cost factor: 10)
//   - Email uniqueness enforced by database constraint
//   - Never stores plaintext passwords
//
// HTTP Status Codes:
//   - 201 Created: User successfully registered
//   - 400 Bad Request: Invalid input or duplicate email
//   - 500 Internal Server Error: Hash generation or token creation failed
func (cfg *Config) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// Validate required fields
	if params.Email == "" || params.Password == "" || params.Name == "" {
		models.RespondWithError(w, http.StatusBadRequest, "Name, email and password are required")
		return
	}

	// Hash password with bcrypt
	// Uses DefaultCost (10) which provides good security/performance balance
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user in database
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Email: sql.NullString{
			String: params.Email,
			Valid:  true,
		},
		PasswordHash: sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		},
	})
	if err != nil {
		// Database constraint errors (e.g., duplicate email) will be caught here
		models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create user: %v", err))
		return
	}

	// Generate JWT token for immediate authentication
	token, err := auth.GenerateJWT(user.ID, user.Email.String)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return user data and authentication token
	type response struct {
		User  models.User `json:"user"`
		Token string      `json:"token"`
	}

	models.RespondWithJSON(w, http.StatusCreated, response{
		User:  models.DatabaseUserToUser(user),
		Token: token,
	})
}

// HandlerLogin handles user authentication (sign in).
//
// Flow:
//  1. Parse and validate email and password from request
//  2. Retrieve user from database by email
//  3. Verify password using bcrypt comparison
//  4. Generate JWT token upon successful authentication
//  5. Return user data and authentication token
//
// Security:
//   - Uses constant-time password comparison (bcrypt)
//   - Returns generic error message to prevent user enumeration
//   - Implements secure password verification flow
//
// HTTP Status Codes:
//   - 200 OK: Authentication successful
//   - 400 Bad Request: Missing required fields
//   - 401 Unauthorized: Invalid credentials
//   - 500 Internal Server Error: Token generation failed
func (cfg *Config) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// Validate required fields
	if params.Email == "" || params.Password == "" {
		models.RespondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Find user by email
	user, err := cfg.DB.GetUserByEmail(r.Context(), sql.NullString{
		String: params.Email,
		Valid:  true,
	})
	if err != nil {
		// Return generic error to prevent user enumeration attacks
		models.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Verify password using constant-time comparison
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(params.Password))
	if err != nil {
		// Return same generic error for invalid password
		models.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Email.String)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return user data and authentication token
	type response struct {
		User  models.User `json:"user"`
		Token string      `json:"token"`
	}

	models.RespondWithJSON(w, http.StatusOK, response{
		User:  models.DatabaseUserToUser(user),
		Token: token,
	})
}
