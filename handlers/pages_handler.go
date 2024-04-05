package handlers

import (
	"context"
	"net/http"

	"github.com/DeepAung/anon-chat/utils"
	"github.com/DeepAung/anon-chat/views"
)

type pagesHandler struct{}

func NewPagesHandler() *pagesHandler {
	return &pagesHandler{}
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
	if !utils.HasCookie(r, "id") {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	}

	if err := views.Chat().Render(context.Background(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
