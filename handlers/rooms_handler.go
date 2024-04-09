package handlers

import (
	"net/http"

	"github.com/DeepAung/anon-chat/hub"
	"github.com/DeepAung/anon-chat/utils"
	"github.com/DeepAung/anon-chat/views"
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
		utils.Render(views.ErrorMsg("no username"), w)
		return
	} else if roomName == "" {
		utils.Render(views.ErrorMsg("no room name"), w)
		return
	}

	roomId := h.hub.CreateRoom(roomName)
	url, err := utils.SetQueries("/chat", map[string]string{
		"username": username,
		"roomId":   roomId,
	})
	if err != nil {
		utils.Render(views.ErrorBody("gen url error: "+err.Error()), w)
		return
	}

	w.Header().Add("HX-Redirect", url)
}

func (h *roomsHandler) Connect(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	roomId := r.FormValue("roomId")

	if username == "" {
		utils.Render(views.ErrorMsg("no username"), w)
		return
	} else if roomId == "" {
		utils.Render(views.ErrorMsg("no room id"), w)
		return
	}

	url, err := utils.SetQueries("/chat", map[string]string{
		"username": username,
		"roomId":   roomId,
	})
	if err != nil {
		utils.Render(views.ErrorBody("gen url error: "+err.Error()), w)
		return
	}

	w.Header().Add("HX-Redirect", url)
}
