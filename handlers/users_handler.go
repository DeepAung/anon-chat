package handlers

import (
	"net/http"
	"time"

	"github.com/DeepAung/anon-chat/server/utils"
	"github.com/google/uuid"
)

type usersHandler struct{}

func NewUsersHandler() *usersHandler {
	return &usersHandler{}
}

func (h *usersHandler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	id := uuid.NewString()

	expires := time.Now().AddDate(0, 0, 1)
	utils.SetCookie(w, "username", username, expires, 86400)
	utils.SetCookie(w, "id", id, expires, 86400)

	http.Redirect(w, r, "/index", http.StatusFound)
}

func (h *usersHandler) Logout(w http.ResponseWriter, r *http.Request) {
	utils.DeleteCookie(w, "username")
	utils.DeleteCookie(w, "id")

	http.Redirect(w, r, "/login", http.StatusFound)
}
