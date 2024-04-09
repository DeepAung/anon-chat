package utils

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func Render(component templ.Component, w http.ResponseWriter) {
	if err := component.Render(context.Background(), w); err != nil {
		http.Error(w, "error: cannot render templ component", http.StatusInternalServerError)
	}
}
