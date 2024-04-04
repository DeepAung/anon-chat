package handlers

import (
	"fmt"
	"net/http"

	"github.com/DeepAung/anon-chat/server/hub"
	"github.com/DeepAung/anon-chat/server/services"
)

type roomsHandler struct {
	hub      *hub.Hub
	usersSvc *services.UsersService
}

func NewRoomsHandler(hub *hub.Hub, usersSvc *services.UsersService) *roomsHandler {
	return &roomsHandler{
		hub:      hub,
		usersSvc: usersSvc,
	}
}

func (h *roomsHandler) CreateConnect(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	roomName := r.FormValue("roomName")
	roomId := r.FormValue("roomId")

	if username == "" {
		http.Error(w, "no username", http.StatusBadRequest)
		return
	}

	h.usersSvc.Login(w, username)

	if roomId != "" {
		http.Redirect(w, r, fmt.Sprintf("/chat?roomId=%s", roomId), http.StatusMovedPermanently)
	} else if roomName != "" {
		http.Redirect(w, r, fmt.Sprintf("/chat?roomName=%s", roomName), http.StatusMovedPermanently)
	} else {
		http.Error(w, "no room id or room name", http.StatusBadRequest)
	}
}

func (h *roomsHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	h.usersSvc.Logout(w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
