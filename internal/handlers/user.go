package handlers

import (
	"net/http"

	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

// HandlerGetUser returns the authenticated user's information
// Bu endpoint JWT ile korunuyor - sadece giriş yapmış kullanıcılar erişebilir
// Middleware user bilgisini zaten pass ediyor, direkt döndürüyoruz
// @Summary     Get current user
// @Description Get the authenticated user's information
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Success     200  {object}  object  "User information"
// @Failure     401  {object}  object  "Unauthorized"
// @Router      /v1/users/me [get]
func (cfg *Config) HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	models.RespondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}
