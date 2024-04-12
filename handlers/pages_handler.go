package handlers

import (
	"context"
	"net/http"

	"github.com/DeepAung/anon-chat/pkg/hub"
	"github.com/DeepAung/anon-chat/pkg/utils"
	"github.com/DeepAung/anon-chat/views"
)

type pagesHandler struct {
	hub *hub.Hub
}

func NewPagesHandler(hub *hub.Hub) *pagesHandler {
	return &pagesHandler{
		hub: hub,
	}
}

func (h *pagesHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		utils.Render(views.ErrorPage("path not found"), w)
		return
	}

	if err := views.Index().Render(context.Background(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *pagesHandler) Chat(w http.ResponseWriter, r *http.Request) {
	connectUrl := "/ws/connect?" + r.URL.RawQuery

	roomId := r.URL.Query().Get("roomId")
	room, ok := h.hub.GetRoom(roomId)
	if !ok {
		utils.Render(views.ErrorPage("room id not found"), w)
		return
	}

	if err := views.Chat(room, *room.MsgHistory.Iter(), connectUrl).Render(context.Background(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
