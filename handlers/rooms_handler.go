package handlers

import (
	"net/http"

	"github.com/DeepAung/anon-chat/hub"
	"github.com/DeepAung/anon-chat/services"
)

type roomsHandler struct {
	hub      *hub.Hub
	usersSvc *services.UsersService
	roomsSvc *services.RoomsService
}

func NewRoomsHandler(
	hub *hub.Hub,
	usersSvc *services.UsersService,
	roomsSvc *services.RoomsService,
) *roomsHandler {
	return &roomsHandler{
		hub:      hub,
		usersSvc: usersSvc,
		roomsSvc: roomsSvc,
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

	h.usersSvc.Login(w, username)
	h.roomsSvc.SetRoomName(w, roomName)

	w.Header().Add("HX-Redirect", "/chat")
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

	h.usersSvc.Login(w, username)
	h.roomsSvc.SetRoomId(w, roomId)

	w.Header().Add("HX-Redirect", "/chat")
}

func (h *roomsHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	h.usersSvc.Logout(w)
	h.roomsSvc.DeleteCookies(w)

	w.Header().Add("HX-Redirect", "/")
}
