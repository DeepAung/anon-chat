package handlers

import (
	"net/http"

	"github.com/DeepAung/anon-chat/hub"
	"github.com/DeepAung/anon-chat/utils"
)

type roomsHandler struct {
	hub *hub.Hub
}

func NewRoomsHandler(hub *hub.Hub) *roomsHandler {
	return &roomsHandler{
		hub: hub,
	}
}

func (h *roomsHandler) CreateAndConnect(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	roomName := r.FormValue("roomName")

	if username == "" {
		http.Error(w, "no username", http.StatusBadRequest)
		return
	} else if roomName == "" {
		http.Error(w, "no room name", http.StatusBadRequest)
		return
	}

	roomId := h.hub.Create(roomName)
	url, err := utils.SetQueries("/chat", map[string]string{
		"username": username,
		"roomId":   roomId,
	})
	if err != nil {
		http.Error(w, "gen url error: "+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("HX-Redirect", url)
}

func (h *roomsHandler) Connect(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	roomId := r.FormValue("roomId")

	if username == "" {
		http.Error(w, "no username", http.StatusBadRequest)
		return
	} else if roomId == "" {
		http.Error(w, "no room id", http.StatusBadRequest)
		return
	}

	url, err := utils.SetQueries("/chat", map[string]string{
		"username": username,
		"roomId":   roomId,
	})
	if err != nil {
		http.Error(w, "gen url error: "+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("HX-Redirect", url)
}

func (h *roomsHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Redirect", "/")
}
