package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
	"github.com/mehmettalhairmak/rss-aggregator/internal/realtime"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary     WebSocket connection
// @Description Establishes a WebSocket connection for real-time updates. The connection requires authentication via JWT token passed as query parameter. Once connected, clients receive real-time notifications when new posts are available from their followed feeds.
// @Tags        websocket
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Param       token  query     string  true  "JWT access token for authentication"
// @Success     101    {string}  string  "Switching Protocols - WebSocket connection established"
// @Failure     400     {object}  object  "Bad request - Invalid token or connection error"
// @Failure     401     {object}  object  "Unauthorized - Invalid or missing token"
// @Failure     500     {object}  object  "Internal server error"
// @Router      /v1/ws [get]
// @Note        This endpoint upgrades HTTP connection to WebSocket. Use WebSocket client libraries (e.g., gorilla/websocket) to connect. The connection remains open and receives JSON messages with new post updates in real-time.
func (cfg *Config) HandlerWebsocket(w http.ResponseWriter, r *http.Request, user database.User) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	client := realtime.NewClient(cfg.Hub, conn, user.ID)
	cfg.Hub.RegisterClient(client)

	go client.WritePump()
	client.ReadPump()
}
