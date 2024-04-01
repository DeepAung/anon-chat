package handlers

import (
	"net/http"

	"github.com/DeepAung/anon-chat/server/hub"
)

type roomsHandler struct {
	hub *hub.Hub
}

func NewRoomsHandler(hub *hub.Hub) *roomsHandler {
	return &roomsHandler{
		hub: hub,
	}
}

func (h *roomsHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
}

func (h *roomsHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
}

func (h *roomsHandler) LeaveRoom(w http.ResponseWriter, r *http.Request) {
}
