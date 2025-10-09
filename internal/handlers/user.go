package handlers

import (
	"net/http"

	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

// HandlerGetUser returns the authenticated user's information
// Bu endpoint JWT ile korunuyor - sadece giriş yapmış kullanıcılar erişebilir
// Middleware user bilgisini zaten pass ediyor, direkt döndürüyoruz
func (cfg *Config) HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	models.RespondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}
