package handlers

import (
	"context"
	"net/http"

	"github.com/DeepAung/anon-chat/server/utils"
	"github.com/DeepAung/anon-chat/server/views"
)

type pagesHandler struct{}

func NewPagesHandler() *pagesHandler {
	return &pagesHandler{}
}

func (h *pagesHandler) Login(w http.ResponseWriter, r *http.Request) {
	if utils.HasCookie(r, "username") {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	}

	if err := views.Login().Render(context.Background(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *pagesHandler) Index(w http.ResponseWriter, r *http.Request) {
	if !utils.HasCookie(r, "username") {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}

	err := views.Index().Render(context.Background(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
