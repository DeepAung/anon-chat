package handlers

import (
	"context"
	"net/http"

	"github.com/DeepAung/anon-chat/services"
	"github.com/DeepAung/anon-chat/utils"
	"github.com/DeepAung/anon-chat/views"
)

type pagesHandler struct {
	usersSvc *services.UsersService
	roomsSvc *services.RoomsService
}

func NewPagesHandler(
	usersSvc *services.UsersService,
	roomsSvc *services.RoomsService,
) *pagesHandler {
	return &pagesHandler{
		usersSvc: usersSvc,
		roomsSvc: roomsSvc,
	}
}

func (h *pagesHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := views.Index().Render(context.Background(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *pagesHandler) Chat(w http.ResponseWriter, r *http.Request) {
	if !utils.HasCookie(r, "userId") {
		h.usersSvc.Logout(w)
		h.roomsSvc.DeleteCookies(w)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	roomName := utils.GetCookieValue(r, "roomName")
	roomId := utils.GetCookieValue(r, "roomId")
	if roomName == "" && roomId == "" {
		h.usersSvc.Logout(w)
		h.roomsSvc.DeleteCookies(w)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	if err := views.Chat(roomName, roomId).Render(context.Background(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
