package middleware

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/mehmettalhairmak/rss-aggregator/internal/auth"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

// AuthedHandler is a handler function that requires authentication
// Unlike normal http.HandlerFunc, it takes database.User as a third parameter
// Açıklama: Protected endpoint'ler için özel handler tipi
type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

// Config holds dependencies for middleware
type Config struct {
	DB *database.Queries
}

// NewConfig creates a new middleware config
func NewConfig(db *database.Queries) *Config {
	return &Config{
		DB: db,
	}
}

// Auth wraps an authenticated handler with JWT authentication logic
// Açıklama: JWT Middleware - Her protected endpoint'e gelen request'i kontrol eder
// Flow:
// 1. Authorization header'dan "Bearer token" formatında JWT alır
// 2. Token'ı validate eder (signature, expiration check)
// 3. Token'dan user_id'yi çıkarır
// 4. User'ı database'den bulur
// 5. Handler'a user bilgisini pass eder
func (cfg *Config) Auth(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authorization header'ı al
		// Format: "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			models.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// "Bearer " prefix'ini kaldır ve token'ı al
		token, err := auth.GetBearerToken(authHeader)
		if err != nil {
			models.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid authorization header: %v", err))
			return
		}

		// JWT token'ı validate et
		// Bu fonksiyon token'ın:
		// - Signature'ını kontrol eder (bizim secret key ile imzalanmış mı)
		// - Expiration time'ını kontrol eder (süresi dolmuş mu)
		// - Claim'leri parse eder (user_id, email vs.)
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			models.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			return
		}

		// Token'dan aldığımız user_id ile database'den user'ı bul
		user, err := cfg.DB.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				models.RespondWithError(w, http.StatusUnauthorized, "User not found")
				return
			}
			models.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
			return
		}

		// Kullanıcı bulundu! Handler'ı çağır ve user bilgisini pass et
		// Artık handler içinde user.ID, user.Email vs. kullanılabilir
		handler(w, r, user)
	}
}
