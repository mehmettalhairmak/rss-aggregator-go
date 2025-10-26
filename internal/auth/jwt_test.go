package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func init() {
	// Set JWT_SECRET for tests
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-only")
}

func TestGenerateJWT_ValidUser_ReturnsToken(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"

	token, err := GenerateJWT(userID, email)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	// Token should be parseable
	if len(token) < 10 {
		t.Error("Token seems too short to be valid")
	}
}

func TestValidateJWT_ValidToken_ReturnsClaims(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"

	token, err := GenerateJWT(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %v, got %v", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
}

func TestValidateJWT_InvalidToken_ReturnsError(t *testing.T) {
	invalidToken := "invalid.jwt.token"

	_, err := ValidateJWT(invalidToken)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestValidateJWT_ExpiredToken_ReturnsError(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"

	// Create an expired token manually
	expirationTime := time.Now().Add(-1 * time.Hour) // 1 hour ago
	claims := &CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(getJWTSecret())

	_, err := ValidateJWT(tokenString)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestGenerateRefreshToken_ReturnsNonEmptyToken(t *testing.T) {
	token, err := GenerateRefreshToken()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty refresh token")
	}

	// Token should be 44 characters (base64 URL encoded 32 bytes)
	if len(token) != 44 {
		t.Errorf("Expected token length 44, got %d", len(token))
	}
}

func TestGenerateRefreshToken_UniqueTokens(t *testing.T) {
	token1, err1 := GenerateRefreshToken()
	token2, err2 := GenerateRefreshToken()

	if err1 != nil || err2 != nil {
		t.Fatal("Unexpected error generating tokens")
	}

	if token1 == token2 {
		t.Error("Expected different tokens, got same token")
	}
}

func TestHashRefreshToken_ConsistentHash(t *testing.T) {
	token := "test-token-123"
	hash1 := HashRefreshToken(token)
	hash2 := HashRefreshToken(token)

	if hash1 != hash2 {
		t.Error("Expected same hash for same token")
	}
}

func TestHashRefreshToken_ValidHashFormat(t *testing.T) {
	token := "test-token-123"
	hash := HashRefreshToken(token)

	// SHA256 hash should be 64 hex characters
	if len(hash) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}

	// Should be valid hex
	_, err := hex.DecodeString(hash)
	if err != nil {
		t.Errorf("Expected valid hex string, got error: %v", err)
	}
}

func TestGetBearerToken_ValidHeader_ReturnsToken(t *testing.T) {
	authHeader := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

	token, err := GetBearerToken(authHeader)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token != expectedToken {
		t.Errorf("Expected token %s, got %s", expectedToken, token)
	}
}

func TestGetBearerToken_EmptyHeader_ReturnsError(t *testing.T) {
	_, err := GetBearerToken("")

	if err == nil {
		t.Error("Expected error for empty header")
	}
}

func TestGetBearerToken_InvalidPrefix_ReturnsError(t *testing.T) {
	testCases := []struct {
		name   string
		header string
	}{
		{"Missing prefix", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"},
		{"Wrong prefix", "Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"},
		{"Too short", "Bea"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GetBearerToken(tc.header)
			if err == nil {
				t.Errorf("Expected error for header: %s", tc.header)
			}
		})
	}
}

func TestGetBearerToken_EmptyToken_ReturnsError(t *testing.T) {
	_, err := GetBearerToken("Bearer ")

	if err == nil {
		t.Error("Expected error for empty token")
	}
}

func BenchmarkGenerateJWT(b *testing.B) {
	userID := uuid.New()
	email := "test@example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateJWT(userID, email)
	}
}

func BenchmarkGenerateRefreshToken(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateRefreshToken()
	}
}

func BenchmarkHashRefreshToken(b *testing.B) {
	token := "test-token-123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HashRefreshToken(token)
	}
}

// TestHashRefreshToken_ConsistencyWithSha256 tests that our hash function produces correct SHA256
func TestHashRefreshToken_ConsistencyWithSha256(t *testing.T) {
	token := "8UMBun-uTzNTmFoPErpxvqIZeI3UuFazIA3bjwp0S7w="

	// Calculate expected hash
	hasher := sha256.New()
	hasher.Write([]byte(token))
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	// Calculate actual hash
	actualHash := HashRefreshToken(token)

	if actualHash != expectedHash {
		t.Errorf("Expected hash %s, got %s", expectedHash, actualHash)
	}

	t.Logf("Hash: %s", actualHash)
}
