package handlers

import (
	"net/http"

	"github.com/DeepAung/anon-chat/pkg/hub"
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
		h.hub.ConnectAndListen(ws, username, roomId)
	}).ServeHTTP(w, r)
}
