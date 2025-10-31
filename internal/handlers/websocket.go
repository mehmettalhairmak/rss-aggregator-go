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
