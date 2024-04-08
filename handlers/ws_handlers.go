package handlers

import (
	"context"
	"net/http"

	"github.com/DeepAung/anon-chat/hub"
	"github.com/DeepAung/anon-chat/views"
	"golang.org/x/net/websocket"
)

type wsHandler struct {
	hub *hub.Hub
}

func NewWsHandler(hub *hub.Hub) *wsHandler {
	return &wsHandler{
		hub: hub,
	}
}

func (h *wsHandler) Connect(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	roomId := query.Get("roomId")
	username := query.Get("username")
	if username == "" {
		username = "Anonymous User"
	}

	websocket.Handler(func(ws *websocket.Conn) {
		err := h.hub.Connect(ws, username, roomId)
		defer ws.Close()

		if err != nil {
			_ = views.ErrorMessage(err.Error()).Render(context.Background(), ws)
			return
		}

		h.hub.Listen(ws, roomId)
	}).ServeHTTP(w, r)
}
