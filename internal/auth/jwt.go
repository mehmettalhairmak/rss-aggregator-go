package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// getJWTSecret retrieves the JWT secret key from environment variables.
// The application will panic if JWT_SECRET is not set, preventing insecure startup.
// In production, this should be a strong, randomly generated secret (min 32 bytes).
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable must be set")
	}
	return []byte(secret)
}

// CustomClaims represents the JWT payload structure.
// It embeds jwt.RegisteredClaims to include standard fields (exp, iat, sub, etc.)
// and adds custom fields specific to our application.
type CustomClaims struct {
	UserID               uuid.UUID `json:"user_id"` // Unique identifier for the user
	Email                string    `json:"email"`   // User's email address
	jwt.RegisteredClaims           // Standard JWT claims (expiry, issued at, etc.)
}

// GenerateJWT creates a signed JWT token for authenticated users.
// The token includes user identification data and is valid for 24 hours.
//
// Parameters:
//   - userID: Unique identifier of the user
//   - email: User's email address
//
// Returns:
//   - string: Base64-encoded JWT token
//   - error: Any error encountered during token generation
//
// Security:
//   - Uses HMAC-SHA256 (HS256) signing algorithm
//   - Token expires after 24 hours
//   - Secret key loaded from environment variable
func GenerateJWT(userID uuid.UUID, email string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken crates a secure, random string to be used as a refresh token.
// It generates 32bytes of random data and encodes it to a URL-safe base64 string.
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateJWT verifies and parses a JWT token string.
// It performs comprehensive validation including signature verification,
// expiration check, and algorithm validation.
//
// Parameters:
//   - tokenString: Raw JWT token string (without "Bearer " prefix)
//
// Returns:
//   - *CustomClaims: Parsed claims if token is valid
//   - error: Validation error if token is invalid, expired, or malformed
//
// Security considerations:
//   - Validates signing algorithm to prevent algorithm substitution attacks
//   - Checks token expiration automatically
//   - Verifies signature using secret key from environment
func ValidateJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method to prevent algorithm substitution attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims using type assertion
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}

	return claims, nil
}

// GetBearerToken extracts the JWT token from an HTTP Authorization header.
// Expected format: "Authorization: Bearer <token>"
//
// Parameters:
//   - authHeader: Value of the Authorization HTTP header
//
// Returns:
//   - string: Extracted JWT token
//   - error: Error if header is missing, empty, or malformed
//
// Example:
//
//	token, err := GetBearerToken("Bearer eyJhbGciOiJIUzI1...")
func GetBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	const prefix = "Bearer "
	if len(authHeader) < len(prefix) {
		return "", errors.New("authorization header format is invalid")
	}

	if authHeader[:len(prefix)] != prefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	token := authHeader[len(prefix):]
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}

func HashRefreshToken(tokenString string) string {
	hasher := sha256.New()
	hasher.Write([]byte(tokenString))

	return hex.EncodeToString(hasher.Sum(nil))
}
