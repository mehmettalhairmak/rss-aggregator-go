package realtime

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Hub struct {
	clients    map[uuid.UUID]*Client
	register   chan *Client
	unregister chan *Client
	signal     chan map[uuid.UUID][]byte
	Logger     zerolog.Logger
}

func NewHub(l zerolog.Logger) *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		signal:     make(chan map[uuid.UUID][]byte),
		Logger:     l,
	}
}

func (hub *Hub) Run() {
	hub.Logger.Info().Msg("Realtime Hub started running.")

	for {
		select {
		case client := <-hub.register:
			hub.clients[client.userID] = client

			hub.Logger.Info().
				Str("user_id", client.userID.String()).
				Int("total_clients", len(hub.clients)).
				Msg("Client registered successfully.")
		case client := <-hub.unregister:
			if _, ok := hub.clients[client.userID]; ok {
				delete(hub.clients, client.userID)
				close(client.send)

				hub.Logger.Warn().
					Str("user_id", client.userID.String()).
					Int("total_clients", len(hub.clients)).
					Msg("Client unregistered. Connection closed.")
			}
		case signals := <-hub.signal:
			for userID, payload := range signals {
				if client, ok := hub.clients[userID]; ok {
					select {
					case client.send <- payload:
					default:
						hub.Logger.Error().
							Str("user_id", userID.String()).
							Msg("Client send channel is full! Disconnecting misbehaving client.")

						close(client.send)
						delete(hub.clients, userID)
					}
				}
			}
		}
	}
}

func (hub *Hub) RegisterClient(c *Client) {
	hub.register <- c
}

func (hub *Hub) SendSignal(signals map[uuid.UUID][]byte) {
	hub.signal <- signals
}
