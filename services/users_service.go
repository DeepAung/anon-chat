package services

import (
	"net/http"
	"time"

	"github.com/DeepAung/anon-chat/utils"
	"github.com/google/uuid"
)

type UsersService struct{}

func NewUsersService() *UsersService {
	return &UsersService{}
}

func (h *UsersService) Login(w http.ResponseWriter, username string) {
	id := uuid.NewString()

	expires := time.Now().AddDate(0, 0, 1)
	utils.SetCookie(w, "username", username, expires, 86400)
	utils.SetCookie(w, "id", id, expires, 86400)
}

func (h *UsersService) Logout(w http.ResponseWriter) {
	utils.DeleteCookie(w, "username")
	utils.DeleteCookie(w, "id")
}
