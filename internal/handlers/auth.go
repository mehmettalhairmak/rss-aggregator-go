package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
//
// @Summary     Register a new user
// @Description Creates a new user account with email and password
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       user  body      object  true  "User registration data" schema(parameters)
// @Success     201   {object}  object  "User registered successfully"
// @Failure     400   {object}  object  "Invalid input or duplicate email"
// @Failure     500   {object}  object  "Server error"
// @Router      /v1/auth/register [post]
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
	accessToken, err := auth.GenerateJWT(user.ID, user.Email.String)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	_, errSaveRefreshTokenDb := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: auth.HashRefreshToken(refreshToken),
		ExpiresAt: time.Now().Add(24 * time.Hour * 7).UTC(),
		CreatedAt: time.Now().UTC(),
	})
	if errSaveRefreshTokenDb != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	// Return user data and authentication token
	type response struct {
		User         models.User `json:"user"`
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		ExpiresIn    int64       `json:"expires_in"`
	}

	models.RespondWithJSON(w, http.StatusCreated, response{
		User:         models.DatabaseUserToUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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
//
// @Summary     Login user
// @Description Authenticate user with email and password
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       credentials  body      object  true  "Login credentials"
// @Success     200          {object}  object  "Login successful"
// @Failure     400          {object}  object  "Invalid input"
// @Failure     401          {object}  object  "Invalid credentials"
// @Router      /v1/auth/login [post]
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

	refreshToken, errRefreshToken := auth.GenerateRefreshToken()
	if errRefreshToken != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	errDeleteGenerateRefreshToken := cfg.deleteAndGenerateRefreshTokenFromDB(r.Context(), &user, refreshToken)
	if errDeleteGenerateRefreshToken != nil {
		models.RespondWithError(w, http.StatusInternalServerError, errDeleteGenerateRefreshToken.Error())
		return
	}

	// Return user data and authentication token
	type response struct {
		User         models.User `json:"user"`
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
	}

	models.RespondWithJSON(w, http.StatusOK, response{
		User:         models.DatabaseUserToUser(user),
		AccessToken:  token,
		RefreshToken: refreshToken,
	})
}

// @Summary     Logout user
// @Description Logout user and invalidate refresh token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Success     200  {object}  object  "Logout successful"
// @Failure     500  {object}  object  "Server error"
// @Router      /v1/auth/logout [get]
func (cfg *Config) HandlerLogout(w http.ResponseWriter, r *http.Request, user database.User) {
	err := cfg.DB.DeleteRefreshToken(r.Context(), user.ID)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to delete refresh token")
		return
	}

	models.RespondWithJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "Successfully logged out",
	})
}

// HandlerRefreshToken handles issuing a new JWT access token using a valid refresh token.
//
// Flow:
//  1. Parse and validate refresh token from request body
//  2. Hash the provided refresh token for secure comparison
//  3. Retrieve the corresponding refresh token record from the database
//  4. Check if the refresh token is expired
//  5. Retrieve the associated user from the database
//  6. Generate a new JWT access token
//  7. Generate a new refresh token and update the database record
//  8. Return the new access token, refresh token, and user data
//
// Security:
//   - Refresh tokens are hashed before storage and comparison to prevent leakage
//   - Expired refresh tokens are rejected to prevent reuse
//   - New refresh tokens are generated upon each use to limit lifespan
//
// HTTP Status Codes:
//   - 200 OK: New tokens successfully issued
//   - 400 Bad Request: Missing or invalid refresh token, or expired token
//   - 401 Unauthorized: Invalid token
//   - 500 Internal Server Error: Database or token generation failure
//
// @Summary     Refresh access token
// @Description Get new access token using refresh token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       refresh_token  body      object  true  "Refresh token"
// @Success     200            {object}  object  "New tokens issued"
// @Failure     400            {object}  object  "Invalid or expired token"
// @Router      /v1/auth/refresh [post]
func (cfg *Config) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		models.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.RefreshToken == "" {
		models.RespondWithError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	hashedRefreshTokenPayload := auth.HashRefreshToken(params.RefreshToken)
	if hashedRefreshTokenPayload == "" {
		models.RespondWithError(w, http.StatusBadRequest, "Refresh token is required")
	}

	refreshTokenObject, errGetRefreshTokenFromDb := cfg.DB.GetRefreshTokenByHash(r.Context(), hashedRefreshTokenPayload)
	if errGetRefreshTokenFromDb != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to get refresh token from DB")
		return
	}

	if time.Now().UTC().After(refreshTokenObject.ExpiresAt) {
		models.RespondWithError(w, http.StatusBadRequest, "Refresh token is expired")
		return
	}

	user, errFindUser := cfg.DB.GetUserByID(r.Context(), refreshTokenObject.UserID)
	if errFindUser != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to find user")
		return
	}

	accessToken, err := auth.GenerateJWT(user.ID, user.Email.String)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	errGenerateRefToken := cfg.deleteAndGenerateRefreshTokenFromDB(r.Context(), &user, refreshToken)
	if errGenerateRefToken != nil {
		models.RespondWithError(w, http.StatusInternalServerError, errGenerateRefToken.Error())
		return
	}

	type response struct {
		User         models.User `json:"user"`
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
	}

	models.RespondWithJSON(w, http.StatusOK, response{
		User:         models.DatabaseUserToUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}

// deleteAndGenerateRefreshTokenFromDB deletes any existing refresh token for the user
// and creates a new refresh token record in the database within a transaction.
//
// Parameters:
//   - context: The context for database operations
//   - user: The user for whom the refresh token is being managed
//   - refreshTokenString: The new refresh token string to be hashed and stored
//
// Returns:
//   - error: Any error encountered during the process, or nil if successful
//
// Transaction Management:
//   - Begins a new database transaction
//   - Deletes existing refresh token for the user
//   - Inserts the new refresh token record
//   - Commits the transaction if all operations succeed
//   - Rolls back the transaction in case of any errors
func (cfg *Config) deleteAndGenerateRefreshTokenFromDB(context context.Context, user *database.User, refreshTokenString string) error {
	tx, errorTx := cfg.DBConn.BeginTx(context, nil)
	if errorTx != nil {
		return fmt.Errorf("Failed to start transaction: %v", errorTx)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("transaction rollback failed: %v", err)
		}
	}()

	qtx := cfg.DB.WithTx(tx)

	errDeleteRefreshTokenDb := qtx.DeleteRefreshToken(context, user.ID)
	if errDeleteRefreshTokenDb != nil {
		return fmt.Errorf("Failed to delete refresh token: %v", errDeleteRefreshTokenDb)
	}

	_, errSaveRefreshTokenDb := qtx.CreateRefreshToken(context, database.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: auth.HashRefreshToken(refreshTokenString),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).UTC(),
		CreatedAt: time.Now().UTC(),
	})
	if errSaveRefreshTokenDb != nil {
		return fmt.Errorf("Failed to save refresh token: %v", errSaveRefreshTokenDb)
	}

	if errTxCommit := tx.Commit(); errTxCommit != nil {
		return fmt.Errorf("Failed to commit transaction: %v", errTxCommit)
	}

	return nil
}
